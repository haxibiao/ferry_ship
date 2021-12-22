/*
 * @Author: Bin
 * @Date: 2021-07-08
 * @FilePath: /ferry_ship/bot/plugin/movie.go
 */
package plugin

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/Mrs4s/MiraiGo/binary"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/beego/beego/v2/client/httplib"

	beego "github.com/beego/beego/v2/server/web"

	"ferry_ship/bot"
	"ferry_ship/bot/utils"
	"ferry_ship/helper"
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

		if botObj := bot.Instances[c.Uin]; botObj == nil {
			// 机器人已下线，直接结束回复流程
			fmt.Println("【收到消息】机器人已下线，直接结束回复流程")
			return
		}

		groupInfo := c.FindGroup(msg.GroupCode)
		if groupInfo == nil {
			// QQ 群信息获取失败，结束流程
			fmt.Println("【收到消息】QQ 群信息获取失败，直接结束回复流程")
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
			// if elem.Type() == message.At {
			if elem.Type() != message.Voice {

				// 判断是否 @ 当前机器人并触发搜索关键词
				// mKeys := []string{"@" + botName + " 搜索 ", "@" + botName + " 搜索"}
				mKeys := []string{"搜索 ", "搜索", "@" + botName + " 搜索 ", "@" + botName + " 搜索"}

				for _, value := range mKeys {

					// 正则匹配电影名称
					flysnowRegexp := regexp.MustCompile(value + `(.+)$`)
					params := flysnowRegexp.FindStringSubmatch(msg.ToString())
					if len(params) <= 0 {
						// 判断没有匹配到关键词，进入下一次循环
						continue
					}
					movieKey := params[1] // 提取匹配到的关键词
					fmt.Printf("群电影名=%+v\n", params[1])
					logger.Infof("群电影名: %v", params[1])

					if movieKey != "" {

						// 开始搜索电影
						// out := autoreply(movieKey)

						// TODO: 为了满足业务需求，增加一个聚合搜索接口需要传 QQ 账号给后端
						out := autoreply(movieKey, strconv.FormatInt(c.Uin, 10))
						if out == "" {
							return
						}

						// 生成回复消息并发送
						m := message.NewSendingMessage().Append(message.NewText(out))
						msgt := c.SendGroupMessage(msg.GroupCode, m)
						logger.Infof("回复消息: %d", msgt.Id)
						// fmt.Printf("回复=%+v\n", msgt.Id)
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

// TODO: 解析影视分享卡片
func (mov *movie) ServeCard(b *bot.Bot) {
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {

		msgEles := message.ToSrcProtoElems(msg.Elements)
		for i := 0; i < len(msgEles); i++ {
			var msgItem = msgEles[i].GetLightApp()
			if msgItem != nil {
				var content []byte
				if msgItem.Data[0] == 0 {
					content = msgItem.Data[1:]
				}
				if msgItem.Data[0] == 1 {
					content = binary.ZlibUncompress(msgItem.Data[1:])
				}
				if len(content) > 0 && len(content) < 1024*1024*1024 { // 解析出错 or 非法内容
					// TODO: 解析具体的APP
					fmt.Printf("message=%+v\n", string(content))

					var data_obj interface{}
					json.Unmarshal(content, &data_obj)
					if data_obj != nil {
						data := data_obj.(map[string]interface{})

						imgUrlMetaObj := data["meta"].(map[string]interface{})

						if imgUrlMetaObj != nil {
							imgUrlNewsObj := imgUrlMetaObj["news"].(map[string]interface{})

							if imgUrlNewsObj != nil {
								infoTitleStr := imgUrlNewsObj["title"].(string)
								imgUrlStr := imgUrlNewsObj["preview"].(string)
								playUrlStr := imgUrlNewsObj["jumpUrl"].(string)

								if imgFile, err := downloadLoadImage(imgUrlStr); err == nil {
									if imgObj, err := c.UploadGroupImage(msg.GroupCode, imgFile); err == nil {

										m := message.NewSendingMessage().Append(message.NewText(data["prompt"].(string)))
										m.Append(imgObj)
										m.Append(message.NewText(fmt.Sprintf("视频名称: %s\n视频地址: %s", infoTitleStr, playUrlStr)))

										c.SendGroupMessage(msg.GroupCode, m)
									} else {

										fmt.Printf("\nerr=%+v\n", err)
									}
								}

							}
						}

					}

					return
				}
			}
		}

	})
}

func (mov *movie) Start(bot *bot.Bot) {
}

func (mov *movie) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func autoreply(in string, qq string) string {

	_, _, configFail, configErr := getTemplateConfig()
	if configErr != nil {
		// configSuccess = BaseBotName + "（" + BaseWebURL + "）帮您搜索到 ${movie.total} 条《${movie.keywords}》相关内容，建议复制链接打开手机浏览器观看：\n\n${movie.list}"
		// configEmpty = BaseBotName + "（" + BaseWebURL + "）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://" + BaseWebURL + "/app"
		configFail = BaseBotName + "（" + BaseWebURL + "）好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://" + BaseWebURL + "/app"
	}

	if in == "" {
		return configFail
	}

	APIID, apiid_err := beego.AppConfig.Int("apiid")
	if apiid_err != nil {
		APIID = 0
	}

	var out string
	var ok bool
	switch APIID {
	case 1:
		// xiamaoshipin.com
		out, ok = SearchMovie(in)
	case 2:
		// xiaocaihong.tv
		out, ok = SearchMovie2(in)
	case 3:
		// juhaokan.renzaichazai.cn
		out, ok = SearchMovie3(in)
	case 4:
		// juhaokantv.com
		out, ok = SearchMovie4(in)
	case 5:
		// mangguoshipin.info
		out, ok = SearchMovie5All(in, qq)
	default:
		// 默认
		out, ok = SearchMovie4(in)
	}

	// out, ok := SearchMovie4(in)
	if !ok {
		return "搜索失败，服务器异常"
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

	configSuccess, configEmpty, _, configErr := getTemplateConfig()
	if configErr != nil {
		configSuccess = BaseBotName + "（" + BaseWebURL + "）帮您搜索到 ${movie.total} 条《${movie.keywords}》相关内容，建议复制链接打开手机浏览器观看：\n\n${movie.list}"
		configEmpty = BaseBotName + "（" + BaseWebURL + "）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://" + BaseWebURL + "/app"
		// configFail = BaseBotName + "（" + BaseWebURL + "）好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://" + BaseWebURL + "/app"
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
			// 搜索到电影了，使用成功的消息模版
			text = configSuccess
		} else {
			// 没有搜索到电影，使用空数据的消息模版
			text = configEmpty
		}

		// 将数据传递进去进行匹配
		text = replaceTemplateCharacters(text, map[string]interface{}{
			"total":    quantity,   // 搜索到的电影数量
			"keywords": keywords,   // 搜索的关键词
			"list":     seekMovies, // 电影列表
		}, "movie")
	}

	// fmt.Println("【消息】" + text)
	return text, true
}

func SearchMovie2(keywords string) (callback string, ok bool) {
	seekApi := "https://xiaocaihong.tv/api/movie/qq/search/"
	req := httplib.Post(seekApi)
	req.Param("q", keywords)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	configSuccess, configEmpty, _, configErr := getTemplateConfig()
	if configErr != nil {
		configSuccess = BaseBotName + "（" + BaseWebURL + "）帮您搜索到 ${movie.total} 条《${movie.keywords}》相关内容，建议复制链接打开手机浏览器观看：\n\n${movie.list}"
		configEmpty = BaseBotName + "（" + BaseWebURL + "）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://" + BaseWebURL + "/app"
		// configFail = BaseBotName + "（" + BaseWebURL + "）好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://" + BaseWebURL + "/app"
	}

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
			// 搜索到电影了，使用成功的消息模版
			text = configSuccess
		} else {
			// 没有搜索到电影，使用空数据的消息模版
			text = configEmpty
		}

		// 将数据传递进去进行匹配
		text = replaceTemplateCharacters(text, map[string]interface{}{
			"total":    quantity,   // 搜索到的电影数量
			"keywords": keywords,   // 搜索的关键词
			"list":     seekMovies, // 电影列表
		}, "movie")
	}

	// fmt.Println("【消息】" + text)
	return text, true
}

func SearchMovie3(keywords string) (callback string, ok bool) {
	seekApi := "http://juhaokan.renzaichazai.cn/api/movie/search?name=" + keywords
	req := httplib.Post(seekApi)
	// req.Param("q", keywords)
	// req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	configSuccess, configEmpty, _, configErr := getTemplateConfig()
	if configErr != nil {
		configSuccess = BaseBotName + "（" + BaseWebURL + "）帮您搜索到 ${movie.total} 条《${movie.keywords}》相关内容，建议复制链接打开手机浏览器观看：\n\n${movie.list}"
		configEmpty = BaseBotName + "（" + BaseWebURL + "）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://" + BaseWebURL + "/app"
		// configFail = BaseBotName + "（" + BaseWebURL + "）好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://" + BaseWebURL + "/app"
	}

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
			// 搜索到电影了，使用成功的消息模版
			text = configSuccess
		} else {
			// 没有搜索到电影，使用空数据的消息模版
			text = configEmpty
		}

		// 将数据传递进去进行匹配
		text = replaceTemplateCharacters(text, map[string]interface{}{
			"total":    quantity,   // 搜索到的电影数量
			"keywords": keywords,   // 搜索的关键词
			"list":     seekMovies, // 电影列表
		}, "movie")
	}

	// fmt.Println("【消息】" + text)
	return text, true
}

