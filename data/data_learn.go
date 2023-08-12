package database

import (
	"fmt"
	"time"

	"github.com/2mf8/QQBotOffical/config"
	"github.com/2mf8/QQBotOffical/public"
	_ "github.com/denisenkom/go-mssqldb"
	"gopkg.in/guregu/null.v3"
	_ "gopkg.in/guregu/null.v3/zero"
)

type Learn struct {
	Id          int64
	Ask         string
	GuildId     string
	ChannelId   string
	AdminId     string
	Answer      null.String
	GmtModified time.Time
	Pass        bool
}

func LearnGet(guildId, channelId, ask string) (l Learn, err error) {
	l = Learn{}
	err = Db.QueryRow("select ID, ask, guild_id, channel_id, admin_id, answer, gmt_modified, pass from [$4].[dbo].[guild_learn] where guild_id = $1 and ask = $2 and channel_id = $3", guildId, ask, channelId, config.Conf.DatabaseName).Scan(&l.Id, &l.Ask, &l.GuildId, &l.ChannelId, &l.AdminId, &l.Answer, &l.GmtModified, &l.Pass)
	fmt.Println(l, err)
	info := fmt.Sprintf("%s", err)
	if public.StartsWith(info, "sql") || public.StartsWith(info, "unable") {
		if public.StartsWith(info, "unable") {
			fmt.Println(info)
		}
	}
	return
}

func (learn *Learn) LearnCreate() (err error) {
	statement := "insert into [$8].[dbo].[guild_learn] values ($1, $2, $3, $4, $5, $6, $7) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(learn.Ask, learn.GuildId, learn.ChannelId, learn.AdminId, learn.Answer, learn.GmtModified, learn.Pass, config.Conf.DatabaseName).Scan(&learn.Id)
	return
}

func (learn *Learn) LearnUpdate(answer null.String) (err error) {
	_, err = Db.Exec("update [$9].[dbo].[guild_learn] set ask = $2, guild_id = $3, channel_id = $4, admin_id = $5, answer = $6, gmt_modified = $7, pass = $8 where ID = $1", learn.Id, learn.Ask, learn.GuildId, learn.ChannelId, learn.AdminId, answer.String, learn.GmtModified, learn.Pass, config.Conf.DatabaseName)

	return
}

// LearnDeleteByGuildIdAndAskAndChannelId
func (learn *Learn) LDBGIAAACI() (err error) {
	_, err = Db.Exec("delete from [$2].[dbo].[guild_learn] where ID = $1", learn.Id, config.Conf.DatabaseName)
	return
}

func LearnSave(ask, guildId, channelId, adminId string, answer null.String, gmtModified time.Time, pass bool) (err error) {
	learn := Learn{
		Ask:         ask,
		GuildId:     guildId,
		ChannelId:   channelId,
		AdminId:     adminId,
		Answer:      answer,
		GmtModified: gmtModified,
		Pass:        pass,
	}
	learn_get, err := LearnGet(guildId, channelId, ask)
	fmt.Println(learn_get, err)
	if err != nil {
		err = learn.LearnCreate()
		fmt.Println("创建", err)
		return
	}
	err = learn_get.LearnUpdate(answer)
	fmt.Println("更新", err)
	return
}

func LDBGAA(guildId, channelId, ask string) (err error) {
	learn_get, err := LearnGet(guildId, channelId, ask)
	if err != nil {
		return
	}
	learn_get.LDBGIAAACI()
	return
}
