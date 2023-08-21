package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	database "github.com/2mf8/QQBotOffical/data"
	"github.com/2mf8/QQBotOffical/status"
	"github.com/gin-gonic/gin"
)

func IndexApi(c *gin.Context) {
	c.String(int(status.OK), "It works")
}

func PriceAddAndUpdateByItemApi(c *gin.Context) {
	citem := c.Param("item")
	sn := c.Param("service_number")
	sng, _ := c.Get("server_number")
	if !(sn == sng || sng == "10000") {
		c.JSON(int(status.ExpectationFailed), gin.H{
			"code": status.Forbidden,
			"msg":  "禁止访问，您权限不足。",
		})
		c.Abort()
		return
	}
	cp, err := database.GetItem(sn, sn, citem)
	if err != nil {
		len := c.Request.ContentLength
		body := make([]byte, len)
		c.Request.Body.Read(body)
		json.Unmarshal(body, &cp)
		cid := cp.Id
		if cp.Item == "" {
			c.JSON(int(status.NoContent), gin.H{
				"code": status.GetError,
				"msg":  "item不能为空",
			})
			c.Abort()
			return
		} else {
			cu, err := database.GetItem(sn, sn, cp.Item)
			gid := cu.Id
			if err != nil {
				err = cp.ItemCreate()
				if err != nil {
					msg := fmt.Sprintf("创建%s失败", cp.Item)
					c.JSON(int(status.InternalServerError), gin.H{
						"code": status.CreateError,
						"msg":  msg,
					})
					c.Abort()
					return
				} else {
					msg := fmt.Sprintf("创建%s成功", cp.Item)
					c.JSON(int(status.OK), gin.H{
						"code": status.OK,
						"msg":  msg,
					})
					c.Abort()
					return
				}
			} else {
				if cid == gid {
					err = cp.ItemUpdate()
				} else {
					err1 := cu.ItemDeleteById()
					fmt.Println(err1)
					err = cp.ItemUpdate()
				}
				if err != nil {
					msg := fmt.Sprintf("更新%s失败", cp.Item)
					c.JSON(int(status.InternalServerError), gin.H{
						"code": status.UpdateError,
						"msg":  msg,
					})
					c.Abort()
					return
				} else {
					msg := fmt.Sprintf("更新%s成功", cp.Item)
					c.JSON(int(status.OK), gin.H{
						"code": status.OK,
						"msg":  msg,
					})
					c.Abort()
					return
				}
			}
		}
	} else {
		len := c.Request.ContentLength
		body := make([]byte, len)
		c.Request.Body.Read(body)
		json.Unmarshal(body, &cp)
		err = cp.ItemUpdate()
		if err != nil {
			msg := fmt.Sprintf("更新%s失败", cp.Item)
			c.JSON(int(status.InternalServerError), gin.H{
				"code": status.UpdateError,
				"msg":  msg,
			})
			c.Abort()
			return
		} else {
			msg := fmt.Sprintf("更新%s成功", cp.Item)
			c.JSON(int(status.OK), gin.H{
				"code": status.OK,
				"msg":  msg,
			})
			c.Abort()
			return
		}
	}
}

func PriceDeleteByItemApi(c *gin.Context) {
	citem := c.Param("item")
	sn := c.Param("service_number")
	sng, _ := c.Get("server_number")
	if !(sn == sng || sng == "10000") {
		c.JSON(int(status.ExpectationFailed), gin.H{
			"code": status.Forbidden,
			"msg":  "禁止访问，您权限不足。",
		})
		c.Abort()
		return
	}
	cp, err := database.GetItem(sn, sn, citem)
	if err != nil {
		fmt.Println(err)
		msg := fmt.Sprintf("获取%s失败", cp.Item)
		c.JSON(int(status.InternalServerError), gin.H{
			"code": status.GetError,
			"msg":  msg,
		})
		c.Abort()
		return
	} else {
		err = cp.ItemDeleteById()
		if err != nil {
			msg := fmt.Sprintf("删除%s失败", cp.Item)
			c.JSON(int(status.InternalServerError), gin.H{
				"code": status.DeleteError,
				"msg":  msg,
			})
			c.Abort()
			return
		} else {
			msg := fmt.Sprintf("删除%s成功", cp.Item)
			c.JSON(int(status.OK), gin.H{
				"code": status.OK,
				"msg":  msg,
			})
			c.Abort()
			return
		}
	}
}