func SearchMovie4(keywords string) (callback string, ok bool) {
	seekApi := "https://juhaokantv.com/api/movie/search?name=" + url.QueryEscape(keywords)
	req := httplib.Post(seekApi).SetTimeout(100*time.Second, 30*time.Second)
	// req.Param("q", keywords)
	// req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	configSuccess, configEmpty, _, configErr := getTemplateConfig()
	if configErr != nil {
		configSuccess = BaseBotName + "（" + BaseWebURL + "）帮您搜索到 ${movie.total} 条《${movie.keywords}》相关内容，建议复制链接打开手机浏览器观看：\n\n${movie.list}"
		configEmpty = BaseBotName + "（" + BaseWebURL + "）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://" + BaseWebURL + "/app"
		// configFail = BaseBotName + "（" + BaseWebURL + "）好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://" + BaseWebURL + "/app"
	}

	// 判断请求是否成功
	if res, err := req.Response(); err != nil || res.StatusCode > 299 || res.StatusCode < 200 {
		return "", false
	}

	str, err := req.String()
	if err != nil {
		return "", false
		// t.Fatal(err)
	}

	text := ""
	seekMovies := ""

	var data_obj interface{}
	if err = json.Unmarshal([]byte(str), &data_obj); err != nil {
		return "", false
	}

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
			// 搜索到电影了，使用成功的消息模版
			text = configSuccess
		} else {
			// 没有搜索到电影，使用空数据的消息模版
			text = configEmpty
		}

		// 将数据传递进去进行匹配
		text = replaceTemplateCharacters(text, map[string]interface{}{
			"total":    quantity,   // 搜索到的电影数量
			"keywords": keywords,   // 搜索的关键词
			"list":     seekMovies, // 电影列表
		}, "movie")
	}

	// fmt.Println("【消息】" + text)
	return text, true
}

