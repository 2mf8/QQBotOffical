package plugins

import (
	"context"
	"fmt"
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

func (price *PricePlugin) Do(ctx *context.Context, messageType public.MessageType, admins []string, gmap map[string][]string, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, useRole []string, isBot, isDirectMessage, botIsAdmin bool, priceSearch string, attachments []string) utils.RetStuct {

	isBotAdmin := public.IsBotAdmin(userId, admins)
	isAdmin := public.IsAdmin(useRole)
	reg1 := regexp.MustCompile("％")
	reg2 := regexp.MustCompile("＃")
	reg3 := regexp.MustCompile("＆")
	str1 := strings.TrimSpace(reg1.ReplaceAllString(msg, "%"))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, "#"))
	str3 := strings.TrimSpace(reg3.ReplaceAllString(str2, "&"))
	is_magnetism := false

	messageType = public.Undefined
	s, b := public.Prefix(str3, "%", messageType)
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

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

	if len(attachments) == 0 {
		attachments = append(attachments, "HJdhuhjd")
	}
	_pass := database.IsExist(guildId)
	fmt.Println("通过？", _pass)
	if _pass {
		fmt.Println("通过")
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
					replyText := "删除失败，该商品不存在。"
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
				_, err := database.ItemSave("10001", "10001", null.String{}, str5[0], null.NewString(str6[0], true), null.String{}, null.NewString(userId, true), time.Now().Unix(), is_magnetism, null.String{})
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
			_, err := database.ItemSave("10001", "10001", null.String{}, str5[0], null.NewString(str6[0], true), null.NewString(str6[1], true), null.NewString(userId, true), time.Now().Unix(), is_magnetism, null.String{})
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
		ti := 0
		topStart := 0
		ss := strings.Split(s, "#")
		if len(ss) > 1 {
			topStart, _ = strconv.Atoi(strings.TrimSpace(ss[1]))
		}
		if topStart == 0 {
			topStart = 1
		}
		cps, _ := database.GetItems("10001", "10001", ss[0])
		for ii, i := range cps {
			if !(ii < topStart-1 || ii > topStart+18) {
				ic++
				if i.Shipping.String == "" {
					ps += "\n" + i.Item + " | " + i.Price.String
					ti++
				} else {
					ps += "\n" + i.Item + " | " + i.Price.String + " | " + i.Shipping.String
					ti++
				}
			} else if ii > topStart+18 {
				ic++
				break
			}
		}
		if ic > 20 {
			ps += "\n...\n\n翻页请使用\n%[品名]#[序号]\n指令。例如：\n%三#21\n"
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
			if ti == 0 {
				replyText := fmt.Sprintf("%s Top %d - %d ", ss[0], topStart, topStart+19) + "暂无相关记录"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, replyText)

				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: replyText,
					},
					ReqType: utils.GuildMsg,
				}
			}
			psc = "共搜到" + strconv.Itoa(len(cps)) + "条记录" + fmt.Sprintf("\nTop %d - %d", topStart, topStart+ti-1) + "\n品名 | 价格 | 备注" + ps + "\n价格源自 奇乐魔方坊"
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
