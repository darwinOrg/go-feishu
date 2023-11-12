package dgfeishu_test

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dgfeishu "github.com/darwinOrg/go-feishu"
	dglogger "github.com/darwinOrg/go-logger"
	"testing"
)

func TestSimpleSendTextMessage(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: "123"}
	resp, err := dgfeishu.SimpleSendTextMessage(ctx, "https://open.feishu.cn/open-apis/bot/v2/hook/102301d5-cf25-4269-b296-bf6df737643a", "验证码：000000")
	if err != nil {
		dglogger.Error(ctx, err)
		return
	}
	dglogger.Info(ctx, resp)
}
