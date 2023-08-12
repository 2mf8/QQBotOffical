## 简介

爱魔方吧是一个基于 [QQ官方API](https://q.qq.com/#/) 和 [BotGo](https://github.com/tencent-connect/botgo) 实现的频道机器人。

[![QQ魔方频道](https://img.shields.io/static/v1?label=QQ%E9%AD%94%E6%96%B9%E9%A2%91%E9%81%93&message=66tz506re5&color=blue)](https://pd.qq.com/s/4syyazec6)

## 消息支持

仅支持在频道中使用，若后续开通了群聊，将会在此文档中有所体现。

## 机器人获取

魔方频道主若想自己的频道接入爱魔方吧，请加入 [![QQ魔方频道](https://img.shields.io/static/v1?label=QQ%E9%AD%94%E6%96%B9%E9%A2%91%E9%81%93&message=66tz506re5&color=blue)](https://pd.qq.com/s/4syyazec6) ，然后私信频道主你的QQ，等待频道主与你联系。

由于是频道私域机器人，能添加的频道数量有限，故申请的门槛有点高。要求频道大于200人且是魔方频道。想要申请请抓紧，先到先得。(承若接入永久免费。若有人向你许诺达到某种条件即可添加，请不要相信)

## 使用机器人

使用本机器人即代表您同意 爱魔方吧 的 [使用条款](https://2mf8.cn/docs/bot/%E4%BD%BF%E7%94%A8%E6%9D%A1%E6%AC%BE.html#%E6%9D%A1%E6%AC%BE%E8%AF%B4%E6%98%8E) 和 [隐私策略](https://2mf8.cn/docs/bot/%E9%9A%90%E7%A7%81%E7%AD%96%E7%95%A5.html) 。

## 特点

### 稳定、可靠

全天24小时在线，稳定可靠。SDK和API由QQ官方维护，不用担心跑路的风险。

### 功能强大

有频道管理和守卫功能，可以拦截广告并撤回，同时予以广告发布者以禁言加警告。

### 可控行强

每个功能都有开关，可以精准的控制某功能的开启和关闭。

## 指令目录

1. [功能开关](https://2mf8.cn/docs/bot/%E5%8A%9F%E8%83%BD%E5%BC%80%E5%85%B3.html#%E7%AE%80%E4%BB%8B)
2. [频道管理](https://2mf8.cn/docs/bot/%E9%A2%91%E9%81%93%E7%AE%A1%E7%90%86.html#%E7%AE%80%E4%BB%8B)
3. [打乱获取](https://2mf8.cn/docs/bot/%E6%89%93%E4%B9%B1%E8%8E%B7%E5%8F%96.html#%E7%AE%80%E4%BB%8B)
4. [魔方赛季](https://2mf8.cn/docs/bot/%E9%AD%94%E6%96%B9%E8%B5%9B%E5%AD%A3.html#%E7%AE%80%E4%BB%8B)
5. [频道守卫](https://2mf8.cn/docs/bot/%E9%A2%91%E9%81%93%E5%AE%88%E5%8D%AB.html#%E7%AE%80%E4%BB%8B)
6. [频道屏蔽](https://2mf8.cn/docs/bot/%E9%A2%91%E9%81%93%E5%B1%8F%E8%94%BD.html#%E7%AE%80%E4%BB%8B)
7. [随机复读](https://2mf8.cn/docs/bot/%E9%9A%8F%E6%9C%BA%E5%A4%8D%E8%AF%BB.html#%E7%AE%80%E4%BB%8B)
8. [频道学习](https://2mf8.cn/docs/bot/%E9%A2%91%E9%81%93%E5%AD%A6%E4%B9%A0.html#%E7%AE%80%E4%BB%8B)

## 自己部署

1. 安装 MicroSoft Sql Server (Express)  数据库, 使用 sql 文件夹里的 sql 文件生成数据库。
2. 安装 Redis 服务
3. 开始使用。（首次打开会生成配置文件，修改配置文件后即可使用）