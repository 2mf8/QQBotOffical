package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/public"
	"github.com/2mf8/QQBotOffical/utils"
)

type WCA struct {
}

func (wca *WCA) Do(ctx *context.Context, gmap map[string][]string, guildId, channelId, userId, msg, msgId, username, avatar, srcGuildID string, useRole []string, isBot, isDirectMessage, botIsAdmin bool, priceSearch string, attachments []string) utils.RetStuct {
	s, b := public.Prefix(msg, ".")
	if !b {
		return utils.RetStuct{
			RetVal: utils.MESSAGE_IGNORE,
		}
	}
	if strings.HasPrefix(s, "wca") {
		w_m := strings.TrimSpace(strings.TrimSpace(string([]byte(s)[len("wca"):])))
		fmt.Println(w_m)
		url := "http://www.2mf8.cn:8100/wcaPerson/searchPeople?q=" + url.QueryEscape(w_m)
		resp, _ := http.Get(url)
		s := database.Info{}
		gen := ""
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal([]byte(body), &s)
		//fmt.Println(fmt.Sprintf("%+v", s))
		//fmt.Println(len(s.Data))
		if s.TotalElements == 1 {
			//fmt.Println(s.Data[0].Id)
			if s.Data[0].Gender == "m" {
				gen = "Male"
			} else {
				gen = "Female"
			}
			s_r := s.Data[0].Name + "\n" + s.Data[0].Id + "," + s.Data[0].CountryId + "," + gen + database.WcaPersonHandler(s.Data[0].Id)
			// fmt.Println(s_r)
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: s_r,
				},
				MsgId: msgId,
			}
		} else if s.TotalElements > 99 {
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: "搜索范围太大！",
				},
				MsgId: msgId,
			}
		} else {
			rankList := ""
			s_r := ""
			count := 0
			for _, l := range s.Data {
				if count < 4 {
					rankList += "\n" + l.Id + " | " + l.Name
				} else {
					rankList += "\n" + l.Id + " | " + l.Name
					if count == 19 {
						rankList += "\n..."
						break
					}
				}
				count++
			}
			if s.TotalElements == 0 {
				s_r = "暂无相关记录，请换个名字或输入对应的WCAID进行搜索。"
			} else {
				s_r = strconv.Itoa(s.TotalElements) + "条记录" + rankList + "\n请换个名字或输入对应的WCAID进行搜索。"
			}
			return utils.RetStuct{
				RetVal: utils.MESSAGE_BLOCK,
				ReplyMsg: &utils.Msg{
					Text: s_r,
				},
				MsgId: msgId,
			}
		}
	}
	return utils.RetStuct{
		RetVal: utils.MESSAGE_IGNORE,
	}
}

func init() {
	utils.Register("WCA", &WCA{})
}
