package plugins

import (
	"context"
	"math/rand"
	"time"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
	log "github.com/sirupsen/logrus"
)

type Repeat struct {
}

func (rep *Repeat) Do(ctx *context.Context, messageType public.MessageType, admins []string, gmap map[string][]string, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, useRole []string, isBot, isDirectMessage, botIsAdmin bool, priceSearch string, attachments []string) utils.RetStuct {

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(101)

	isBotAdmin := public.IsBotAdmin(userId, admins)
	ggk, _ := database.GetJudgeKeys()
	containsJudgeKeys := database.Judge(msg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" && !isBotAdmin {
		msg := "消息触发守卫，已被拦截"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
		}
	}

	if len(msg) < 20 && r%70 == 0 && !(public.StartsWith(msg, ".") || public.StartsWith(msg, "%") || public.StartsWith(msg, "％")) {
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType: utils.GuildMsg,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("复读", &Repeat{})
}
