package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/2mf8/QQBotOffical/public"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gomodule/redigo/redis"
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

type LearnSync struct {
	IsTrue    bool `json:"synchronization"`
	LearnSync *Learn
}

func LearnGet(guildId, channelId, ask string) (learnSync LearnSync, err error) {
	learn := Learn{}
	learnSync = LearnSync{
		IsTrue:    true,
		LearnSync: &learn,
	}

	var vb []byte
	var bw_set []byte

	bw := guildId + "_" + channelId + "_" + ask
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", bw)
	c.Flush()
	vb, err = redis.Bytes(c.Receive())
	if err != nil {
		fmt.Println("[查询] 首次查询-学习", bw)
		err = Db.QueryRow("select ID, ask, guild_id, channel_id, admin_id, answer, gmt_modified, pass from [kequ5060].[dbo].[guild_learn] where guild_id = $1 and ask = $2 and channel_id = $3", guildId, ask, channelId).Scan(&learn.Id, &learn.Ask, &learn.GuildId, &learn.ChannelId, &learn.AdminId, &learn.Answer, &learn.GmtModified, &learn.Pass)
		info := fmt.Sprintf("%s", err)
		if public.StartsWith(info, "sql") || public.StartsWith(info, "unable") {
			if public.StartsWith(info, "unable") {
				fmt.Println(info)
			}
			learnSync = LearnSync{
				IsTrue:    false,
				LearnSync: &learn,
			}
		}
		bw_set, _ = json.Marshal(&learnSync)
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, _ := c.Receive()
		fmt.Printf("[收到] %#v\n", v)
		return
	}
	err = json.Unmarshal(vb, &learnSync)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	// fmt.Println("[Redis] Key(", bw, ") Value(", learnSync.IsTrue, *learnSync.LearnSync, ")") //测试用
	return
}

func (learn *Learn) LearnCreate() (err error) {
	statement := "insert into [kequ5060].[dbo].[guild_learn] (ask, guild_id, channel_id, admin_id, answer, gmt_modified, pass) values ($1, $2, $3, $4, $5, $6, $7) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(learn.Ask, learn.GuildId, learn.ChannelId, learn.AdminId, learn.Answer, learn.GmtModified, learn.Pass).Scan(&learn.Id)

	learnSync := LearnSync{
		IsTrue: true,
		LearnSync: &Learn{
			Id:          learn.Id,
			Ask:         learn.Ask,
			GuildId:     learn.GuildId,
			ChannelId:   learn.ChannelId,
			AdminId:     learn.AdminId,
			Answer:      learn.Answer,
			GmtModified: learn.GmtModified,
			Pass:        learn.Pass,
		},
	}

	bw := learn.GuildId + "_" + learn.ChannelId + "_" + learn.Ask
	var bw_set []byte
	bw_set, _ = json.Marshal(&learnSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Sprintf("%#v", v)
	return
}

func (learn *Learn) LearnUpdate(answer null.String) (err error) {
	_, err = Db.Exec("update [kequ5060].[dbo].[guild_learn] set ask = $2, guild_id = $3, channel_id = $4, admin_id = $5, answer = $6, gmt_modified = $7, pass = $8 where ID = $1", learn.Id, learn.Ask, learn.GuildId, learn.ChannelId, learn.AdminId, answer.String, learn.GmtModified, learn.Pass)

	learnSync := LearnSync{
		IsTrue: true,
		LearnSync: &Learn{
			Id:          learn.Id,
			Ask:         learn.Ask,
			GuildId:     learn.GuildId,
			ChannelId:   learn.ChannelId,
			AdminId:     learn.AdminId,
			Answer:      answer,
			GmtModified: learn.GmtModified,
			Pass:        learn.Pass,
		},
	}

	bw := learn.GuildId + "_" + learn.ChannelId + "_" + learn.Ask
	var bw_set []byte
	bw_set, _ = json.Marshal(&learnSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Sprintf("%#v", v)

	return
}

// LearnDeleteByGuildIdAndAskAndChannelId
func (learn *Learn) LDBGIAAACI() (err error) {
	_, err = Db.Exec("delete from [kequ5060].[dbo].[guild_learn] where ID = $1", learn.Id)
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
	if err != nil || learn_get.IsTrue == false {
		err = learn.LearnCreate()
		return
	}
	err = learn_get.LearnSync.LearnUpdate(answer)
	return
}

func LDBGAA(guildId, channelId, ask string) (err error) {
	learn_get, err := LearnGet(guildId, channelId, ask)
	if err != nil {
		return
	}
	learn_get.LearnSync.LDBGIAAACI()

	learnSync := LearnSync{
		IsTrue:    true,
		LearnSync: &Learn{},
	}

	bw := learn_get.LearnSync.GuildId + "_" + learn_get.LearnSync.ChannelId + "_" + ask
	var bw_set []byte
	bw_set, _ = json.Marshal(&learnSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	fmt.Sprintf("%#v", v)
	return
}