func SearchMovie5(keywords string) (callback string, ok bool) {
	seekApi := "https://mangguoshipin.info/api/movie/search?name=" + keywords
	req := httplib.Post(seekApi).SetTimeout(100*time.Second, 30*time.Second)
	// req.Param("q", keywords)
	// req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	configSuccess, configEmpty, _, configErr := getTemplateConfig()
	if configErr != nil {
		configSuccess = BaseBotName + "（" + BaseWebURL + "）帮您搜索到 ${movie.total} 条《${movie.keywords}》相关内容，建议复制链接打开手机浏览器观看：\n\n${movie.list}"
		configEmpty = BaseBotName + "（" + BaseWebURL + "）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://" + BaseWebURL + "/app"
		// configFail = BaseBotName + "（" + BaseWebURL + "）好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://" + BaseWebURL + "/app"
	}

	// 判断请求是否成功
	if res, err := req.Response(); err != nil || res.StatusCode > 299 || res.StatusCode < 200 {
		return "", false
	}

	str, err := req.String()
	if err != nil {
		return "", false
		// t.Fatal(err)
	}

	text := ""
	seekMovies := ""

	var data_obj interface{}
	if err = json.Unmarshal([]byte(str), &data_obj); err != nil {
		return "", false
	}

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
			// 搜索到电影了，使用成功的消息模版
			text = configSuccess
		} else {
			// 没有搜索到电影，使用空数据的消息模版
			text = configEmpty
		}

		// 将数据传递进去进行匹配
		text = replaceTemplateCharacters(text, map[string]interface{}{
			"total":    quantity,   // 搜索到的电影数量
			"keywords": keywords,   // 搜索的关键词
			"list":     seekMovies, // 电影列表
		}, "movie")
	}

	// fmt.Println("【消息】" + text)
	return text, true
}

