package plugins

import (
	"context"

	"github.com/2mf8/QQBotOffical/utils"
)

type Reply struct {
}

func (rep *Reply) Do(ctx *context.Context, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, isBot, isDirectMessage, botIsAdmin, isBotAdmin, isAdmin bool, priceSearch string, imgs []string) utils.RetStuct {
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("回复", &Reply{})
}
