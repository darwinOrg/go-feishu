package dgfeishu

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/go-lark/lark"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

func SimpleSendTextMessage(ctx *dgctx.DgContext, webhook, content string) (*lark.PostNotificationV2Resp, error) {
	dglogger.Infof(ctx, "simple send text message, content: %s", content)
	bot := lark.NewNotificationBot(webhook)
	return bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text(content).Build())
}

func SimpleSendImageMessage(ctx *dgctx.DgContext, appId, appSecret, imageFile string) (*lark.PostMessageResponse, error) {
	dglogger.Infof(ctx, "simple send image message, imageFile: %s", imageFile)
	chatBot := lark.NewChatBot(appId, appSecret)
	_, _ = chatBot.GetTenantAccessTokenInternal(true)
	resp, err := chatBot.UploadImage(imageFile)
	if err != nil {
		dglogger.Errorf(ctx, "upload image[%s] error: %s", imageFile, err)
		return nil, err
	}
	if resp.Code != 0 {
		dglogger.Errorf(ctx, "upload image[%s] error: %s", imageFile, resp.Msg)
		return nil, errors.New(resp.Msg)
	}

	return chatBot.PostImage(resp.Data.ImageKey, lark.WithChatID(""))
}

var (
	appID             = os.Getenv("APP_ID")
	appSecret         = os.Getenv("APP_SECRET")
	verificationToken = os.Getenv("VERIFICATION_TOKEN")
	encryptKey        = os.Getenv("ENCRYPT_KEY")
	// 匹配 @_user_数字 的正则
	atUserRegex = regexp.MustCompile(`@_user_\d+`)
)

// FeishuMessage 提取需要的消息内容结构体
type FeishuMessage struct {
	ChatID    string `json:"chat_id"`    // 群ID
	MessageID string `json:"message_id"` // 消息ID
	OpenID    string `json:"open_id"`    // 发送人open_id
	Content   string `json:"content"`    // 消息文本
}

// InnerMsgContent 解析content内嵌json
type InnerMsgContent struct {
	Text string `json:"text"`
}

// GetRealText 从FeishuMessage中取出纯净消息文本
func (m *FeishuMessage) GetRealText() (string, error) {
	var inner InnerMsgContent
	if err := json.Unmarshal([]byte(m.Content), &inner); err != nil {
		return "", fmt.Errorf("解析content内嵌json失败: %w", err)
	}
	// 替换所有@_user_xxx为空字符串
	cleanText := atUserRegex.ReplaceAllString(inner.Text, "")
	// 清除前后多余空格
	return strings.TrimSpace(cleanText), nil
}

// GetFirstWsMsg 公共函数：建立WS，收到第一条消息后立即关闭连接返回内容
// 使用官方标准 OnP2MessageReceiveV1 事件分发器写法
func GetFirstWsMsg(ctx *dgctx.DgContext, waitTimeout time.Duration) (*FeishuMessage, error) {
	if ctx.GetInnerContext() == nil {
		ctx.SetInnerContext(context.Background())
	}
	parentCtx := ctx.GetInnerContext()

	msgChan := make(chan *FeishuMessage, 1)
	errChan := make(chan error, 1)

	// 全局控制ws生命周期的cancel ctx，用于强制终止重连
	wsCtx, wsCancel := context.WithCancel(context.Background())
	defer wsCancel()

	// 初始化事件分发器
	eventDispatcher := dispatcher.NewEventDispatcher(verificationToken, encryptKey).
		// 标准v2消息事件注册
		OnP2MessageReceiveV1(func(evCtx context.Context, event *larkim.P2MessageReceiveV1) error {
			// 提取消息主体
			msg := event.Event.Message
			sender := event.Event.Sender

			res := &FeishuMessage{
				ChatID:    *msg.ChatId,
				MessageID: *msg.MessageId,
				OpenID:    *sender.SenderId.OpenId,
				Content:   *msg.Content,
			}
			// 投递消息，主线程收到后关闭WS
			msgChan <- res
			return nil
		})

	// 关键：禁用自动重连，一次性连接不需要重试
	wsCli := larkws.NewClient(
		appID,
		appSecret,
		larkws.WithEventHandler(eventDispatcher),
		larkws.WithLogLevel(larkcore.LogLevelInfo),
		larkws.WithAutoReconnect(false), // 关闭断线自动重连！
	)

	// 协程启动WS连接
	go func() {
		dglogger.Info(ctx, "飞书WebSocket长连接启动...")
		// 阻塞运行，直到连接关闭/上下文取消
		_ = wsCli.Start(wsCtx)
		dglogger.Info(ctx, "ws内部循环已退出")
	}()

	// 等待消息/超时
	timeoutCtx, cancel := context.WithTimeout(parentCtx, waitTimeout)
	defer cancel()

	select {
	case <-timeoutCtx.Done():
		dglogger.Info(ctx, "等待消息超时，开始关闭ws")
		// 双重关闭：先取消上下文，再调用Close
		wsCli.Close()
		wsCancel()
		return nil, errors.New("等待飞书消息超时")
	case err := <-errChan:
		log.Printf("ws连接异常: %v", err)
		wsCli.Close()
		wsCancel()
		return nil, err
	case msg := <-msgChan:
		dglogger.Info(ctx, "收到消息，立即终止ws连接")
		// 收到消息，强制停止ws循环与重连
		wsCli.Close()
		wsCancel()
		return msg, nil
	}
}
