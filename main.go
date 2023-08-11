//go:generate goversioninfo
package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/plugins"
	_ "github.com/2mf8/QQBotOffical/plugins"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
)

func main() {
	InitLog()
	log.Infoln("欢迎您使用QQBotOffical")
	_, err := os.Stat("conf.toml")
	if err != nil {
		_ = os.WriteFile("conf.toml", []byte("Plugins = [\"守卫\",\"屏蔽\",\"开关\",\"复读\",\"回复\",\"频道管理\",\"赛季\",\"查价\",\"打乱\",\"学习\"]   # 插件管理\nAppId = 0   # 机器人AppId\nAccessToken = \"\"   # 机器人AccessToken\nAdmins = []   # 机器人管理员管理\nDatabaseUser = \"\"   # MSSQL数据库用户名\nDatabasePassword = \"\"   # MSSQL数据库密码\nDatabasePort = 1433   # MSSQL数据库服务端口\nDatabaseServer = \"127.0.0.1\"   # MSSQL数据库服务网址\nServerPort = 8081   # 服务端口\nScrambleServer = \"http://localhost:2014\"   # 打乱服务地址"), 0644)
		os.Exit(1)
	}
	plugin, _ := public.TbotConf()
	pluginString := fmt.Sprintf("%s", plugin.Conf)
	botLoginInfo, _ := public.BotLoginInfo()
	log.Infof("已加载插件 %s", pluginString)

	token := token.BotToken(botLoginInfo.AppId, botLoginInfo.AccessToken)
	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Printf("%+v, err:%v", ws, err)
	}
	// 监听哪类事件就需要实现哪类的 handler，定义：websocket/event_handler.go
	var rolesMap = map[string][]string{}
	// roles":[{"id":"4","name":"频道主","color":4294917938,"hoist":1,"number":1,"member_limit":1},{"id":"2","name":"超级管理员","color":4294936110,"hoist":1,"number":17,"member_limit":50},{"id":"7","name":"分组管理员","color":4283608319,"hoist":1,"number":0,"member_limit":50},{"id":"5","name":"子频道管理员","color":4288922822,"hoist":1,"number":16,"member_limit":50},{"id":"10012668","name":"直播组","color":4283249526,"hoist":0,"number":0,"member_limit":3000},{"id":"10012638","name":"魔方官方","color":4293221280,"hoist":1,"number":7,"member_limit":3000},{"id":"10012648","name":"知名选手","color":4294920704,"hoist":1,"number":6,"member_limit":3000},{"id":"10012655","name":"资深魔友","color":4290852578,"hoist":1,"number":40,"member_limit":3000},{"id":"10012214","name":"一个头衔","color":4288044306,"hoist":0,"number":18,"member_limit":3000},{"id":"10015793","name":"魔方店家","color":4279419354,"hoist":1,"number":2,"member_limit":3000},{"id":"13719410","name":"开发者","color":4285672924,"hoist":1,"number":2,"member_limit":3000},{"id":"13818102","name":"赛季巡查员","color":4292095291,"hoist":1,"number":2,"member_limit":3000},{"id":"13818124","name":"广告巡查员","color":4289887999,"hoist":1,"number":7,"member_limit":3000},{"id":"14102869","name":"热心魔友","color":4279419354,"hoist":1,"number":4,"member_limit":3000},{"id":"6","name":"访客","color":4286151052,"hoist":0,"number":0,"member_limit":3000},{"id":"1","name":"普通成员","color":4286151052,"hoist":0,"number":0,"member_limit":1000}],"role_num_limit":"32"}

	var message event.MessageEventHandler = func(event *dto.WSPayload, data *dto.WSMessageData) error {
		me, _ := api.Me(ctx)
		atBot := fmt.Sprintf("<@!%s>", me.ID)
		imgStr := ""
		imgs := []string{}
		if len(data.Attachments) != 0 {
			for _, imgUrl := range data.Attachments {
				imgStr += "<img:\"" + imgUrl.URL + "\">"
				imgs = append(imgs, imgUrl.URL)
			}
		}
		guildId := data.GuildID               // 频道Id
		channelId := data.ChannelID           // 子频道Id
		userId := data.Author.ID              // 用户Id
		msg := data.Content                   // 消息内容
		msgId := data.ID                      // 消息Id
		username := data.Author.Username      // 消息发送者频道昵称
		avatar := data.Author.Avatar          // 消息发送者频道头像
		isBot := data.Author.Bot              // 消息发送者是否是机器人
		srcGuildID := data.SrcGuildID         // 私信下确定频道来源
		isDirectMessage := data.DirectMessage // 是否是私信
		roles := data.Member.Roles
		br, _ := api.GuildMember(ctx, guildId, me.ID)
		botIsAdmin := public.IsAdmin(br.Roles)
		isBotAdmin := public.IsBotAdmin(userId)
		rawMsg := strings.TrimSpace(strings.ReplaceAll(msg, atBot, ""))
		reg1 := regexp.MustCompile("％")
		reg2 := regexp.MustCompile("＃")
		reg3 := regexp.MustCompile("？")
		reg4 := regexp.MustCompile("/")
		// ．
		reg5 := regexp.MustCompile("．")
		rawMsg = strings.TrimSpace(reg1.ReplaceAllString(rawMsg, "%"))
		rawMsg = strings.TrimSpace(reg2.ReplaceAllString(rawMsg, "#"))
		rawMsg = strings.TrimSpace(reg3.ReplaceAllString(rawMsg, "?"))

		fmt.Println("通过？", plugins.Pass(roles))

		if len(rolesMap[guildId]) == 0 {
			var gRoles []string
			guildRoles, err := api.Roles(ctx, guildId)
			if err == nil {
				for _, r := range guildRoles.Roles {
					if r.Name != "普通成员" && r.Name != "访客" {
						gRoles = append(gRoles, string(r.ID))
					}
				}
			}
			rolesMap[guildId] = gRoles
			fmt.Println(rolesMap)
		}

		if public.IsAdmin(roles) && strings.TrimSpace(rawMsg) == "更新" {
			var gRoles []string
			guildRoles, err := api.Roles(ctx, guildId)
			defer api.Roles(ctx, guildId)
			if err == nil {
				for _, r := range guildRoles.Roles {
					if r.Name != "普通成员" && r.Name != "访客" {
						gRoles = append(gRoles, string(r.ID))
					}
				}
			}
			rolesMap[guildId] = gRoles
			fmt.Println(rolesMap)
			api.PostMessage(ctx, channelId, &dto.MessageToCreate{Content: "更新成功"})
		}

		if isBotAdmin && public.StartsWith(rawMsg, "信") {
			dmsg, err := api.CreateDirectMessage(ctx, &dto.DirectMessageToCreate{
				SourceGuildID: guildId,
				RecipientID:   data.Author.ID,
			})
			if err != nil {
				log.Warnf("私信出错了，err = ", err)
				return nil
			}
			api.PostDirectMessage(ctx, dmsg, &dto.MessageToCreate{Content: "hello", MsgID: data.ID})
		}

		if !isBotAdmin {
			rawMsg = strings.TrimSpace(reg4.ReplaceAllString(rawMsg, ""))
		}
		rawMsg = strings.TrimSpace(reg5.ReplaceAllString(rawMsg, "."))

		if public.Contains(msg, atBot) {
			if !public.StartsWith(rawMsg, "%") {
				rawMsg = "." + rawMsg
			}
		}

		if isBotAdmin && public.StartsWith(rawMsg, "转") {
			rawMsg, _ = public.Prefix(rawMsg, "转")
			if rawMsg != "" {
				if public.Contains(rawMsg, ".") {
					api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "成功，已确认机器人管理身份"})
				} else {
					api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: rawMsg})
				}
			}
		}

		if isBotAdmin && public.StartsWith(rawMsg, "频道") {
			api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "https://pd.qq.com/s/4syyazec6"})
		}

		channel, err := api.Channel(ctx, channelId)
		if err != nil {
			log.Warnf("获取子频道信息出错， err = %+v", err)
			return nil
		}
		priceSearch := channel.Name

		if imgStr == "" {
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) <- %s", guildId, channelId, userId, rawMsg)
		} else {
			log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) <- %s %s", guildId, channelId, userId, rawMsg, imgStr)
		}
		ctx := context.WithValue(context.Background(), botLoginInfo, plugin)
		sg, _ := database.SGBGIACI(guildId, channelId)
		for _, i := range plugin.Conf {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(database.PluginNameToIntent(i))
			if intent == int64(database.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			retStuct := utils.PluginSet[i].Do(&ctx, rolesMap, guildId, channelId, userId, rawMsg, msgId, username, avatar, srcGuildID, roles, isBot, isDirectMessage, botIsAdmin, priceSearch, imgs)
			if retStuct.RetVal == utils.MESSAGE_BLOCK {
				if retStuct.ReqType == utils.GuildMsg {
					if retStuct.ReplyMsg != nil {
						newMsg := &dto.MessageToCreate{
							Content: retStuct.ReplyMsg.Text,
							MsgID:   data.ID,
						}
						if retStuct.ReplyMsg.Image != "" {
							newMsg = &dto.MessageToCreate{
								Content: retStuct.ReplyMsg.Text,
								Image:   retStuct.ReplyMsg.Image,
								MsgID:   data.ID,
							}
						}
						if len(retStuct.ReplyMsg.Images) != 0 {
							newMsg = &dto.MessageToCreate{
								Content: retStuct.ReplyMsg.Text,
								Image:   "https://" + retStuct.ReplyMsg.Images[0],
								MsgID:   data.ID,
							}
						}
						if public.Contains(retStuct.ReplyMsg.Text, "奇乐") {
							api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Ark: &dto.Ark{TemplateID: 23, KV: []*dto.ArkKV{
								{
									Key:   "#DESC#",
									Value: "消息",
								},
								{
									Key:   "#PROMPT#",
									Value: "查价",
								},
								{
									Key: "#LIST#",
									Obj: []*dto.ArkObj{
										{
											ObjKV: []*dto.ArkObjKV{
												{
													Key:   "desc",
													Value: newMsg.Content,
												},
											},
										},
										{
											ObjKV: []*dto.ArkObjKV{
												{
													Key:   "desc",
													Value: "🔗奇乐最新价格",
												},
												{
													Key:   "link",
													Value: "https://2mf8.cn/webview/#/pages/index/webview?url=https%3A%2F%2Fqilecube.gitee.io%2F",
												},
											},
										},
									},
								},
							}}})
						} else if len(retStuct.ReplyMsg.Images) == 0 && retStuct.ReplyMsg.Image == "" {
							var results [][2]string
							//s := "测试1[ss](https://2mf8.cn)test1[百度](https://www.baidu.com)jkdhi是"
							//reg := regexp.MustCompile(`(\[[^x00-xff]+\])(\([a-zA-Z0-9:/.]*\))`)
							//reg1 := regexp.MustCompile(`[一-龥a-zA-Z]+`)

							reg := regexp.MustCompile(`(\[[一-龥a-zA-Z]+\])(\([a-zA-Z0-9:/.]*\))`)
							strs := reg.FindAllString(retStuct.ReplyMsg.Text, -1)
							texts := reg.Split(retStuct.ReplyMsg.Text, -1)
							if len(strs) == 0 {
								results = append(results, [2]string{retStuct.ReplyMsg.Text})
							}
							for i, iv := range texts {
								for j, jv := range strs {
									if i == j {
										results = append(results, [2]string{iv})
										var result [2]string
										link := strings.Split(jv, "](")
										result[0] = strings.ReplaceAll(link[0], "[", "")
										result[1] = strings.ReplaceAll(link[1], ")", "")
										results = append(results, result)
									}
								}
								if i != 0 && i > len(strs)-1 && texts[i] != "" {
									results = append(results, [2]string{iv})
								}
							}
							fmt.Println(strs, texts, len(texts), results, retStuct.ReplyMsg.Text)

							var _msg []*dto.ArkObj
							for _, v := range results {
								if v[0] != "" && v[1] == "" {
									kv := &dto.ArkObj{
										ObjKV: []*dto.ArkObjKV{
											{
												Key:   "desc",
												Value: strings.TrimSpace(v[0]),
											},
										},
									}
									_msg = append(_msg, kv)
								}
								if strings.TrimSpace(v[0]) != "" && strings.TrimSpace(v[1]) != "" {
									_key := "🔗" + strings.TrimSpace(v[0])
									if strings.HasPrefix(v[1], "http") {
										_url := "https://2mf8.cn/webview/#/pages/index/webview?url=" + url.QueryEscape(strings.TrimSpace(v[1]))
										kv := &dto.ArkObj{
											ObjKV: []*dto.ArkObjKV{
												{
													Key:   "desc",
													Value: _key,
												},
												{
													Key:   "link",
													Value: _url,
												},
											},
										}
										_msg = append(_msg, kv)
									}
								}
							}
							api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Ark: &dto.Ark{TemplateID: 23, KV: []*dto.ArkKV{
								{
									Key:   "#DESC#",
									Value: "消息",
								},
								{
									Key:   "#PROMPT#",
									Value: "问答",
								},
								{
									Key: "#LIST#",
									Obj: _msg,
								},
							}}})
						} else {
							api.PostMessage(ctx, channelId, newMsg)
						}
						if len(retStuct.ReplyMsg.Images) == 2 {
							api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Image: "https://" + retStuct.ReplyMsg.Images[1]})
						}
						if len(retStuct.ReplyMsg.Images) >= 3 {
							api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Image: "https://" + retStuct.ReplyMsg.Images[1]})
							api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Image: "https://" + retStuct.ReplyMsg.Images[2]})
						}
					}
					break
				}
				if retStuct.ReqType == utils.GuildBan {
					if len(retStuct.BanId) == 0 {
						if retStuct.ReplyMsg != nil {
							newMsg := &dto.MessageToCreate{
								Content: retStuct.ReplyMsg.Text,
								MsgID:   data.ID,
							}
							api.PostMessage(ctx, channelId, newMsg)
							break
						}
						break
					} else if len(retStuct.BanId) == 1 {
						err := api.MemberMute(ctx, guildId, retStuct.BanId[0], &dto.UpdateGuildMute{MuteSeconds: retStuct.Duration})
						if err != nil {
							reply := "禁言用户<@!" + retStuct.BanId[0] + ">出错，请确认被禁言身份后重试"
							log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
							newMsg := &dto.MessageToCreate{
								Content: reply,
								MsgID:   data.ID,
							}
							api.PostMessage(ctx, channelId, newMsg)
							break
						}
						if retStuct.ReplyMsg != nil {
							reply := "已禁言用户<@!" + retStuct.BanId[0] + ">" + retStuct.ReplyMsg.Text
							log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
							newMsg := &dto.MessageToCreate{
								Content: reply,
								MsgID:   data.ID,
							}
							api.PostMessage(ctx, channelId, newMsg)
							break
						}
						break
					} else {
						_, err := api.MultiMemberMute(ctx, guildId, &dto.UpdateGuildMute{
							MuteSeconds: retStuct.Duration,
							UserIDs:     retStuct.BanId,
						})
						if err != nil {
							reply := "批量禁言用户出错，请确认被禁言用户身份后重试"
							log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
							newMsg := &dto.MessageToCreate{
								Content: reply,
								MsgID:   data.ID,
							}
							api.PostMessage(ctx, channelId, newMsg)
							break
						}
						if retStuct.ReplyMsg != nil {
							reply := "已批量禁言用户" + retStuct.ReplyMsg.Text
							log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
							newMsg := &dto.MessageToCreate{
								Content: reply,
								MsgID:   data.ID,
							}
							api.PostMessage(ctx, channelId, newMsg)
							break
						}
						break
					}
				}
				if retStuct.ReqType == utils.RelieveBan {
					if len(retStuct.BanId) == 0 {
						if retStuct.ReplyMsg != nil {
							newMsg := &dto.MessageToCreate{
								Content: retStuct.ReplyMsg.Text,
								MsgID:   data.ID,
							}
							api.PostMessage(ctx, channelId, newMsg)
							break
						}
						break
					} else if len(retStuct.BanId) == 1 {
						err := api.MemberMute(ctx, guildId, retStuct.BanId[0], &dto.UpdateGuildMute{MuteSeconds: retStuct.Duration})
						if err != nil {
							reply := "解除禁言用户<@!" + retStuct.BanId[0] + ">出错"
							log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
							newMsg := &dto.MessageToCreate{
								Content: reply,
								MsgID:   data.ID,
							}
							api.PostMessage(ctx, channelId, newMsg)
							break
						}
						reply := "已解除禁言用户<@!" + retStuct.BanId[0] + ">的禁言"
						log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
						newMsg := &dto.MessageToCreate{
							Content: reply,
							MsgID:   data.ID,
						}
						api.PostMessage(ctx, channelId, newMsg)
						break
					} else {
						_, err := api.MultiMemberMute(ctx, guildId, &dto.UpdateGuildMute{
							MuteSeconds: retStuct.Duration,
							UserIDs:     retStuct.BanId,
						})
						if err != nil {
							reply := "批量解除禁言用户出错"
							log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
							newMsg := &dto.MessageToCreate{
								Content: reply,
								MsgID:   data.ID,
							}
							api.PostMessage(ctx, channelId, newMsg)
							break
						}
						reply := "已批量解除禁言用户的禁言"
						log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
						newMsg := &dto.MessageToCreate{
							Content: reply,
							MsgID:   data.ID,
						}
						api.PostMessage(ctx, channelId, newMsg)
						break
					}
				}
				if retStuct.ReqType == utils.GuildKick {
					tMsg := ""
					if len(retStuct.BanId) == 0 {
						if retStuct.ReplyMsg != nil {
							newMsg := &dto.MessageToCreate{
								Content: retStuct.ReplyMsg.Text,
								MsgID:   data.ID,
							}
							api.PostMessage(ctx, channelId, newMsg)
							break
						}
						break
					}
					for _, ban_id := range retStuct.BanId {
						err = api.DeleteGuildMember(ctx, guildId, ban_id, dto.WithAddBlackList(retStuct.RejectAddAgain), dto.WithDeleteHistoryMsg(retStuct.Retract))
						if err != nil {
							reply := "移除用户<@!" + ban_id + ">出错，请确认被移除用户身份后重试"
							log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, reply)
							newMsg := &dto.MessageToCreate{
								Content: reply,
								MsgID:   data.ID,
							}
							api.PostMessage(ctx, channelId, newMsg)
							break
						}
					}
					if len(retStuct.BanId) == 1 {
						if retStuct.RejectAddAgain {
							tMsg = "移除用户成功，并已加入黑名单"
						} else {
							tMsg = "移除用户成功"
						}
					} else {
						if retStuct.RejectAddAgain {
							tMsg = "批量移除用户成功，并已加入黑名单"
						} else {
							tMsg = "批量移除用户成功"
						}
					}
					log.Infof("GuildId(%s) ChannelId(%s) UserId(%s) -> %s", guildId, channelId, userId, tMsg)
					newMsg := &dto.MessageToCreate{
						Content: tMsg,
						MsgID:   data.ID,
					}
					api.PostMessage(ctx, channelId, newMsg)
					break
				}
				if retStuct.ReqType == utils.GuildLeave {
				}
				if retStuct.ReqType == utils.DeleteMsg {
					api.MemberMute(ctx, guildId, userId, &dto.UpdateGuildMute{MuteSeconds: retStuct.Duration})
					newMsg := &dto.MessageToCreate{
						Content: retStuct.ReplyMsg.Text,
						MsgID:   data.ID,
					}
					api.PostMessage(ctx, channelId, newMsg)
					api.RetractMessage(ctx, channelId, msgId, openapi.RetractMessageOptionHidetip)
					break
				}
			}
		}
		return nil

	}

	var dm event.DirectMessageEventHandler = func(event *dto.WSPayload, data *dto.WSDirectMessageData) error {
		fmt.Println(event.Data)
		fmt.Println(data.GuildID, data.ChannelID, data.Content, data.ID)
		fmt.Println(data.Author.ID, data.Author.Username)
		fmt.Println(data.SrcGuildID)
		gs, _ := api.MeGuilds(ctx, &dto.GuildPager{})
		member, err := api.GuildMember(ctx, gs[1].ID, data.Author.ID)
		if err != nil || member == nil {
			log.Warnf("%s(%s) 不是频道 %s(%s) 成员，无法发送私信消息", data.Author.Username, data.Author.ID, gs[1].Name, gs[1].ID)
			return nil
		}
		log.Infof("GuildId(%s) UserId(%s)：%s", gs[1].ID, data.Author.ID, data.Content)
		dmsg, err := api.CreateDirectMessage(ctx, &dto.DirectMessageToCreate{
			SourceGuildID: gs[1].ID,
			RecipientID:   data.Author.ID,
		})
		if err != nil {
			log.Warnf("私信出错了，err = ", err)
			return nil
		}
		api.PostDirectMessage(ctx, dmsg, &dto.MessageToCreate{Content: "hello", MsgID: data.ID})
		return nil
	}

	intent := websocket.RegisterHandlers(message, dm)

	// 启动 session manager 进行 ws 连接的管理，如果接口返回需要启动多个 shard 的连接，这里也会自动启动多个
	botgo.NewSessionManager().Start(ws, token, &intent)
}

func InitLog() {
	// 输出到命令行
	customFormatter := &log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
	}
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)

	// 输出到文件
	rl, err := rotatelogs.New(path.Join("logs", "%Y-%m-%d.log"),
		rotatelogs.WithLinkName(path.Join("logs", "latest.log")), // 最新日志软链接
		rotatelogs.WithRotationTime(time.Hour*24),                // 每天一个新文件
		rotatelogs.WithMaxAge(time.Hour*24*3),                    // 日志保留3天
	)
	if err != nil {
		utils.FatalError(err)
		return
	}
	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			log.InfoLevel:  rl,
			log.WarnLevel:  rl,
			log.ErrorLevel: rl,
			log.FatalLevel: rl,
			log.PanicLevel: rl,
		},
		&easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%time%] [%lvl%]: %msg% \r\n",
		},
	))
}
