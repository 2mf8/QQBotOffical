package database

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/2mf8/QQBotOffical/public"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gomodule/redigo/redis"
	_ "gopkg.in/guregu/null.v3/zero"
)

type PBlock struct {
	Id          int `json:"id"`
	GuildId     string `json:"guild_id"`
	UserId      string `json:"user_id"`
	IsPBlock    bool `json:"ispblock"`
	AdminId     string `json:"admin_id"`
	GmtModified time.Time `json:"gmt_modified"`
}

type PBlockSync struct {
	IsTrue     bool `json:"synchronization"`
	PBlockSync *PBlock
}

func PBlockGet(guildId, userId string) (pblockSync PBlockSync, err error) {
	pblock := PBlock{}
	pblockSync = PBlockSync{
		IsTrue:     true,
		PBlockSync: &pblock,
	}

	bw := guildId + "_pblock_" + userId
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", bw)
	c.Flush()

	var vb []byte
	vb, err = redis.Bytes(c.Receive())
	if err != nil {
		fmt.Println("[查询] 首次查询-个人屏蔽", bw)
		err = Db.QueryRow("select ID, guild_id, user_id, admin_id, gmt_modified, ispblock from [kequ5060].[dbo].[guild_pblock] where user_id = $1 and ispblock = $2 and guild_id = $3", userId, true, guildId).Scan(&pblock.Id, &pblock.GuildId, &pblock.UserId, &pblock.AdminId, &pblock.GmtModified, &pblock.IsPBlock)
		info := fmt.Sprintf("%s", err)
		if public.StartsWith(info, "sql") || public.StartsWith(info, "unable") {
			if public.StartsWith(info, "unable") {
				fmt.Println(info)
			}
			pblockSync = PBlockSync{
				IsTrue:     false,
				PBlockSync: &pblock,
			}
		}
		var bw_set []byte
		bw_set, _ = json.Marshal(&pblockSync)
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, _ := c.Receive()
		fmt.Printf("[收到] %#v\n", v)
		return
	}
	err = json.Unmarshal(vb, &pblockSync)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	//fmt.Println("[Redis] Key(", bw, ") Value(", pblockSync.IsTrue, *pblockSync.PBlockSync, ")") //测试用
	return
}

func (pBlock *PBlock) PBlockCreate() (err error) {
	statement := "insert into [kequ5060].[dbo].[guild_pblock] (guild_id, user_id, admin_id, gmt_modified, ispblock) values ($1, $2, $3, $4, $5) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(pBlock.GuildId, pBlock.UserId, pBlock.AdminId, pBlock.GmtModified, pBlock.IsPBlock).Scan(&pBlock.Id)

	pblockSync := PBlockSync{
		IsTrue: true,
		PBlockSync: &PBlock{
			Id:          pBlock.Id,
			GuildId:     pBlock.GuildId,
			UserId:      pBlock.UserId,
			IsPBlock:    pBlock.IsPBlock,
			AdminId:     pBlock.AdminId,
			GmtModified: pBlock.GmtModified,
		},
	}

	bw := pBlock.GuildId + "_pblock_" + pBlock.UserId
	var bw_set []byte
	bw_set, _ = json.Marshal(&pblockSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	_ = fmt.Sprintf("%#v", v)
	return
}

func (pBlock *PBlock) PBlockUpdate(ispblock bool) (err error) {
	_, err = Db.Exec("update [kequ5060].[dbo].[guild_pblock] set guild_id = $2, user_id = $3, ispblock = $4, admin_id = $5, gmt_modified = $6 where ID = $1", pBlock.Id, pBlock.GuildId, pBlock.UserId, pBlock.IsPBlock, pBlock.AdminId, pBlock.GmtModified)

	pblockSync := PBlockSync{
		IsTrue: true,
		PBlockSync: &PBlock{
			Id:          pBlock.Id,
			GuildId:     pBlock.GuildId,
			UserId:      pBlock.UserId,
			IsPBlock:    ispblock,
			AdminId:     pBlock.AdminId,
			GmtModified: pBlock.GmtModified,
		},
	}

	bw := pBlock.GuildId + "_pblock_" + pBlock.UserId
	var bw_set []byte
	bw_set, _ = json.Marshal(&pblockSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	_ = fmt.Sprintf("%#v", v)
	return
}

func PBlockSave(guildId, userId, adminId string, ispblock bool, gmtModified time.Time) (err error) {
	pblock := PBlock{
		GuildId:     guildId,
		UserId:      userId,
		IsPBlock:    ispblock,
		AdminId:     adminId,
		GmtModified: gmtModified,
	}
	pblock_get, err := PBlockGet(guildId, userId)
	if err != nil || !pblock_get.IsTrue {
		pblock.PBlockCreate()
		return
	}
	pblock_get.PBlockSync.PBlockUpdate(ispblock)
	return
}
