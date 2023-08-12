package plugins

import (
	"context"
	"fmt"
	"strings"
	"time"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

type LearnPlugin struct {
}

func (learnPlugin *LearnPlugin) Do(ctx *context.Context, admins []string, gmap map[string][]string, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, useRole []string, isBot, isDirectMessage, botIsAdmin bool, priceSearch string, imgs []string) utils.RetStuct {

	s, b := public.Prefix(msg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	isAdmin := public.IsAdmin(useRole)
	isBotAdmin := public.IsBotAdmin(userId, admins)
	ggk, _ := database.GetJudgeKeys()
	containsJudgeKeys := database.Judge(msg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" && !isBotAdmin {
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

	if public.StartsWith(s, "#+") && (isBotAdmin || isAdmin) {
		if !isBotAdmin {
			suspectedUrl := strings.Split(s, ".")
			if len(suspectedUrl) > 1 && !isBotAdmin {
				if suspectedUrl[0] != "" && suspectedUrl[1] != "" {
					msg := "疑似网址(暂不支持网址)，已被拦截"
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
		}
		str2 := strings.TrimSpace(strings.TrimPrefix(s, "#+"))
		str3 := strings.Split(str2, "##")
		if len(str3) != 2 {
			if strings.TrimSpace(str3[0]) == "" {
				replyText := "问指令不能为空"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			err := database.LDBGAA(guildId, channelId, str3[0])
			if err != nil {
				replyText := "问答删除失败"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			replyText := "问答删除成功"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GuildMsg,
			}
		}
		if strings.TrimSpace(str3[0]) == "" {
			replyText := "问指令不能为空"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GuildMsg,
			}
		}
		iua := ""
		if len(imgs) > 0 {
			for _, iu := range imgs {
				iua += "#url#" + iu
			}
		}
		ans := str3[1] + iua
		err := database.LearnSave(strings.TrimSpace(str3[0]), guildId, channelId, userId, null.NewString(ans, true), time.Now(), true)
		fmt.Println(err)
		if err != nil {
			replyText := "添加失败"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GuildMsg,
			}
		}
		replyText := "学习已完成，下次触发有效"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyText,
			},
			ReqType: utils.GuildMsg,
		}
	}
	if public.StartsWith(s, "++") && isBotAdmin {
		str2 := strings.TrimSpace(strings.TrimPrefix(s, "++"))
		str3 := strings.Split(str2, "##")
		if len(str3) != 2 {
			if strings.TrimSpace(str3[0]) == "" {
				replyText := "系统问指令不能为空"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			err := database.LDBGAA("sys", "sys", str3[0])
			if err != nil {
				replyText := "系统问答删除失败"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			replyText := "系统问答删除成功"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GuildMsg,
			}
		}
		if strings.TrimSpace(str3[0]) == "" {
			replyText := "系统问指令不能为空"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GuildMsg,
			}
		}
		iua := ""
		if len(imgs) > 0 {
			for _, iu := range imgs {
				iua += "#url#" + iu
			}
		}
		ans := str3[1] + iua
		err := database.LearnSave(strings.TrimSpace(str3[0]), "sys", "sys", userId, null.NewString(ans, true), time.Now(), true)
		fmt.Println(err)
		if err != nil {
			replyText := "系统问答添加失败"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GuildMsg,
			}
		}
		replyText := "系统问答学习已完成，下次触发有效"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyText,
			},
			ReqType: utils.GuildMsg,
		}
	}
	if strings.TrimSpace(msg) == "" {
		replyText := "指令不能为空"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: replyText,
			},
			ReqType: utils.GuildMsg,
		}
	}
	learn_get, err := database.LearnGet(guildId, channelId, strings.TrimSpace(s))
	log.Println(learn_get.Answer.String, "ceshil", err)
	if err != nil || learn_get.Answer.String == "" {
		imgs := []string{}
		sys_learn_get, err := database.LearnGet("sys", "sys", strings.TrimSpace(s))
		log.Println(sys_learn_get.Answer.String, "ceshilsys", err)
		if sys_learn_get.Answer.String != "" {
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, sys_learn_get.Answer.String)
			if public.Contains(sys_learn_get.Answer.String, "#url#") {
				anst := ""
				anscontainimg := strings.Split(sys_learn_get.Answer.String, "#url#")
				for i, aci := range anscontainimg {
					if i > 0 {
						imgs = append(imgs, aci)
					} else {
						anst += aci
					}
				}
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:   anst,
						Images: imgs,
					},
					ReqType: utils.GuildMsg,
				}
			}
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: sys_learn_get.Answer.String,
				},
				ReqType: utils.GuildMsg,
			}
		}
	}
	if learn_get.Answer.String != "" {
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, learn_get.Answer.String)
		if public.Contains(learn_get.Answer.String, "#url#") {
			anst := ""
			anscontainimg := strings.Split(learn_get.Answer.String, "#url#")
			for i, aci := range anscontainimg {
				if i > 0 {
					imgs = append(imgs, aci)
				} else {
					anst += aci
				}
			}
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text:   anst,
					Images: imgs,
				},
				ReqType: utils.GuildMsg,
			}
		}
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: learn_get.Answer.String,
			},
			ReqType: utils.GuildMsg,
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}
func init() {
	utils.Register("学习", &LearnPlugin{})
}