func PriceGetItemApi(c *gin.Context) {
	citem, _ := url.QueryUnescape(c.Param("item"))
	sn := c.Param("service_number")
	cp, err := database.GetItem(sn, sn, citem)
	if err != nil {
		fmt.Println(err)
		msg := fmt.Sprintf("获取%s失败", citem)
		c.JSON(int(status.InternalServerError), gin.H{
			"code": status.GetError,
			"msg":  msg,
		})
		c.Abort()
		return
	} else {
		/*op, err := json.MarshalIndent(&cp, "", "\t")
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"msg": "转换为JSON失败",
			})
		}
		fmt.Println(string(op))*/
		c.JSON(int(status.OK), gin.H{
			"code": status.OK,
			"msg":  fmt.Sprintf("获取%s成功", citem),
			"data": cp,
		})
		c.Abort()
		return
	}
}

func PriceGetItemsApi(c *gin.Context) {
	citem := c.Param("key")
	sn := c.Param("service_number")
	shop := ""
	QQGuild := ""
	if sn == "10001" {
		shop = "黄小姐的魔方店"
		QQGuild = "https://pd.qq.com/s/9ngvzfmbg"
	}
	if sn == "10002" {
		shop = "奇乐魔方坊"
		QQGuild = "https://pd.qq.com/s/af5gmzqhh"
	}
	cp, err := database.GetItems(sn, sn, citem)
	//fmt.Println(cp)
	if err != nil {
		fmt.Println(err)
		msg := fmt.Sprintf("获取%s失败", citem)
		c.JSON(int(status.InternalServerError), gin.H{
			"code":    status.GetError,
			"msg":     msg,
			"shop":    shop,
			"QQGuild": QQGuild,
		})
		c.Abort()
		return
	} else if len(cp) == 0 {
		c.JSON(int(status.ExpectationFailed), gin.H{
			"code":    status.ExpectationFailed,
			"msg":     fmt.Sprintf("获取%s成功, 但未查询到数据。", citem),
			"shop":    shop,
			"QQGuild": QQGuild,
		})
		c.Abort()
		return
	} else {
		/*op, err := json.MarshalIndent(&cp, "", "\t")
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"msg": "转换为JSON失败",
			})
		}
		fmt.Println(string(op))*/
		c.JSON(int(status.OK), gin.H{
			"code":    status.OK,
			"msg":     fmt.Sprintf("获取%s成功", citem),
			"shop":    shop,
			"QQGuild": QQGuild,
			"data":    cp,
		})
		c.Abort()
		return
	}
}

func PriceGetItemsAllApi(c *gin.Context) {
	citem := c.Param("key")
	sn := c.Param("service_number")
	shop := ""
	QQGuild := ""
	if sn == "10001" {
		shop = "黄小姐的魔方店"
		QQGuild = "https://pd.qq.com/s/9ngvzfmbg"
	}
	if sn == "10002" {
		shop = "奇乐魔方坊"
		QQGuild = "https://pd.qq.com/s/af5gmzqhh"
	}
	cp, err := database.GetItems(sn, sn, citem)
	//fmt.Println(cp)
	if err != nil {
		fmt.Println(err)
		msg := fmt.Sprintf("获取%s失败", citem)
		c.JSON(int(status.InternalServerError), gin.H{
			"code":    status.GetError,
			"msg":     msg,
			"shop":    shop,
			"QQGuild": QQGuild,
		})
		c.Abort()
		return
	} else if len(cp) == 0 {
		c.JSON(int(status.ExpectationFailed), gin.H{
			"code":    status.ExpectationFailed,
			"msg":     fmt.Sprintf("获取%s成功, 但未查询到数据。", citem),
			"shop":    shop,
			"QQGuild": QQGuild,
		})
		c.Abort()
		return
	} else {
		/*op, err := json.MarshalIndent(&cp, "", "\t")
		if err != nil {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"msg": "转换为JSON失败",
			})
		}
		fmt.Println(string(op))*/
		c.JSON(int(status.OK), gin.H{
			"code":    status.OK,
			"msg":     fmt.Sprintf("获取%s成功", citem),
			"shop":    shop,
			"QQGuild": QQGuild,
			"data":    cp,
		})
		c.Abort()
		return
	}
}
