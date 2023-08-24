package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gomodule/redigo/redis"
)

type ServerNumberSet struct {
	ServerNumbers map[string]string
	Intent        map[string]int
}

type ServerNumberSetSync struct {
	IsTrue              bool
	ServerNumberSetSync *ServerNumberSet
}

func (s *ServerNumberSet) ServerNumbersSet() error {
	output, err := json.MarshalIndent(&s, "", "\t")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return err
	}
	err = os.WriteFile("server_numbers.json", output, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file", err)
		return err
	}
	return nil
}

func ServerNumbersRead() (s ServerNumberSet, err error) {
	jsonFile, err := os.Open("server_numbers.json")
	if err != nil {
		fmt.Println("Error reading JSON File:", err)
		return
	}
	defer jsonFile.Close()
	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}
	json.Unmarshal(jsonData, &s)
	//fmt.Println(k)
	return
}

func ServerNumbersGet() (s ServerNumberSetSync, err error) {
	sns := ServerNumberSet{
		ServerNumbers: map[string]string{},
		Intent:        map[string]int{},
	}
	snss := ServerNumberSetSync{
		IsTrue:              true,
		ServerNumberSetSync: &sns,
	}
	var vb []byte
	var bw_set []byte

	bw := "server_numbers"
	c := Pool.Get()
	defer c.Close()
	c.Send("Get", bw)
	c.Flush()
	vb, err = redis.Bytes(c.Receive())
	if err != nil {
		fmt.Println("[查询] 首次查询-守卫", bw)
		snr, err := ServerNumbersRead()
		if err != nil {
			snss = ServerNumberSetSync{
				IsTrue:              false,
				ServerNumberSetSync: &sns,
			}
			snss.ServerNumberSetSync.ServerNumbersSet()
		}
		snss.ServerNumberSetSync = &snr
		bw_set, _ = json.Marshal(&snss)
		c.Send("Set", bw, bw_set)
		c.Flush()
		v, _ := c.Receive()
		fmt.Printf("[收到] %#v\n", v)
		return snss, err
	}
	err = json.Unmarshal(vb, &snss)
	if err != nil {
		fmt.Println("[错误] Unmarshal出错")
	}
	fmt.Println("[Redis] Key(", bw, ") Value(", snss.IsTrue, *snss.ServerNumberSetSync, ")") //测试用
	return snss, err
}

func (s *ServerNumberSetSync) ServerNumberUpdate(sn, value string, intent int) error {
	if s.ServerNumberSetSync.ServerNumbers == nil {
		var _tsnm map[string]string = map[string]string{}
		var _tinm map[string]int = map[string]int{}
		_tsnm[sn] = value
		_tinm[value] = intent
		tsnm := ServerNumberSet{
			ServerNumbers: _tsnm,
		}
		s.ServerNumberSetSync = &tsnm
	}
	if s.ServerNumberSetSync.Intent == nil {
		var _tinm map[string]int = map[string]int{}
		_tinm[value] = intent
		tsnm := ServerNumberSet{
			Intent: _tinm,
		}
		s.ServerNumberSetSync = &tsnm
	}
	s.ServerNumberSetSync.ServerNumbers[sn] = value
	s.ServerNumberSetSync.Intent[value] = intent
	bw := "server_numbers"
	var bw_set []byte
	serverNumberSetSync := ServerNumberSetSync{
		IsTrue:              true,
		ServerNumberSetSync: s.ServerNumberSetSync,
	}
	bw_set, _ = json.Marshal(&serverNumberSetSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	_ = fmt.Sprintf("%#v", v)
	err = s.ServerNumberSetSync.ServerNumbersSet()
	return err
}

func (s *ServerNumberSetSync) ServerNumbersDelete(sn, value string) {
	delete(s.ServerNumberSetSync.Intent, value)
	delete(s.ServerNumberSetSync.ServerNumbers, sn)
	bw := "server_numbers"
	var bw_set []byte
	serverNumberSetSync := ServerNumberSetSync{
		IsTrue:              true,
		ServerNumberSetSync: s.ServerNumberSetSync,
	}
	bw_set, _ = json.Marshal(&serverNumberSetSync)
	c := Pool.Get()
	defer c.Close()
	c.Send("Set", bw, bw_set)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		fmt.Println("[错误] Receive出错")
	}
	_ = fmt.Sprintf("%#v", v)
	s.ServerNumberSetSync.ServerNumbersSet()
}
