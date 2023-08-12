package plugins

import (
	"context"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
)

type Block struct{}

func (block *Block) Do(ctx *context.Context, admins []string, gmap map[string][]string, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, useRole []string, isBot, isDirectMessage, botIsAdmin bool, priceSearch string, attachments []string) utils.RetStuct {

	isBotAdmin := public.IsBotAdmin(userId, admins)
	ispblock, _ := database.PBlockGet(guildId, userId)
	if ispblock.UserId == userId && ispblock.IsPBlock {
		if !isBotAdmin {
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
			}
		}
	}

	s, b := public.Prefix(msg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	reg1 := regexp.MustCompile("@!")
	reg2 := regexp.MustCompile("@")
	reg3 := regexp.MustCompile(">")
	reg4 := regexp.MustCompile("  ")

	str1 := strings.TrimSpace(reg1.ReplaceAllString(s, "at qq=\""))
	str1 = strings.TrimSpace(reg2.ReplaceAllString(str1, "at qq=\""))
	str2 := strings.TrimSpace(reg3.ReplaceAllString(str1, "\"/>"))

	replyMsg := ""

	for public.Contains(str2, "  ") {
		str2 = strings.TrimSpace(reg4.ReplaceAllString(str2, " "))
	}

	_, cstrs := public.GuildAtConvert(str2)

	if public.StartsWith(s, "屏蔽+") && isBotAdmin {
		if len(cstrs) == 0 {
			replyMsg := "用户不存在"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyMsg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyMsg,
				},
				ReqType: utils.GuildMsg,
			}
		}

		for _, pb := range cstrs {
			err := database.PBlockSave(guildId, pb, userId, true, time.Now())
			if err != nil {
				replyMsg := "屏蔽<@!" + pb + ">失败"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyMsg)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyMsg,
					},
					ReqType: utils.GuildMsg,
				}
			}
			replyMsg += " <@!" + pb + ">"
		}
		replyMsg = "屏蔽" + replyMsg + " 成功"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyMsg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyMsg,
			},
			ReqType: utils.GuildMsg,
		}
	}
	if public.StartsWith(s, "屏蔽-") && isBotAdmin {
		if len(cstrs) == 0 {
			replyMsg := "用户不存在"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyMsg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyMsg,
				},
				ReqType: utils.GuildMsg,
			}
		}

		for _, pb := range cstrs {
			err := database.PBlockSave(guildId, pb, userId, false, time.Now())
			if err != nil {
				replyMsg := "解除屏蔽<@!" + pb + ">失败"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyMsg)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyMsg,
					},
					ReqType: utils.GuildMsg,
				}
			}
			replyMsg += " <@!" + pb + ">"
		}
		replyMsg = "解除屏蔽" + replyMsg + " 成功"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyMsg)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyMsg,
			},
			ReqType: utils.GuildMsg,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("屏蔽", &Block{})
}
