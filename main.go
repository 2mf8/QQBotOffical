//go:generate goversioninfo
package main

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	database "github.com/2mf8/QQBotOffical/data"
	_dto "github.com/2mf8/QQBotOffical/dto"
	_ "github.com/2mf8/QQBotOffical/plugins"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/router"
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
	"google.golang.org/protobuf/proto"
	"gopkg.in/guregu/null.v3"
)

func main() {
	InitLog()
	go GinRun()

	fmt.Println(public.RandomString(6))
	//go database.GetAll()

	tomlData := `
	Plugins = ["守卫","开关","复读","服务号","WCA","回复","频道管理","赛季","查价","打乱","学习"]   # 插件管理
	AppId = 0 # 机器人AppId
	AccessToken = "" # 机器人AccessToken
	Admins = [""]   # 机器人管理员管理
	DatabaseUser = "sa"   # MSSQL数据库用户名
	DatabasePassword = ""   # MSSQL数据库密码
	DatabasePort = 1433   # MSSQL数据库服务端口
	DatabaseServer = "127.0.0.1"   # MSSQL数据库服务网址
	DatabaseName = ""  # 数据库名
	ServerPort = 8081   # 服务端口
	ScrambleServer = "http://localhost:2014"   # 打乱服务地址
	RedisServer = "127.0.0.1" # Redis服务网址
	RedisPort = 6379 # Redis端口
	RedisPassword = "" # Redis密码
	RedisTable = 0 # Redis数据表
	RedisPoolSize = 1000 # Redis连接池数量
	JwtKey = ""
	RefreshKey = ""
	`

	log.Infoln("欢迎您使用QQBotOffical")
	_, err := os.Stat("conf.toml")
	if err != nil {
		_ = os.WriteFile("conf.toml", []byte(tomlData), 0644)
		log.Warn("已生成配置文件 conf.toml ,请修改后重新启动程序。")
		log.Info("该程序将于5秒后退出！")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}
	allconfig := database.AllConfig
	log.Info("[配置信息]", allconfig)
	pluginString := fmt.Sprintf("%s", allconfig.Plugins)
	botLoginInfo := &public.BotLogin{
		AppId:       allconfig.AppId,
		AccessToken: allconfig.AccessToken,
	}
	log.Infof("已加载插件 %s", pluginString)

	token := token.BotToken(botLoginInfo.AppId, botLoginInfo.AccessToken)
	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Warn("登录失败，请检查 appid 和 AccessToken 是否正确。")
		log.Info("该程序将于5秒后退出！")
		time.Sleep(time.Second * 5)
		log.Printf("%+v, err:%v", ws, err)
	}
	// 监听哪类事件就需要实现哪类的 handler，定义：websocket/event_handler.go
	var rolesMap = map[string][]string{}
	// roles":[{"id":"4","name":"频道主","color":4294917938,"hoist":1,"number":1,"member_limit":1},{"id":"2","name":"超级管理员","color":4294936110,"hoist":1,"number":17,"member_limit":50},{"id":"7","name":"分组管理员","color":4283608319,"hoist":1,"number":0,"member_limit":50},{"id":"5","name":"子频道管理员","color":4288922822,"hoist":1,"number":16,"member_limit":50},{"id":"10012668","name":"直播组","color":4283249526,"hoist":0,"number":0,"member_limit":3000},{"id":"10012638","name":"魔方官方","color":4293221280,"hoist":1,"number":7,"member_limit":3000},{"id":"10012648","name":"知名选手","color":4294920704,"hoist":1,"number":6,"member_limit":3000},{"id":"10012655","name":"资深魔友","color":4290852578,"hoist":1,"number":40,"member_limit":3000},{"id":"10012214","name":"一个头衔","color":4288044306,"hoist":0,"number":18,"member_limit":3000},{"id":"10015793","name":"魔方店家","color":4279419354,"hoist":1,"number":2,"member_limit":3000},{"id":"13719410","name":"开发者","color":4285672924,"hoist":1,"number":2,"member_limit":3000},{"id":"13818102","name":"赛季巡查员","color":4292095291,"hoist":1,"number":2,"member_limit":3000},{"id":"13818124","name":"广告巡查员","color":4289887999,"hoist":1,"number":7,"member_limit":3000},{"id":"14102869","name":"热心魔友","color":4279419354,"hoist":1,"number":4,"member_limit":3000},{"id":"6","name":"访客","color":4286151052,"hoist":0,"number":0,"member_limit":3000},{"id":"1","name":"普通成员","color":4286151052,"hoist":0,"number":0,"member_limit":1000}],"role_num_limit":"32"}

	var message event.MessageEventHandler = func(event *dto.WSPayload, data *dto.WSMessageData) error {

		/*gss, _ := api.MeGuilds(ctx, &dto.GuildPager{})
		for _, v := range gss {
			fmt.Println(v.OwnerID, v.Name)
		}*/

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
		isBotAdmin := public.IsBotAdmin(userId, allconfig.Admins)
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

		u, b := public.Prefix(rawMsg, ".创建账号")
		if b {
			role := 0
			reg11 := regexp.MustCompile("@!")
			reg12 := regexp.MustCompile("@")
			reg13 := regexp.MustCompile(">")
			reg14 := regexp.MustCompile("  ")

			str1 := strings.TrimSpace(reg11.ReplaceAllString(u, "at qq=\""))
			str1 = strings.TrimSpace(reg12.ReplaceAllString(str1, "at qq=\""))
			str2 := strings.TrimSpace(reg13.ReplaceAllString(str1, "\"/>"))

			for public.Contains(str2, "  ") {
				str2 = strings.TrimSpace(reg14.ReplaceAllString(str2, " "))
			}
			t, cs := public.GuildAtConvert(str2)
			if isBotAdmin {
				fmt.Println(t, cs)
				if strings.TrimSpace(t) == "10000" {
					role = 1 << 30
				} else {
					sng, _ := database.ServerNumbersGet()
					role = 1 << sng.ServerNumberSetSync.Intent[sng.ServerNumberSetSync.ServerNumbers[strings.TrimSpace(t)]]
				}
				for _, _ui := range cs {
					_u, err := api.GuildMember(ctx, guildId, _ui)
					if err != nil {
						continue
					}
					err = database.UserInfoSave(null.NewString(_ui, true), null.NewString(_u.Nick, true), null.NewString(_u.User.Avatar, true), null.NewString(strings.TrimSpace(t), true), null.NewString("", true), null.NewString("", true), null.NewString("", true), null.NewString("", true), role)
					if err != nil {
						api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "账号创建失败"})
						fmt.Println(err)
						return nil
					} else {
						api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "账号创建成功"})
						return nil
					}
				}
			}
			if userId == "1161014622077006888" {
				role = 1 << 1 // 黄小姐
				for _, _ui := range cs {
					_u, err := api.GuildMember(ctx, guildId, _ui)
					if err != nil {
						continue
					}
					err = database.UserInfoSave(null.NewString(_ui, true), null.NewString(_u.Nick, true), null.NewString(_u.User.Avatar, true), null.NewString("10001", true), null.NewString("", true), null.NewString("", true), null.NewString("", true), null.NewString("", true), role)
					if err != nil {
						api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "账号创建失败"})
						fmt.Println(err)
						return nil
					} else {
						api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "账号创建成功"})
						return nil
					}
				}
			}
			if userId == "18155629338841245002" {
				role = 1 << 1 //奇乐
				for _, _ui := range cs {
					_u, err := api.GuildMember(ctx, guildId, _ui)
					if err != nil {
						continue
					}
					err = database.UserInfoSave(null.NewString(_ui, true), null.NewString(_u.Nick, true), null.NewString(_u.User.Avatar, true), null.NewString("10002", true), null.NewString("", true), null.NewString("", true), null.NewString("", true), null.NewString("", true), role)
					if err != nil {
						api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "账号创建失败"})
						fmt.Println(err)
						return nil
					} else {
						api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "账号创建成功"})
						return nil
					}
				}
			}

			err = database.UserInfoSave(null.NewString(userId, true), null.NewString(username, true), null.NewString(avatar, true), null.NewString("", true), null.NewString("", true), null.NewString("", true), null.NewString("", true), null.NewString("", true), role)
			if err != nil {
				api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "账号创建失败"})
				fmt.Println(err)
				return nil
			} else {
				api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "账号创建成功"})
				return nil
			}
		}

		if rawMsg == ".登录" {
			randomString := public.RandomString(6)
			database.RedisSet(randomString, []byte(userId))
			//go GetT(userId)
			byteCode, _ := proto.Marshal(&_dto.CodeLoginReq{Code: randomString})
			fmt.Printf("curl -X POST -H \"Content-Type:application/x-protobuf\" -d %v http://localhost:8200/login\n", bytes.NewBuffer(byteCode))
			dmsg, err := api.CreateDirectMessage(ctx, &dto.DirectMessageToCreate{
				SourceGuildID: guildId,
				RecipientID:   data.Author.ID,
			})
			if err != nil {
				log.Warnf("私信出错了，err = %v", err)
				return nil
			}
			api.PostDirectMessage(ctx, dmsg, &dto.MessageToCreate{Content: fmt.Sprintf("登录信息\n验证码：%s\n注：该验证码五分钟内有效。", randomString), MsgID: data.ID})
			api.PostMessage(ctx, channelId, &dto.MessageToCreate{Content: "登录信息已私发，请查看私信。", MsgID: msgId})
		}

		if len(rolesMap[guildId]) == 0 {
			var gRoles []string
			guildRoles, err := api.Roles(ctx, guildId)
			if err == nil {
				for _, r := range guildRoles.Roles {
					if r.Name != "普通成员" && r.Name != "访客" {
						gRoles = append(gRoles, string(r.ID))
						fmt.Println(r.Name, r.ID)
					}
				}
			}
			rolesMap[guildId] = gRoles
			fmt.Println(rolesMap)
		}

		if (public.IsAdmin(roles) || public.IsBotAdmin(userId, allconfig.Admins)) && strings.TrimSpace(rawMsg) == "更新" {
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
			// fmt.Println(rolesMap)
			api.PostMessage(ctx, channelId, &dto.MessageToCreate{Content: "更新成功", MsgID: data.ID})
		}

		if isBotAdmin && public.StartsWith(rawMsg, "信") {
			dmsg, err := api.CreateDirectMessage(ctx, &dto.DirectMessageToCreate{
				SourceGuildID: guildId,
				RecipientID:   data.Author.ID,
			})
			if err != nil {
				log.Warnf("私信出错了，err = %v", err)
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
		ctx := context.WithValue(context.Background(), botLoginInfo, allconfig.Plugins)
		sg, _ := database.SGBGIACI(guildId, channelId)
		for _, i := range allconfig.Plugins {
			intent := sg.PluginSwitch.IsCloseOrGuard & int64(database.PluginNameToIntent(i))
			if intent == int64(database.PluginReply) {
				break
			}
			if intent > 0 {
				continue
			}
			retStuct := utils.PluginSet[i].Do(&ctx, allconfig.Admins, rolesMap, guildId, channelId, userId, rawMsg, msgId, username, avatar, srcGuildID, roles, isBot, isDirectMessage, botIsAdmin, priceSearch, imgs)
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
				/*if retStuct.ReqType == utils.GuildLeave {
				}*/
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
		//fmt.Println(event.Data)
		//fmt.Println(data.GuildID, data.ChannelID, data.Content, data.ID)
		//fmt.Println(data.Author.ID, data.Author.Username)
		// fmt.Println(data.SrcGuildID)
		gs, _ := api.MeGuilds(ctx, &dto.GuildPager{})

		for _, g := range gs {
			member, err := api.GuildMember(ctx, g.ID, data.Author.ID)
			if err != nil || member == nil {
				log.Warnf("%s(%s) 不是频道 %s(%s) 成员，无法发送私信消息", data.Author.Username, data.Author.ID, g.Name, g.ID)
				continue
			}
			log.Infof("GuildId(%s) UserId(%s)：%s", g.ID, data.Author.ID, data.Content)
			dmsg, err := api.CreateDirectMessage(ctx, &dto.DirectMessageToCreate{
				SourceGuildID: g.ID,
				RecipientID:   data.Author.ID,
			})
			if err != nil {
				log.Warn("私信出错了，err = ", err)
				continue
			}
			api.PostDirectMessage(ctx, dmsg, &dto.MessageToCreate{Content: "hello", MsgID: data.ID})
			break
		}
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

func GinRun() {
	defer database.Db.Close()
	r := router.InitRouter()
	r.Run(":8200")
}

func GetT(userId string) {
	for i := 10; i > 0; i-- {
		v, err := database.RedisGet(userId)
		fmt.Println(i, v, err)
		time.Sleep(time.Second)
	}
}
