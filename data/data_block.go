package database

import (
	"fmt"
	"time"

	"github.com/2mf8/QQBotOffical/public"
	_ "github.com/denisenkom/go-mssqldb"
	_ "gopkg.in/guregu/null.v3/zero"
)

type PBlock struct {
	Id          int       `json:"id"`
	GuildId     string    `json:"guild_id"`
	UserId      string    `json:"user_id"`
	IsPBlock    bool      `json:"ispblock"`
	AdminId     string    `json:"admin_id"`
	GmtModified time.Time `json:"gmt_modified"`
}

func PBlockGet(guildId, userId string) (p PBlock, err error) {
	pblock := PBlock{}
	err = Db.QueryRow("select ID, guild_id, user_id, admin_id, gmt_modified, ispblock from [kequ5060].[dbo].[guild_pblock] where user_id = $1 and ispblock = $2 and guild_id = $3", userId, true, guildId).Scan(&pblock.Id, &pblock.GuildId, &pblock.UserId, &pblock.AdminId, &pblock.GmtModified, &pblock.IsPBlock)
	info := fmt.Sprintf("%s", err)
	if public.StartsWith(info, "sql") || public.StartsWith(info, "unable") {
		if public.StartsWith(info, "unable") {
			fmt.Println(info)
		}
		return
	}
	return
}

func (pBlock *PBlock) PBlockCreate() (err error) {
	statement := "insert into [kequ5060].[dbo].[guild_pblock] values ($1, $2, $3, $4, $5) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(pBlock.GuildId, pBlock.UserId, pBlock.AdminId, pBlock.GmtModified, pBlock.IsPBlock).Scan(&pBlock.Id)
	return
}

func (pBlock *PBlock) PBlockUpdate(ispblock bool) (err error) {
	_, err = Db.Exec("update [kequ5060].[dbo].[guild_pblock] set guild_id = $2, user_id = $3, ispblock = $4, admin_id = $5, gmt_modified = $6 where ID = $1", pBlock.Id, pBlock.GuildId, pBlock.UserId, pBlock.IsPBlock, pBlock.AdminId, pBlock.GmtModified)
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
	if err != nil {
		pblock.PBlockCreate()
		return
	}
	pblock_get.PBlockUpdate(ispblock)
	return
}
