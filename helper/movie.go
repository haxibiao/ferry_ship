/*
 * @Author: Bin
 * @Date: 2021-07-08
 * @FilePath: /ferry_ship/helper/movie.go
 */
package helper

import (
	"crypto/tls"
	"encoding/json"
	"regexp"
	"strconv"
	"sync"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/beego/beego/v2/client/httplib"

	"ferry_ship/bot"
	"ferry_ship/bot/utils"
)

func init() {
	// bot.RegisterModule(instance)
}

var MovieInstance = &movie{}

var logger = utils.GetModuleLogger("bin.moviereply")

// var tem map[string]string

// var BaseWebURL = "xiamaoshipin.com"
// var BaseBotName = "瞎猫视频"

var BaseWebURL = "xiaocaihong.tv"
var BaseBotName = "小彩虹视频"

type movie struct {
}

func (mov *movie) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "bin.moviereply",
		Instance: MovieInstance,
	}
}

func (mov *movie) Init() {
	// path := config.GlobalConfig.GetString("logiase.autoreply.path")

	// if path == "" {
	// 	path = "./autoreply.yaml"
	// }

	// bytes := utils.ReadFile(path)
	// err := yaml.Unmarshal(bytes, &tem)
	// if err != nil {
	// 	logger.WithError(err).Errorf("unable to read config file in %s", path)
	// }
}

func (mov *movie) PostInit() {
}

func (mov *movie) Serve(b *bot.Bot) {
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {

		// fmt.Printf("message=%+v\n", msg.Sender.Nickname)

		if botObj := bot.Instances[c.Uin]; botObj == nil {
			// 机器人已下线，直接结束回复流程
			// fmt.Println("【收到消息】机器人已下线，直接结束回复流程")
			return
		}

		groupInfo := c.FindGroup(msg.GroupCode)
		if groupInfo == nil {
			// QQ 群信息获取失败，结束流程
			// fmt.Println("【收到消息】QQ 群信息获取失败，直接结束回复流程")
			return
		}
		groupMemberInfo := groupInfo.FindMember(c.Uin)
		botName := ""
		if groupMemberInfo == nil {
			// QQ 群我的数据获取失败，直接赋值昵称
			botName = b.Nickname
		} else {
			botName = groupMemberInfo.DisplayName()
		}

		// fmt.Printf("群昵称=%+v\n", botName)

		for _, elem := range msg.Elements {

			// 判断是 @ 用户消息类型
			if elem.Type() == message.At {

				// 判断是否 @ 当前机器人并触发搜索关键词
				mKeys := []string{"@" + botName + " 搜索 ", "@" + botName + " 搜索"}

				for _, value := range mKeys {

					// 正则匹配电影名称
					flysnowRegexp := regexp.MustCompile(value + `(.+)$`)
					params := flysnowRegexp.FindStringSubmatch(msg.ToString())
					if len(params) <= 0 {
						// 判断没有匹配到关键词，进入下一次循环
						continue
					}
					movieKey := params[1] // 提取匹配到的关键词
					// fmt.Printf("群电影名=%+v\n", params[1])

					if movieKey != "" {

						// 开始搜索电影
						out := autoreply(movieKey)
						if out == "" {
							return
						}

						// 生成回复消息并发送
						m := message.NewSendingMessage().Append(message.NewText(out))
						c.SendGroupMessage(msg.GroupCode, m)

						// 匹配成功一次之后就跳出匹配
						return
					}
				}

				// 未触发关键词，生成默认回复消息并发送
				// out := autoreply("")
				// m := message.NewSendingMessage().Append(message.NewText(out))
				// c.SendGroupMessage(msg.GroupCode, m)

			}
		}

		// fmt.Println("【收到消息】" + msg.ToString())

	})

	// b.OnPrivateMessage(func(c *client.QQClient, msg *message.PrivateMessage) {
	// 	out := autoreply(msg.ToString())
	// 	if out == "" {
	// 		return
	// 	}
	// 	m := message.NewSendingMessage().Append(message.NewText(out))
	// 	c.SendPrivateMessage(msg.Sender.Uin, m)
	// })
}

func (mov *movie) Start(bot *bot.Bot) {
}

func (mov *movie) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func autoreply(in string) string {

	if in == "" {
		return BaseBotName + "（" + BaseWebURL + "）好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://" + BaseWebURL + "/app"
	}

	out, ok := SearchMovie2(in)
	if !ok {
		return ""
	}
	return out
}

func SearchMovie(keywords string) (callback string, ok bool) {
	seekApi := "https://xiamaoshipin.com/api/movie/search?keyword=" + keywords
	req := httplib.Get(seekApi)
	str, err := req.String()
	if err != nil {
		return "", false
		// t.Fatal(err)
	}

	text := ""
	seekMovies := ""

	var data_obj interface{}
	json.Unmarshal([]byte(str), &data_obj)

	data := data_obj.(map[string]interface{})

	if data["data"] != nil {

		movies := data["data"].([]interface{})
		quantity := len(movies)

		i := 1
		for _, value := range movies {
			seekMovies = seekMovies + "\n" + strconv.Itoa(i) + "，《" + value.(map[string]interface{})["name"].(string) + "》，立即观看：" + value.(map[string]interface{})["watch_url"].(string)
			i++
		}

		if quantity > 0 {
			text = BaseBotName + "（" + BaseWebURL + "）帮您搜索到 " + strconv.Itoa(quantity) + " 条《" + keywords + "》相关内容：\n" + seekMovies
		} else {
			text = BaseBotName + "（" + BaseWebURL + "）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://" + BaseWebURL + "/app"
		}
	}

	// fmt.Println("【消息】" + text)
	return text, true
}

func SearchMovie2(keywords string) (callback string, ok bool) {
	seekApi := "https://xiaocaihong.tv/api/movie/qq/search/"
	req := httplib.Post(seekApi)
	req.Param("q", keywords)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	str, err := req.String()
	if err != nil {
		return "", false
		// t.Fatal(err)
	}

	text := ""
	seekMovies := ""

	var data_obj interface{}
	json.Unmarshal([]byte(str), &data_obj)

	if data_obj != nil {

		movies := data_obj.([]interface{})
		quantity := len(movies)

		i := 1
		for _, value := range movies {

			var movieName = value.(map[string]interface{})["name"].(string)
			var movieUrl = value.(map[string]interface{})["url"].(string)

			seekMovies = seekMovies + "\n" + strconv.Itoa(i) + "，《" + movieName + "》，立即观看：" + movieUrl
			i++
		}

		if quantity > 0 {
			text = BaseBotName + "（" + BaseWebURL + "）帮您搜索到 " + strconv.Itoa(quantity) + " 条《" + keywords + "》相关内容：\n" + seekMovies
		} else {
			text = BaseBotName + "（" + BaseWebURL + "）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://" + BaseWebURL + "/app"
		}
	}

	// fmt.Println("【消息】" + text)
	return text, true
}
