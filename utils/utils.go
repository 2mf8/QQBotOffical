package utils

import (
	"context"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

type ReqType int

const (
	GuildBan   ReqType = iota // 频道禁言
	RelieveBan                // 禁言解除
	GuildKick                 // 频道踢人
	GuildMsg                  // 频道消息
	GuildLeave                // 退频道
	DeleteMsg                 // 消息撤回
	Undefined                 // 未定义
)

type RetStuct struct {
	RetVal         uint
	ReplyMsg       *Msg
	ReqType        ReqType
	Duration       string
	BanId          []string
	RejectAddAgain bool
	Retract        int
	MsgId          string
}

type Msg struct {
	Text  string
	At    bool
	Image string
}

type Plugin interface {
	Do(ctx *context.Context, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, isBot, isDirectMessage, botIsAdmin, isBotAdmin, isAdmin bool, priceSearch string) (retStuct RetStuct)
}

var PluginSet map[string]Plugin

const (
	MESSAGE_BLOCK  uint = 0
	MESSAGE_IGNORE uint = 1
)

func init() {
	PluginSet = make(map[string]Plugin)
}

func Register(k string, v Plugin) {
	PluginSet[k] = v
}

func FatalError(err error) {
	log.Errorf(err.Error())
	buf := make([]byte, 64<<10)
	buf = buf[:runtime.Stack(buf, false)]
	sBuf := string(buf)
	log.Errorf(sBuf)
	time.Sleep(5 * time.Second)
	panic(err)
}
