package database

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func RedisGet(key string, PBStruct any) (any, error) {
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", key)
	c.Flush()
	vb, err := redis.Bytes(c.Receive())
	if err != nil {
		return PBStruct, err
	}
	err = json.Unmarshal(vb, &PBStruct)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	fmt.Println("[Redis] Key(", key, ") Value(", PBStruct, ")") //测试用
	return PBStruct, err
}

func RedisSet(key string, PBStruct any) {
	bw_set, _ := json.Marshal(&PBStruct)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", key, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	_ = fmt.Sprintf("%#v", v)
}
