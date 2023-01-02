package plugins

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
	log "github.com/sirupsen/logrus"
)

type Competition struct {
}

func (rep *Competition) Do(ctx *context.Context, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, isBot, isDirectMessage, botIsAdmin, isBotAdmin, isAdmin bool, priceSearch string, imgs []string) utils.RetStuct {
	var sic []string

	s, b := public.Prefix(msg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}

	sc, b := public.Prefix(s, "新赛季")
	if b && isBotAdmin {
		if sc == "" {
			reply := "格式错误"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		si := strings.Split(strings.TrimSpace(sc), " ")
		d, err := strconv.Atoi(si[0])
		if err != nil || d > 30 || d < 1 {
			d = 30
		}

		for _, v := range si {
			switch v {
			case "222":
				sic = append(sic, v)
			case "333":
				sic = append(sic, v)
			case "444":
				sic = append(sic, v)
			case "555":
				sic = append(sic, v)
			case "666":
				sic = append(sic, v)
			case "777":
				sic = append(sic, v)
			case "skewb":
				sic = append(sic, v)
			case "pyram":
				sic = append(sic, v)
			case "sq1":
				sic = append(sic, v)
			case "clock":
				sic = append(sic, v)
			case "minx":
				sic = append(sic, v)
			case "all":
				sic = append(sic, []string{"222", "333", "444", "555", "666", "777", "skewb", "pyram", "sq1", "clock", "minx"}...)
			}
		}

		if len(sic) == 0 {
			reply := "赛季项目不能为空"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}

		cr, _ := database.CompetitionRead()
		if time.Now().Unix() < cr.EndTime && time.Now().Unix() > cr.StartTime {
			reply := "已存在赛季" + strconv.Itoa(cr.Sessions) + ",请等待赛季结束后再开启新赛季"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		} else {
			cr.Sessions += 1
		}

		err = cr.CompetitionCreate(d, sic)
		if err != nil {
			reply := "创建失败"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		reply := "创建成功"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: reply,
			},
			ReqType: utils.GuildMsg,
		}
	}
	sczj, b := public.Prefix(s, "赛季追加")
	if b && isBotAdmin {
		if sczj == "" {
			reply := "格式错误"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		si := strings.Split(strings.TrimSpace(sc), " ")

		for _, v := range si {
			switch v {
			case "222":
				sic = append(sic, v)
			case "333":
				sic = append(sic, v)
			case "444":
				sic = append(sic, v)
			case "555":
				sic = append(sic, v)
			case "666":
				sic = append(sic, v)
			case "777":
				sic = append(sic, v)
			case "skewb":
				sic = append(sic, v)
			case "pyram":
				sic = append(sic, v)
			case "sq1":
				sic = append(sic, v)
			case "clock":
				sic = append(sic, v)
			case "minx":
				sic = append(sic, v)
			}
		}

		if len(sic) == 0 {
			reply := "赛季项目不能为空"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}

		cr, _ := database.CompetitionRead()
		tip, err := cr.CompetitionUpdate(sic)
		if tip != "" {
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, tip)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: tip,
				},
				ReqType: utils.GuildMsg,
			}
		}
		if err != nil {
			reply := "追加失败"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		reply := "追加成功"
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: reply,
			},
			ReqType: utils.GuildMsg,
		}
	}

	if s == "赛季信息" {
		v, _ := database.CompetitionRead()
		session := v.Sessions
		startTime := time.Unix(v.StartTime, 0)
		endTime := time.Unix(v.EndTime, 0)
		items := strings.Join(v.Items, "、")
		reply := "赛季信息\n场次：" + strconv.Itoa(session) + "\n开始时间：" + startTime.Format("2006-01-02 15:04:05") + "\n结束时间：" + endTime.Format("2006-01-02 15:04:05") + "\n赛季项目：" + items
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: reply,
			},
			ReqType: utils.GuildMsg,
		}
	}

	scr, b := public.Prefix(s, "赛季打乱")
	if b {
		var si []string
		if strings.TrimSpace(scr) == "" {
			reply := "获取出错，格式不对"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		gss, _ := database.CompetitionRead()
		scrs := strings.Split(strings.TrimSpace(scr), " ")
		for _, v := range scrs {
			if v != "" {
				si = append(si, v)
			}
		}
		tgc := database.ToGetScramble(si[0])
		if tgc != "" {
			if len(si) < 2 {
				si = append(si, "-1")
			}
			t := database.ToGetScrambleIndex(si[1])
			if t == 0 {
				if (tgc == "444" && gss.CompContents.Four != "") || (tgc == "555" && gss.CompContents.Five != "") || (tgc == "666" && gss.CompContents.Six != "") || (tgc == "777" && gss.CompContents.Seven != "") || (tgc == "minx" && gss.CompContents.Megaminx != "") {
					reply := "公式较长，请分批获取\n赛季打乱 [项目] [序号]\n注：序号为1-5"
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: reply,
						},
						ReqType: utils.GuildMsg,
					}
				}
				if tgc == "333" && gss.CompContents.Three != "" {
					tsc := strings.Split(gss.CompContents.Three, "\n")
					reply := "3阶\n1、" + tsc[0] + "\n2、" + tsc[1] + "\n3、" + tsc[2] + "\n4、" + tsc[3] + "\n5、" + tsc[4]
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: reply,
						},
						ReqType: utils.GuildMsg,
					}
				}
				if tgc == "222" && gss.CompContents.Two != "" {
					tsc := strings.Split(gss.CompContents.Two, "\n")
					reply := "2阶\n1、" + tsc[0] + "\n2、" + tsc[1] + "\n3、" + tsc[2] + "\n4、" + tsc[3] + "\n5、" + tsc[4]
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: reply,
						},
						ReqType: utils.GuildMsg,
					}
				}
				if tgc == "skewb" && gss.CompContents.Skewb != "" {
					tsc := strings.Split(gss.CompContents.Skewb, "\n")
					reply := "Skewb\n1、" + tsc[0] + "\n2、" + tsc[1] + "\n3、" + tsc[2] + "\n4、" + tsc[3] + "\n5、" + tsc[4]
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: reply,
						},
						ReqType: utils.GuildMsg,
					}
				}
				if tgc == "sq1" && gss.CompContents.Square != "" {
					tsc := strings.Split(gss.CompContents.Square, "\n")
					reply := "Sq1\n1、" + tsc[0] + "\n2、" + tsc[1] + "\n3、" + tsc[2] + "\n4、" + tsc[3] + "\n5、" + tsc[4]
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: reply,
						},
						ReqType: utils.GuildMsg,
					}
				}
				if tgc == "pyram" && gss.CompContents.Pyraminx != "" {
					tsc := strings.Split(gss.CompContents.Pyraminx, "\n")
					reply := "Pyram\n1、" + tsc[0] + "\n2、" + tsc[1] + "\n3、" + tsc[2] + "\n4、" + tsc[3] + "\n5、" + tsc[4]
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: reply,
						},
						ReqType: utils.GuildMsg,
					}
				}
				if tgc == "clock" && gss.CompContents.Clock != "" {
					tsc := strings.Split(gss.CompContents.Clock, "\n")
					reply := "Clock\n1、" + tsc[0] + "\n2、" + tsc[1] + "\n3、" + tsc[2] + "\n4、" + tsc[3] + "\n5、" + tsc[4]
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
					return utils.RetStuct{
						RetVal: utils.MESSAGE_BLOCK,
						ReplyMsg: &utils.Msg{
							Text: reply,
						},
						ReqType: utils.GuildMsg,
					}
				}
				reply := "项目不存在，请使用赛季追加功能追加"
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: reply,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "222" && gss.CompContents.Two != "" {
				tsc := strings.Split(gss.CompContents.Two, "\n")
				reply := "2阶\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "333" && gss.CompContents.Three != "" {
				tsc := strings.Split(gss.CompContents.Three, "\n")
				reply := "3阶\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "444" && gss.CompContents.Four != "" {
				tsc := strings.Split(gss.CompContents.Four, "\n")
				reply := "4阶\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "555" && gss.CompContents.Five != "" {
				tsc := strings.Split(gss.CompContents.Five, "\n")
				reply := "5阶\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "666" && gss.CompContents.Six != "" {
				tsc := strings.Split(gss.CompContents.Six, "\n")
				reply := "6阶\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "777" && gss.CompContents.Seven != "" {
				tsc := strings.Split(gss.CompContents.Seven, "\n")
				reply := "7阶\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "skewb" && gss.CompContents.Skewb != "" {
				tsc := strings.Split(gss.CompContents.Skewb, "\n")
				reply := "Skewb\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "pyram" && gss.CompContents.Pyraminx != "" {
				tsc := strings.Split(gss.CompContents.Pyraminx, "\n")
				reply := "Pyraminx\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "sq1" && gss.CompContents.Square != "" {
				tsc := strings.Split(gss.CompContents.Square, "\n")
				reply := "Square One\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "clock" && gss.CompContents.Clock != "" {
				tsc := strings.Split(gss.CompContents.Clock, "\n")
				reply := "Clock\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			if tgc == "minx" && gss.CompContents.Megaminx != "" {
				tsc := strings.Split(gss.CompContents.Megaminx, "\n")
				tsc[t-1] = strings.Replace(tsc[t-1], "U' ", "#\n", -1)
				tsc[t-1] = strings.Replace(tsc[t-1], "U ", "U\n", -1)
				tsc[t-1] = strings.Replace(tsc[t-1], "#", "U'", -1)
				reply := "Megaminx\n" + strconv.Itoa(t) + "、" + tsc[t-1]
				imgUrl := "http://2mf8.cn:2014/view/" + tgc + ".png?scramble=" + url.QueryEscape(strings.Replace(tsc[t-1], "\n", " ", -1))
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text:  reply,
						Image: imgUrl,
					},
					ReqType: utils.GuildMsg,
				}
			}
			reply := "项目不存在，请使用赛季追加功能追加"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		} else {
			reply := "获取出错，格式不对"
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

	sjcj, b := public.Prefix(s, "赛季成绩")
	if b {
		cji := strings.Split(sjcj, " ")
		v, err := database.CompetitionRead()
		if err != nil {
			reply := "赛季获取出错，请联系管理员添加赛季"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		session := v.Sessions
		tgc := database.ToGetScramble(cji[0])
		if tgc != "" {
			best := -1
			average := -1
			sum := 0
			gs := public.ConvertGradeToInt(sjcj)
			if len(gs) < 4 && len(gs) > 0 {
				minNum := gs[0]
				if len(gs) == 3 && (tgc == "666" || tgc == "777") {
					for _, i := range gs {
						if i < minNum {
							minNum = i
						}
						sum += i
					}
					average = sum / 3
					best = minNum
				} else {
					for _, i := range gs {
						if i < minNum {
							minNum = i
						}
					}
					best = minNum
				}
			}
			if len(gs) == 4 {
				minNum := gs[0]
				for _, i := range gs {
					if i < minNum {
						minNum = i
					}
					sum += i
				}
				best = minNum
				average = (sum - best) / 3
			}
			if len(gs) == 5 {
				minNum := gs[0]
				maxNum := gs[0]
				for _, i := range gs {
					if i < minNum {
						minNum = i
					}
					if i > maxNum {
						maxNum = i
					}
					sum += i
				}
				best = minNum
				average = (sum - best - maxNum) / 3
			}
			err := database.AchievementSave(userId, username, tgc, best, average, session)
			if err != nil {
				reply := "成绩保存出错，请联系管理员或稍后重试"
				fmt.Println(err)
				log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
				return utils.RetStuct{
					RetVal: utils.MESSAGE_BLOCK,
					ReplyMsg: &utils.Msg{
						Text: reply,
					},
					ReqType: utils.GuildMsg,
				}
			}
			bt := "DNF"
			at := "DNF"
			bm := best / 60000
			bs := best % 60000 / 1000
			bms := best % 60000 % 1000
			am := average / 60000
			as := average % 60000 / 1000
			ams := average % 60000 % 1000
			if best > -1 && bm == 0 {
				bt = fmt.Sprintf("%d.%d", bs, bms)
			}
			if bm > 0 {
				bt = fmt.Sprintf("%d:%d.%d", bm, bs, bms)
			}
			if average > -1 && am == 0{
				at = fmt.Sprintf("%d.%d", as, ams)
			}
			if am > 0 {
				at = fmt.Sprintf("%d:%d.%d", am, as, ams)
			}
			pm := ""
			bi, ai, err := database.AchievementGetCount(tgc, best, average, session)
			fmt.Println(bi, ai)
			if err != nil {
				fmt.Println(err)
			}
			if bi >= 0 && bt != "DNF" {
				pm += fmt.Sprintf("单次 %d (%s)，", bi+1, bt)
			}
			if ai >= 0 && at != "DNF" {
				pm += fmt.Sprintf("平均 %d (%s)。", ai+1, at)
			}
			reply := "成绩上传成功。\n" + username + " (" + userId + ") 在赛季 " + strconv.Itoa(session) + "的项目 " + tgc + " 中获得排名：" + pm + "\n" + "详情(项目+成绩)：\n" + sjcj
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		} else {
			reply := "格式错误，缺少项目\n格式为\n赛季成绩 [项目] [成绩] [成绩] [成绩] ..."
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

	_, b = public.Prefix(s, "我的成绩")
	if b {
		v, err := database.CompetitionRead()
		if err != nil {
			reply := "赛季获取出错，请联系管理员添加赛季"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		session := v.Sessions
		as, err := database.AGBUAS(userId, session)
		if err != nil {
			reply := "成绩获取出错，请稍后重试"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		if len(as) == 0 {
			reply := "暂无相关成绩"
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: reply,
				},
				ReqType: utils.GuildMsg,
			}
		}
		it := fmt.Sprintf("%s(%s)在赛季%d中的成绩为\n 项目 最好成绩 || 平均成绩", username, userId, session)
		for _, v := range as {
			bt := "DNF"
			at := "DNF"
			bm := v.Best / 60000
			bs := v.Best % 60000 / 1000
			bms := v.Best % 60000 % 1000
			am := v.Average / 60000
			as := v.Average % 60000 / 1000
			ams := v.Average % 60000 % 1000
			if v.Best > -1 && bm == 0 {
				bt = fmt.Sprintf("%d.%d", bs, bms)
			}
			if bm > 0 {
				bt = fmt.Sprintf("%d:%d.%d", bm, bs, bms)
			}
			if v.Average > -1 && am == 0{
				at = fmt.Sprintf("%d.%d", as, ams)
			}
			if am > 0 {
				at = fmt.Sprintf("%d:%d.%d", am, as, ams)
			}
			it += fmt.Sprintf("\n%s %s || %s", v.Item, bt, at)
		}
		log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, it)
		return utils.RetStuct{
			RetVal: utils.MESSAGE_BLOCK,
			ReplyMsg: &utils.Msg{
				Text: it,
			},
			ReqType: utils.GuildMsg,
		}
	}

	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("赛季", &Competition{})
}