// TODO: 为了满足业务需求，增加一个聚合搜索接口，会传 QQ 账号给后端，请确保后端安全避免泄漏数据
func SearchMovie5All(keywords string, qq string) (callback string, ok bool) {
	seekApi := "http://mangguoshipin.neihancloud.com/api/movie/search?qq=" + qq + "&name=" + url.QueryEscape(keywords)
	req := httplib.Post(seekApi).SetTimeout(100*time.Second, 30*time.Second)
	// req.Param("q", keywords)
	// req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	configSuccess, configEmpty, _, configErr := getTemplateConfig()
	if configErr != nil {
		configSuccess = BaseBotName + "（" + BaseWebURL + "）帮您搜索到 ${movie.total} 条《${movie.keywords}》相关内容，建议复制链接打开手机浏览器观看：\n\n${movie.list}"
		configEmpty = BaseBotName + "（" + BaseWebURL + "）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://" + BaseWebURL + "/app"
		// configFail = BaseBotName + "（" + BaseWebURL + "）好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://" + BaseWebURL + "/app"
	}

	// 判断请求是否成功
	if res, err := req.Response(); err != nil || res.StatusCode > 299 || res.StatusCode < 200 {
		return "", false
	}

	str, err := req.String()
	if err != nil {
		return "", false
		// t.Fatal(err)
	}

	text := ""
	seekMovies := ""

	var data_obj interface{}
	if err = json.Unmarshal([]byte(str), &data_obj); err != nil {
		return "", false
	}

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
			// 搜索到电影了，使用成功的消息模版
			text = configSuccess
		} else {
			// 没有搜索到电影，使用空数据的消息模版
			text = configEmpty
		}

		// 将数据传递进去进行匹配
		text = replaceTemplateCharacters(text, map[string]interface{}{
			"total":    quantity,   // 搜索到的电影数量
			"keywords": keywords,   // 搜索的关键词
			"list":     seekMovies, // 电影列表
		}, "movie")
	}

	// fmt.Println("【消息】" + text)
	return text, true
}

/**
 * @description: 获取缓存的模版数据
 * @param {*}
 * @return {*}
 */
func getTemplateConfig() (configSuccess string, configEmpty string, configFail string, err error) {
	config, err := helper.GetConfigsDataByName("plugin_config_message_template")
	if err != nil || config == nil {
		config = map[string]interface{}{
			"success": BaseBotName + "（" + BaseWebURL + "）帮您搜索到 ${movie.total} 条《${movie.keywords}》相关内容，建议复制链接打开手机浏览器观看：\n\n${movie.list}",                         // 成功的消息模版
			"empty":   BaseBotName + "（" + BaseWebURL + "）很遗憾暂时没有搜索到相关内容，资源马上就上线下载 APP 看看？立即下载：https://" + BaseWebURL + "/app",                                     // 空的消息模版
			"fail":    BaseBotName + "（" + BaseWebURL + "）好像不知道到您想要搜索的关键词，试试热门搜索：\n1，流浪地球\n2，你的名字\n3，我和我的祖国永远在一起\n\n下载APP高清资源无限免费看：https://" + BaseWebURL + "/app", // 失败的消息模版
		}
	}
	// configSuccess, err = strconv.Unquote("\"" + config["success"].(string) + "\"")
	// configEmpty, err = strconv.Unquote("\"" + config["empty"].(string) + "\"")
	// configFail, err = strconv.Unquote("\"" + config["fail"].(string) + "\"")
	configSuccess = config["success"].(string)
	configEmpty = config["empty"].(string)
	configFail = config["fail"].(string)
	return
}

/**
 * @description: 消息模版数据匹配方法
 * @param {string} template
 * @param {map[string]interface{}} data
 * @param {string} flag
 * @return {*}
 */
func replaceTemplateCharacters(template string, data map[string]interface{}, flag string) string {

	// 数据判空
	if template == "" || data == nil || flag == "" {
		return template
	}

	// 正则获取需要替换的参数
	flysnowRegexp := regexp.MustCompile(`\$\{` + flag + `\.(.+?)\}`)
	params := flysnowRegexp.FindAllStringSubmatch(template, -1)
	if params == nil || len(params) <= 0 {
		return template
	}

	for _, param := range params {
		if param == nil || len(param) < 2 || param[1] == "" {
			// 匹配到的关键词异常
			continue
		}
		paramKey := param[1]

		paramValue := ""
		// interface{} 数据转 string
		if data[paramKey] != nil {
			tValue := data[paramKey]
			switch tValue.(type) {
			case string:
				paramValue = tValue.(string)
				break
			case int:
				paramValue = strconv.Itoa(tValue.(int))
				break
			case int64:
				paramValue = strconv.FormatInt(tValue.(int64), 10)
				break
			case float64:
				strconv.FormatFloat(tValue.(float64), 'E', -1, 32)
				break
			case bool:
				strconv.FormatBool(tValue.(bool))
				break
			}
		}

		template = regexp.MustCompile(`\$\{`+flag+`\.`+paramKey+`\}`).ReplaceAllLiteralString(template, paramValue)
	}

	return template
}

func downloadLoadImage(url string) (img io.ReadSeeker, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	file, err := ioutil.ReadAll(resp.Body)
	return bytes.NewReader(file), err
}
