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
	dgfeishu.SimpleSendTextMessage(ctx, os.Getenv("WEBHOOK"), "验证码：000000")
}

func TestSimpleSendImageMessage(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	dgfeishu.SimpleSendImageMessage(ctx, os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), os.Getenv("WEBHOOK"), "test.jpg")
}

func TestSimpleSendFileMessage(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	dgfeishu.SimpleSendFileMessage(ctx, os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), os.Getenv("CHAT_ID"), "test.jpg")
}

func TestSimpleSendTableCard(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	dgfeishu.SimpleSendTableCard(ctx, os.Getenv("WEBHOOK"), "销售数据",
		[]string{"客户", "金额", "状态"},
		[]string{"auto", "120px", "40%"},
		[][]any{
			{"飞书科技", 168, "已签约"},
			{"某公司", 26.5, "跟进中"},
		})
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
