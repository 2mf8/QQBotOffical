package plugins

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
)

type Guard struct {
}

func (guard *Guard) Do(ctx *context.Context, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, isBot, isDirectMessage, botIsAdmin, isBotAdmin, isAdmin bool, priceSearch string) utils.RetStuct {
	if !botIsAdmin {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	guardIntent := int64(database.PluginGuard)
	sg, _ := database.SGBGIACI(guildId, channelId)
	isGuard := sg.PluginSwitch.IsCloseOrGuard & guardIntent

	if isGuard > 0 {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	ggk, _ := database.GetJudgeKeys()

	if public.StartsWith(msg, ".拦截") && (isAdmin || isBotAdmin) {
		vocabulary := strings.TrimPrefix(msg, ".拦截")
		content := strings.Split(vocabulary, " ")
		err := ggk.JudgeKeysUpdate(content...)
		if err != nil {
			log.Warnln(err)
		}
		msg := "拦截词汇添加成功"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType: utils.GuildMsg,
		}
	}

	if public.StartsWith(msg, ".取消拦截") && isBotAdmin {
		vocabulary := strings.TrimPrefix(msg, ".取消拦截")
		content := strings.Split(vocabulary, " ")
		ggk.JudgeKeysDelete(content...)
		msg := "拦截词汇删除成功"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType: utils.GuildMsg,
		}
	}

	containsJudgeKeys := database.Judge(msg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
		if isAdmin {
			msg := "消息触发守卫，已被拦截"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, msg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: msg,
				},
				ReqType: utils.GuildMsg,
			}
		}
		msg := "消息触发守卫，已撤回消息并禁言该用户(<@!" + userId + ">)两分钟, 请文明发言"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType:  utils.DeleteMsg,
			Duration: "120",
			MsgId:    msgId,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("守卫", &Guard{})
}
