package plugins

import (
	"context"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
)

type Guard struct {
}

func (guard *Guard) Do(ctx *context.Context, messageType public.MessageType, admins []string, gmap map[string][]string, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, useRole []string, isBot, isDirectMessage, botIsAdmin bool, priceSearch string, attachments []string) utils.RetStuct {
	if !botIsAdmin {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	isAdmin := public.IsAdmin(useRole)
	isCommonMember := IsCommonMember(guildId, gmap, useRole)
	isBotAdmin := public.IsBotAdmin(userId, admins)
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

	if public.StartsWith(msg, ".取消拦截") && (isAdmin || isBotAdmin) {
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

	if Pass(useRole) {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	containsJudgeKeys := database.Judge(msg, *ggk.JudgekeysSync)
	if (containsJudgeKeys != "" || strings.Contains(msg, "当前版本不支持查看")) && isCommonMember {
		msg := "消息触发守卫，已撤回消息并禁言该用户(<@!" + userId + ">)五分钟, 请文明发言"
		jinDur := "300"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, msg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: msg,
			},
			ReqType:  utils.DeleteMsg,
			Duration: jinDur,
			MsgId:    msgId,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func IsCommonMember(guildId string, gmap map[string][]string, roles []string) bool {
	for _, role := range roles {
		for _, passRole := range gmap[guildId] {
			if role == passRole {
				return false
			}
		}
	}
	return true
}

func Pass(r []string) bool {
	for _, v := range r {
		i, _ := strconv.Atoi(v)
		if i > 21 {
			return true
		}
	}
	return false
}

func init() {
	utils.Register("守卫", &Guard{})
}
