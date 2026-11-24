package dgfeishu_test

import (
	"os"
	"testing"
	"time"

	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-common/utils"
	dgfeishu "github.com/darwinOrg/go-feishu"
	dglogger "github.com/darwinOrg/go-logger"
)

func TestSimpleSendTextMessage(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	resp, err := dgfeishu.SimpleSendTextMessage(ctx, os.Getenv("WEBHOOK"), "验证码：000000")
	if err != nil {
		dglogger.Error(ctx, err)
		return
	}
	dglogger.Info(ctx, resp)
}

func TestSimpleSendImageMessage(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	resp, err := dgfeishu.SimpleSendImageMessage(ctx, os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), os.Getenv("WEBHOOK"), "test.jpg")
	if err != nil {
		dglogger.Error(ctx, err)
		return
	}
	dglogger.Info(ctx, resp)
}

func TestSimpleSendFileMessage(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	resp, err := dgfeishu.SimpleSendFileMessage(ctx, os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), os.Getenv("CHAT_ID"), "test.jpg")
	if err != nil {
		dglogger.Error(ctx, err)
		return
	}
	dglogger.Info(ctx, resp)
}

func TestSimpleSendTableCard(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	resp, err := dgfeishu.SimpleSendTableCard(ctx, os.Getenv("WEBHOOK"), "QA销售数据",
		[]string{"客户", "金额", "状态"},
		[][]string{
			{"飞书科技", "168", "已签约"},
			{"某公司", "26.5", "跟进中"},
		})
	if err != nil {
		dglogger.Error(ctx, err)
		return
	}
	dglogger.Info(ctx, resp)
}

func TestGetFirstWsMsg(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	wsMsg, err := dgfeishu.GetFirstWsMsg(ctx, time.Minute)
	if err != nil {
		dglogger.Error(ctx, err)
		return
	}
	realText, _ := wsMsg.GetRealText()
	dglogger.Infof(ctx, "wsMsg: %s | realText: %s", utils.MustConvertBeanToJsonStringPretty(wsMsg), realText)
}
