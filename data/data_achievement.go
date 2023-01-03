package database

import "fmt"

type Achievement struct {
	Id       int64
	UserId   string
	UserName string
	Item     string
	Best     int
	Average  int
	Session  int
}

func (a *Achievement) AchievementCreate() (err error) {
	statement := "insert into [kequ5060].[dbo].[guild_achievement] values ($1, $2, $3, $4, $5, $6) select @@identity"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(a.UserId, a.UserName, a.Item, a.Best, a.Average, a.Session).Scan(&a.Id)
	return
}

func (a *Achievement) AchievementUpdate(username string, best, average int) (err error) {
	_, err = Db.Exec("update [kequ5060].[dbo].[guild_achievement] set user_id = $2, user_name = $3, item = $4, best = $5, average = $6, session = $7 where Id = $1", a.Id, a.UserId, username, a.Item, best, average, a.Session)
	return
}

func AchievementSave(userId, userName, item string, best, average, session int) (err error) {
	a := Achievement{
		UserId:   userId,
		UserName: userName,
		Item:     item,
		Best:     best,
		Average:  average,
		Session:  session,
	}
	a_get, err := AchievementGet(userId, item, session)
	if err != nil {
		err = a.AchievementCreate()
		return
	}
	if a_get.Best == -1 && a_get.Average == -1 {
		err = a_get.AchievementUpdate(userName, best, average)
		return
	}
	if a_get.Best == -1 {
		err = a_get.AchievementUpdate(userName, best, a_get.Average)
		return
	}
	if a_get.Average == -1 {
		err = a_get.AchievementUpdate(userName, a_get.Best, average)
		return
	}
	if best < a_get.Best && average < a_get.Average {
		err = a_get.AchievementUpdate(userName, best, average)
		return
	}
	if best < a_get.Best {
		err = a_get.AchievementUpdate(userName, best, a_get.Average)
		return
	}
	if average < a_get.Average {
		err = a_get.AchievementUpdate(userName, a_get.Best, average)
		return
	}
	err = a_get.AchievementUpdate(userName, a_get.Best, a_get.Average)
	return
}

func AchievementGet(userId, item string, session int) (a Achievement, err error) {
	a = Achievement{}
	err = Db.QueryRow("select Id, user_id, user_name, item, best, average, session from [kequ5060].[dbo].[guild_achievement] where user_id = $1 and session = $3 and item = $2", userId, item, session).Scan(&a.Id, &a.UserId, &a.UserName, &a.Item, &a.Best, &a.Average, &a.Session)
	return
}

// AchievementDeleteByUserIdAndSession
func ADBUAS(userId string, session int) (err error) {
	_, err = Db.Exec("delete from [kequ5060].[dbo].[guild_achievement] where user_id = $1 and session = $2", userId, session)
	return
}

// AchievementDeleteByUserIdAndItemAndSession
func ADBUAIAS(userId, item string, session int) (err error) {
	_, err = Db.Exec("delete from [kequ5060].[dbo].[guild_achievement] where user_id = $1 and item = $2 and session = $3", userId, item, session)
	return
}

// AchievementGetByUserIdAndSession
func AGBUAS(userId string, session int) (as []Achievement, err error) {
	statment := fmt.Sprintf("select Id, user_id, user_name, item, best, average, session from [kequ5060].[dbo].[guild_achievement] where user_id = %s and session = %d", userId, session)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := Achievement{}
		err = rows.Scan(&a.Id, &a.UserId, &a.UserName, &a.Item, &a.Best, &a.Average, &a.Session)
		as = append(as, a)
	}
	return
}

// AchievementGetByItemAndSessionOrderByBestAsc
// desc 大 → 小
// asc 小 → 大
func AGBIASOBBA(item string, session int) (as []Achievement, err error) {
	statment := fmt.Sprintf("select Id, user_id, user_name, item, best, average, session from [kequ5060].[dbo].[guild_achievement] where item = %s and session = %d and best > -1 order by best asc", item, session)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := Achievement{}
		err = rows.Scan(&a.Id, &a.UserId, &a.UserName, &a.Item, &a.Best, &a.Average, &a.Session)
		as = append(as, a)
	}
	return
}

// AchievementGetByItemAndSessionOrderByAverageAsc
// desc 大 → 小
// asc 小 → 大
func AGBIASOBAA(item string, session int) (as []Achievement, err error) {
	statment := fmt.Sprintf("select Id, user_id, user_name, item, best, average, session from [kequ5060].[dbo].[guild_achievement] where item = %s and session = %d and average > -1 order by average asc", item, session)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := Achievement{}
		err = rows.Scan(&a.Id, &a.UserId, &a.UserName, &a.Item, &a.Best, &a.Average, &a.Session)
		as = append(as, a)
	}
	return
}

func AchievementGetCount(item string, best, average, session int) (i, j int, err error) {
	i = 0
	j = 0
	statment := fmt.Sprintf("select Id, user_id, user_name, item, best, average, session from [kequ5060].[dbo].[guild_achievement] where item = %s and session = %d", item, session)
	rows, err := Db.Query(statment)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := Achievement{}
		err = rows.Scan(&a.Id, &a.UserId, &a.UserName, &a.Item, &a.Best, &a.Average, &a.Session)
		if a.Best > -1 && a.Best < best {
			i++
		}
		if a.Average > -1 && a.Average < average {
			j++
		}
	}
	return
}
