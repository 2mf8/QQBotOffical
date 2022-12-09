package plugins

import (
	"context"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
)

type BotSwitch struct {
}

func (botSwitch *BotSwitch) Do(ctx *context.Context, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, isBot, isDirectMessage, botIsAdmin, isBotAdmin, isAdmin bool, priceSearch string) utils.RetStuct {

	s, b := public.Prefix(msg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	if public.StartsWith(s, "开启") && (isAdmin || botIsAdmin) {
		s = strings.TrimSpace(strings.TrimPrefix(s, "开启"))
		if s == "开关" {
			reply := "开关无法开启或关闭自身"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		i := database.PluginNameToIntent(s)
		if i == 0 {
			reply := "功能不存在"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		err := database.SwitchSave(guildId, channelId, userId, int64(i), time.Now(), false)
		if err != nil {
			reply := "开启失败"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		} else {
			reply := "开启成功"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
	}

	if public.StartsWith(s, "关闭") && (isAdmin || botIsAdmin) {
		s = strings.TrimSpace(strings.TrimPrefix(s, "关闭"))
		if s == "开关" {
			reply := "开关无法开启或关闭自身"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		i := database.PluginNameToIntent(s)
		if i == 0 {
			reply := "功能不存在"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		err := database.SwitchSave(guildId, channelId, userId, int64(i), time.Now(), true)
		if err != nil {
			reply := "关闭失败"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		} else {
			reply := "关闭成功"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}

	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("开关", &BotSwitch{})
}
