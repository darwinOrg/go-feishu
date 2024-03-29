package dgfeishu_test

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dgfeishu "github.com/darwinOrg/go-feishu"
	dglogger "github.com/darwinOrg/go-logger"
	"os"
	"testing"
)

func TestSimpleSendTextMessage(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: "123"}
	resp, err := dgfeishu.SimpleSendTextMessage(ctx, os.Getenv("WEBHOOK"), "验证码：000000")
	if err != nil {
		dglogger.Error(ctx, err)
		return
	}
	dglogger.Info(ctx, resp)
}

func TestSimpleSendImageMessage(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: "123"}
	resp, err := dgfeishu.SimpleSendImageMessage(ctx, os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), os.Getenv("WEBHOOK"), "test.jpg")
	if err != nil {
		dglogger.Error(ctx, err)
		return
	}
	dglogger.Info(ctx, resp)
}
