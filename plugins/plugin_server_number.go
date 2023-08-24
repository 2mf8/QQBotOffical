package plugins

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
)

type SN struct {
}

func (sn *SN) Do(ctx *context.Context, admins []string, gmap map[string][]string, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, useRole []string, isBot, isDirectMessage, botIsAdmin bool, priceSearch string, attachments []string) utils.RetStuct {

	isBotAdmin := public.IsBotAdmin(userId, admins)

	s, b := public.Prefix(msg, ".")

	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	snGet, _ := database.ServerNumbersGet()
	if public.StartsWith(s, "服务号+") && isBotAdmin {
		msg := ""
		snString := strings.TrimPrefix(s, "服务号+")
		fmt.Println(snString)
		sna := strings.SplitN(snString, " ", 3)
		fmt.Println(sna, len(sna))
		if len(sna) == 3 {
			i, err := strconv.Atoi(strings.TrimSpace(sna[2]))
			if err != nil {
				msg = "格式错误"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, msg)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: msg,
					},
					ReqType: utils.GuildMsg,
				}
			}
			err = snGet.ServerNumberUpdate(strings.TrimSpace(sna[0]), strings.TrimSpace(sna[1]), i)
			if err != nil {
				log.Warnln(err)
			}
			msg = "服务号添加成功"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, msg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: msg,
				},
				ReqType: utils.GuildMsg,
			}
		} else {
			msg = "格式错误"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, msg)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: msg,
				},
				ReqType: utils.GuildMsg,
			}
		}
	}

	if public.StartsWith(s, "服务号-") && isBotAdmin {
		snString := strings.TrimPrefix(s, "服务号-")
		value := snGet.ServerNumberSetSync.ServerNumbers[snString]
		snGet.ServerNumbersDelete(strings.TrimSpace(snString), value)
		msg := "服务号删除成功"
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
	utils.Register("服务号", &SN{})
}
