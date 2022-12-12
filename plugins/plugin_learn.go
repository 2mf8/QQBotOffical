package plugins

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

type LearnPlugin struct {
}

func (learnPlugin *LearnPlugin) Do(ctx *context.Context, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, isBot, isDirectMessage, botIsAdmin, isBotAdmin, isAdmin bool, priceSearch string) utils.RetStuct {

	s, b := public.Prefix(msg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	ggk, _ := database.GetJudgeKeys()
	containsJudgeKeys := database.Judge(msg, *ggk.JudgekeysSync)
	if containsJudgeKeys != "" {
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
	reg1 := regexp.MustCompile("＃")
	str1 := strings.TrimSpace(reg1.ReplaceAllString(s, "#"))
	if public.StartsWith(str1, "#+") && (isBotAdmin || isAdmin) {
		suspectedUrl := strings.Split(s, ".")
		if len(suspectedUrl) > 1 {
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
		str2 := strings.TrimSpace(strings.TrimPrefix(str1, "#+"))
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
		err := database.LearnSave(strings.TrimSpace(str3[0]), guildId, channelId, userId, null.NewString(str3[1], true), time.Now(), true)
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
	if public.StartsWith(str1, "++") && isBotAdmin {
		suspectedUrl := strings.Split(s, ".")
		if len(suspectedUrl) > 1 {
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
		str2 := strings.TrimSpace(strings.TrimPrefix(str1, "++"))
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
		err := database.LearnSave(strings.TrimSpace(str3[0]), "sys", "sys", userId, null.NewString(str3[1], true), time.Now(), true)
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
	//log.Println(learn_get.LearnSync.Answer.String,"ceshil", err)
	if err != nil || learn_get.LearnSync.Answer.String == "" {
		sys_learn_get, _ := database.LearnGet("sys", "sys", strings.TrimSpace(s))
		if sys_learn_get.LearnSync.Answer.String != "" {
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, sys_learn_get.LearnSync.Answer.String)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: sys_learn_get.LearnSync.Answer.String,
				},
				ReqType: utils.GuildMsg,
			}
		}
	}
	if learn_get.LearnSync.Answer.String != "" {
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, learn_get.LearnSync.Answer.String)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: learn_get.LearnSync.Answer.String,
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
