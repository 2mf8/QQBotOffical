package plugins

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
)

type PricePlugin struct {
}

func (price *PricePlugin) Do(ctx *context.Context, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, isBot, isDirectMessage, botIsAdmin, isBotAdmin, isAdmin bool, priceSearch string, imgs []string) utils.RetStuct {

	reg1 := regexp.MustCompile("％")
	reg2 := regexp.MustCompile("＃")
	reg3 := regexp.MustCompile("＆")
	str1 := strings.TrimSpace(reg1.ReplaceAllString(msg, "%"))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, "#"))
	str3 := strings.TrimSpace(reg3.ReplaceAllString(str2, "&"))

	s, b := public.Prefix(str3, "%")
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

	if public.Contains(priceSearch, "黄小姐") {
		if public.StartsWith(s, "#+") && (isAdmin || isBotAdmin) {
			str4 := strings.TrimSpace(string([]byte(s)[len("#+"):]))
			str5 := strings.Split(str4, "##")
			if len(str5) != 2 {
				if strings.TrimSpace(str5[0]) == "" {
					replyText := "商品名不能为空"
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: replyText,
						},
						ReqType: utils.GuildMsg,
					}
				}

				err := database.IDBGAN("10001", "10001", str5[0])
				if err != nil {
					replyText := "删除失败"
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: replyText,
						},
						ReqType: utils.GuildMsg,
					}
				}
				replyText := "删除成功"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if strings.TrimSpace(str5[0]) == "" {
				replyText := "商品名不能为空"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			str6 := strings.Split(str5[1], "#?")
			if len(str6) != 2 {
				err := database.ItemSave("10001", "10001", null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, userId, null.NewTime(time.Now(), true))
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
				replyText := "添加成功"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)

				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			err := database.ItemSave("10001", "10001", null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), userId, null.NewTime(time.Now(), true))
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
			replyText := "添加成功"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GuildMsg,
			}
		}

		ps := ""
		psc := ""
		ic := 0
		cps, _ := database.GetItems("10001", "10001", s)
		for _, i := range cps {
			if i.Shipping.String == "" {
				ps += "\n" + i.Item + " | " + i.Price.String
			} else {
				ps += "\n" + i.Item + " | " + i.Price.String + " | " + i.Shipping.String
			}
			if ic == 19 {
				ps += "\n..."
				break
			}
			ic++
		}
		if len(cps) == 0 {
			replyText := "暂无相关记录"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GuildMsg,
			}
		} else {
			psc = "共搜到" + strconv.Itoa(len(cps)) + "条记录" + "\n品名 | 价格 | 备注" + ps + "\n价格源自 黄小姐的魔方店"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, psc)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: psc,
				},
				ReqType: utils.GuildMsg,
			}
		}
	}

	if public.Contains(priceSearch, "奇乐") {
		if public.StartsWith(s, "#+") && (isAdmin || isBotAdmin) {
			str4 := strings.TrimSpace(string([]byte(s)[len("#+"):]))
			str5 := strings.Split(str4, "##")
			if len(str5) != 2 {
				if strings.TrimSpace(str5[0]) == "" {
					replyText := "商品名不能为空"
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: replyText,
						},
						ReqType: utils.GuildMsg,
					}
				}

				err := database.IDBGAN("10002", "10002", str5[0])
				if err != nil {
					replyText := "删除失败"
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: replyText,
						},
						ReqType: utils.GuildMsg,
					}
				}
				replyText := "删除成功"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if strings.TrimSpace(str5[0]) == "" {
				replyText := "商品名不能为空"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			str6 := strings.Split(str5[1], "#?")
			if len(str6) != 2 {
				err := database.ItemSave("10002", "10002", null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, userId, null.NewTime(time.Now(), true))
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
				replyText := "添加成功"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)

				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			err := database.ItemSave("10002", "10002", null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), userId, null.NewTime(time.Now(), true))
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
			replyText := "添加成功"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GuildMsg,
			}
		}

		ps := ""
		psc := ""
		ic := 0
		cps, _ := database.GetItems("10002", "10002", s)
		for _, i := range cps {
			if i.Shipping.String == "" {
				ps += "\n" + i.Item + " | " + i.Price.String
			} else {
				ps += "\n" + i.Item + " | " + i.Price.String + " | " + i.Shipping.String
			}
			if ic == 19 {
				ps += "\n..."
				break
			}
			ic++
		}
		if len(cps) == 0 {
			replyText := "暂无相关记录"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)

			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: replyText,
				},
				ReqType: utils.GuildMsg,
			}
		} else {
			psc = "共搜到" + strconv.Itoa(len(cps)) + "条记录" + "\n品名 | 价格 | 备注" + ps + "\n价格源自 奇乐魔方坊"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, psc)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: psc,
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
	utils.Register("查价", &PricePlugin{})
}
