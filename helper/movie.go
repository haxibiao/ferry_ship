/*
 * @Author: Bin
 * @Date: 2021-07-08
 * @FilePath: /ferry_ship/helper/movie.go
 */
package helper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

		// fmt.Printf("message=%+v\n", msg.Elements)

		for _, elem := range msg.Elements {

			// 判断是 @ 用户消息类型
			if elem.Type() == message.At {

				// 判断是否 @ 当前机器人并触发搜索关键词
				mKeys := []string{"@" + b.Nickname + " 搜索", "@" + b.Nickname + " 搜索"}

				for _, value := range mKeys {
					if strings.Contains(msg.ToString(), value) {

						// 开始搜索电影
						out := autoreply(strings.Replace(msg.ToString(), value, "", -1))
						if out == "" {
							return
						}

						// 生成回复消息并发送
						m := message.NewSendingMessage().Append(message.NewText(out))
						c.SendGroupMessage(msg.GroupCode, m)

						// 匹配成功一次之后就跳出匹配
						break

					}
				}

			}
		}

		fmt.Println("【收到消息】" + msg.ToString())

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
		return "瞎猫视频（xiamaoshipin.com）AI 好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://xiamaoshipin.com/app"
	}

	out, ok := SearchMovie(in)
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
			text = "瞎猫视频（xiamaoshipin.com）AI 帮您搜索到 " + strconv.Itoa(quantity) + " 条《" + keywords + "》相关内容：\n" + seekMovies
		} else {
			text = "瞎猫视频（xiamaoshipin.com）AI 很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://xiamaoshipin.com/app"
		}
	}

	// fmt.Println("【消息】" + text)
	return text, true
}
