package dgfeishu

import (
	"errors"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/go-lark/lark"
)

func SimpleSendTextMessage(ctx *dgctx.DgContext, webhook, content string) (*lark.PostNotificationV2Resp, error) {
	dglogger.Infof(ctx, "simple send text message, content: %s", content)
	bot := lark.NewNotificationBot(webhook)
	return bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text(content).Build())
}

func SimpleSendImageMessage(ctx *dgctx.DgContext, appId, appSecret string, imageFile string) (*lark.PostNotificationV2Resp, error) {
	dglogger.Infof(ctx, "simple send image message, imageFile: %s", imageFile)
	bot := lark.NewChatBot(appId, appSecret)
	_, _ = bot.GetTenantAccessTokenInternal(true)
	//bot := lark.NewNotificationBot(webhook)
	//bot.SetDomain("https://open.feishu.cn")
	resp, err := bot.UploadImage(imageFile)
	if err != nil {
		dglogger.Errorf(ctx, "upload image[%s] error: %s", imageFile, err)
		return nil, err
	}
	if resp.Code != 0 {
		dglogger.Errorf(ctx, "upload image[%s] error: %s", imageFile, resp.Msg)
		return nil, errors.New(resp.Msg)
	}

	return bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgImage).Image(resp.Data.ImageKey).Build())
}
