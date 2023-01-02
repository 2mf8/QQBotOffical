package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	database "github.com/2mf8/QQBotOffical/data"
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
	// role 4 频道主 11 普通成员 2 管理员 5 子频道管理 13 普通成员
	// 监听哪类事件就需要实现哪类的 handler，定义：websocket/event_handler.go
	var message event.MessageEventHandler = func(event *dto.WSPayload, data *dto.WSMessageData) error {
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
		isAdmin := public.IsAdmin(roles)
		br, _ := api.GuildMember(ctx, guildId, "13970278473675774808")
		botIsAdmin := public.IsAdmin(br.Roles)
		isBotAdmin := public.IsBotAdmin(userId)
		rawMsg := strings.TrimSpace(strings.ReplaceAll(msg, "<@!13970278473675774808>", ""))
		reg1 := regexp.MustCompile("％")
		reg2 := regexp.MustCompile("＃")
		reg3 := regexp.MustCompile("？")
		reg4 := regexp.MustCompile("/")
		rawMsg = strings.TrimSpace(reg1.ReplaceAllString(rawMsg, "%"))
		rawMsg = strings.TrimSpace(reg2.ReplaceAllString(rawMsg, "#"))
		rawMsg = strings.TrimSpace(reg3.ReplaceAllString(rawMsg, "?"))
		rawMsg = strings.TrimSpace(reg4.ReplaceAllString(rawMsg, ""))

		if public.Contains(msg, "<@!13970278473675774808>") {
			if !public.StartsWith(rawMsg, "%") {
				rawMsg = "." + rawMsg
			}
		}
		fmt.Println("管理员？", isAdmin, public.IsAdmin(br.Roles)) // 消息发送者角色
		if msg == "机器人权限确认" {
			if public.IsAdmin(br.Roles) {
				api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "机器人是管理员"})
			} else {
				api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Content: "机器人不是管理员"})
			}
			//api.DeleteGuildMember()
		}
		if msg == "url" {
			api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: msgId, Image: "https://gchat.qpic.cn/qmeetpic/14205791670163322/32310159-2985743808-30191304C37B8F47654B38834397C2BF/0/"})
			//api.DeleteGuildMember()
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
			retStuct := utils.PluginSet[i].Do(&ctx, guildId, channelId, userId, rawMsg, msgId, username, avatar, srcGuildID, isBot, isDirectMessage, botIsAdmin, isBotAdmin, isAdmin, priceSearch, imgs)
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
						api.PostMessage(ctx, channelId, newMsg)
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
					api.MemberMute(ctx, guildId, userId, &dto.UpdateGuildMute{MuteSeconds: "120"})
					newMsg := &dto.MessageToCreate{
						Content: retStuct.ReplyMsg.Text,
						MsgID:   data.ID,
					}
					api.PostMessage(ctx, channelId, newMsg)
					api.RetractMessage(ctx, channelId, msgId, openapi.RetractMessageOptionHidetip, openapi.RetractMessageOptionHidetip)
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
