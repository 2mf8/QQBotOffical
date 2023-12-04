package plugins

import (
	"context"
	"net/url"
	"strings"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
	log "github.com/sirupsen/logrus"
)

type ScramblePlugin struct {
}

func (scramble *ScramblePlugin) Do(ctx *context.Context, messageType public.MessageType, admins []string, gmap map[string][]string, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, useRole []string, isBot, isDirectMessage, botIsAdmin bool, priceSearch string, attachments []string) utils.RetStuct {

	s, b := public.Prefix(msg, ".", messageType)
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	messageType = public.Undefined
	tn := database.Tnoodle(s)
	ins := tn.Instruction
	shor := tn.ShortName
	show := tn.ShowName
	if ins != "instruction" {
		gs := database.GetScramble(shor)
		if public.StartsWith(gs, "net") || gs == "获取失败" {
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> 获取打乱失败", guildId, channelId, userId)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: "获取打乱失败",
				},
				ReqType: utils.GuildMsg,
			}
		}
		if shor == "minx" {
			gs = strings.Replace(gs, "U' ", "#\n", -1)
			gs = strings.Replace(gs, "U ", "U\n", -1)
			gs = strings.Replace(gs, "#", "U'", -1)
		}
		imgUrl := "http://2mf8.cn:2014/view/" + shor + ".png?scramble=" + url.QueryEscape(strings.Replace(gs, "\n", " ", -1))
		sc := show + "\n" + gs
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %v\n%v<image url=\"%v\"/>", guildId, channelId, userId, show, gs, imgUrl)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text:  sc,
				Image: imgUrl,
			},
			ReqType: utils.GuildMsg,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("打乱", &ScramblePlugin{})
}
