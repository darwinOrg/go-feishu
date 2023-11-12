package dgfeishu

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/go-lark/lark"
)

func SimpleSendTextMessage(ctx *dgctx.DgContext, webhook string, content string) (*lark.PostNotificationV2Resp, error) {
	dglogger.Infof(ctx, "simple send text message, webhook: %s, content: %s", webhook, content)
	bot := lark.NewNotificationBot(webhook)
	return bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text(content).Build())
}
