package plugins

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/2mf8/QQBotOffical/public"
	. "github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
)

type Admin struct {
}

func (admin *Admin) Do(ctx *context.Context, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, isBot, isDirectMessage, botIsAdmin, isBotAdmin, isAdmin bool) utils.RetStuct {

	s, b := Prefix(msg, ".")

	if !b || !(public.StartsWith(s, "jin") || public.StartsWith(s, "jie") || public.StartsWith(s, "t") || public.StartsWith(s, "T")) {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	if !botIsAdmin {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: "机器人不是管理员，无法进行禁言或踢人操作",
			},
		}
	}

	// <@!13970278473675774808>

	jinTime := "-1"

	reg1 := regexp.MustCompile("@!")
	reg2 := regexp.MustCompile("@")
	reg3 := regexp.MustCompile(">")
	reg4 := regexp.MustCompile("  ")

	str1 := strings.TrimSpace(reg1.ReplaceAllString(s, "at qq=\""))
	str1 = strings.TrimSpace(reg2.ReplaceAllString(str1, "at qq=\""))
	str2 := strings.TrimSpace(reg3.ReplaceAllString(str1, "\"/>"))

	for Contains(str2, "  ") {
		str2 = strings.TrimSpace(reg4.ReplaceAllString(str2, " "))
	}

	cstr, cstrs := public.GuildAtConvert(str2)
	if (public.ConvertTime(cstr) < (30*24*60*60 - 60)) && public.ConvertTime(cstr) > 0 {
		jinTime = strconv.Itoa(int(public.ConvertTime(cstr) + 1))
	}

	if public.StartsWith(str2, "jin") {
		if len(cstrs) == 0 {
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: "禁言用户不能为空",
				},
				ReqType: utils.GuildBan,
			}
		}
		if jinTime == "-1" {
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: "禁言时间有误或超过最大允许范围",
				},
				ReqType: utils.GuildBan,
			}
		}
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: public.ConvertJinTime(int(public.ConvertTime(cstr) + 1)),
			},
			BanId:    cstrs,
			Duration: jinTime,
			ReqType:  utils.GuildBan,
		}
	}

	if public.StartsWith(str2, "jie") {
		if len(cstrs) == 0 {
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: "解禁用户不能为空",
				},
				ReqType: utils.RelieveBan,
			}
		}
		return utils.RetStuct{
			RetVal:   utils.MESSAGE_BLOCK,
			BanId:    cstrs,
			Duration: "0",
			ReqType:  utils.RelieveBan,
		}
	}

	if public.StartsWith(str2, "t") || public.StartsWith(str2, "T") {
		rejectAddAgain := public.StartsWith(str2, "T")
		retract := 0
		if public.StartsWith(str2, "ti") || public.StartsWith(str2, "Ti") {
			retract = -1
		}
		if len(cstrs) == 0 {
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: "被踢用户不能为空",
				},
				ReqType: utils.GuildKick,
			}
		}
		return utils.RetStuct{
			RetVal:         utils.MESSAGE_BLOCK,
			ReqType:        utils.GuildKick,
			BanId:          cstrs,
			RejectAddAgain: rejectAddAgain,
			Retract:        retract,
		}
	}

	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}

}

func init() {
	utils.Register("群管", &Admin{})
}
