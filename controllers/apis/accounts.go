package controllers

import (
	"ferry_ship/bot"
	"ferry_ship/models"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Mrs4s/MiraiGo/client"
	beego "github.com/beego/beego/v2/server/web"
)

type AccountsController struct {
	beego.Controller
}

/**
 * @description: 添加一个 QQ 账号
 * @param {*}
 * @return {*}
 */
func (c *AccountsController) ApiAddBotAccount() {

	_, user, _ := userAssistant(&c.Controller)

	u_account, _ := c.GetInt64("account")
	u_password := c.GetString("password")
	u_auto, _ := c.GetInt("auto")

	if u_account == 0 || u_password == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	account := &models.Accounts{
		Account:   u_account,
		Password:  u_password,
		AutoLogin: u_auto,
		User:      &user,
	}

	account_, existing := models.GetAccountsByAccount(u_account)
	if account_ != nil && existing == nil {

		// go models.AccountLoginQQ(account_)
		callBackResult(&c.Controller, 200, "账号添加失败，账号已存在", nil)

		flash := beego.NewFlash()
		flash.Notice("aaaaa")
		flash.Store(&c.Controller)
		c.Finish()
		return
	}

	id, err := models.AddAccounts(account)
	account, err = models.GetAccountsById(int(id))

	if err != nil {
		callBackResult(&c.Controller, 200, "账号添加失败，请稍后重试", nil)
		c.Finish()
		return
	}

	c.Data["json"] = models.TurnAccountsToMap(account)

	callBackResult(&c.Controller, 200, "", c.Data["json"])
	c.Finish()
}

/**
 * @description: 删除一个机器人 QQ 账号
 * @param {*}
 * @return {*}
 */
func (c *AccountsController) ApiDeleteBotAccount() {
	userAssistant(&c.Controller)

	botID, _ := c.GetInt("id")
	if botID == 0 {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	account_, existing := models.GetAccountsById(botID)
	if account_ == nil && existing != nil {
		callBackResult(&c.Controller, 200, "账号删除失败，账号不存在", nil)
		c.Finish()
		return
	}

	existing = models.DeleteAccounts(account_.Id)
	if existing != nil {
		callBackResult(&c.Controller, 200, "账号删除失败，出现异常", nil)
		c.Finish()
		return
	}

	callBackResult(&c.Controller, 200, "", models.TurnAccountsToMap(account_))
}

// 响应获取全部机器人账号列表
func (c *AccountsController) ApiGetAllAccount() {
	userAssistant(&c.Controller) // 登陆认证

	u_count, _ := c.GetInt("count", 10)
	u_page, _ := c.GetInt("page", 0)

	// 刷新机器人在线状态列表数据
	models.RefreshAccountBotInfo()

	accounts, err := models.AllAccounts(u_count, u_page)

	if err != nil {
		callBackResult(&c.Controller, 403, "服务器错误", nil)
		c.Finish()
		return
	}

	var new_accouts []interface{}

	for item := range accounts {
		i_bot := accounts[item]
		new_bot := models.TurnAccountsToMap(&i_bot)
		new_accouts = append(new_accouts, new_bot)
	}

	callBackResult(&c.Controller, 200, "", new_accouts)
}

// 认证机器人账号登陆滑块
func (c *AccountsController) ApiBotVerifyTicket() {
	userAssistant(&c.Controller) // 认证

	u_ticket := c.GetString("ticket")

	if u_ticket == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	if bot.Instance == nil {
		callBackResult(&c.Controller, 200, "没有待认证的账号", nil)
		c.Finish()
		return
	}

	resp, err := bot.Instance.SubmitTicket(u_ticket)

	if !resp.Success || err != nil {

		if err != nil {
			callBackResult(&c.Controller, 200, "登陆出错，"+err.Error(), nil)
		} else {
			c.Data["json"] = botCallBackToMap(resp)
			callBackResult(&c.Controller, 200, "登陆出错，"+resp.ErrorMessage, c.Data["json"])
		}

		c.Finish()
		return
	}

	c.Data["json"] = ""
	callBackResult(&c.Controller, 200, "成功。", c.Data["json"])
	c.Finish()
}

// 登陆机器人账号
func (c *AccountsController) ApiLoginBotAccount() {

	userAssistant(&c.Controller) // 认证

	a_id, err := c.GetInt("id")

	account, err := models.GetAccountsById(a_id)
	if err != nil || account == nil {
		// 账号不存在
		callBackResult(&c.Controller, 200, "该账号不存在", nil)
		c.Finish()
		return
	}

	// 初始化 Bot
	bot.InitBot(account.Account, account.Password)

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.IPad)

	// 登录
	resp, err := bot.Instance.Login()
	// console := bufio.NewReader(os.Stdin)

	for {
		if err != nil {
			// logger.WithError(err).Fatal("unable to login")
			callBackResult(&c.Controller, 200, "QQ 账号登陆异常，"+err.Error(), nil)
			c.Finish()
			return
		}

		if !resp.Success {
			switch resp.Error {

			case client.NeedCaptcha:
				// img, _, _ := image.Decode(bytes.NewReader(resp.CaptchaImage))

				// fmt.Println(asc2art.New("image", img).Art)

				// fmt.Printf("please input captcha (%x) : \n", resp.CaptchaSign)
				// fmt.Printf("please input captcha (%s) : \n", convert(resp.CaptchaSign))

				// for text == "" {
				// 	// text, _ = console.ReadString('\n')
				// 	log.Println("待输入验证码")
				// 	if text != "" {
				// 		log.Println("输入的验证码：" + text)
				// 		break
				// 	}
				// }
				// log.Println("输入的验证码：" + text)
				// resp, err = bot.Instance.SubmitCaptcha(strings.ReplaceAll(text, "\n", ""), resp.CaptchaSign)

				err := ioutil.WriteFile("static/img/acptcha/log.jpg", resp.CaptchaImage, os.FileMode(0755))
				if err != nil {
					c.Data["json"] = map[string]interface{}{
						"error": 10010,
						"text":  "(验证码获取失败) login failed",
					}
					callBackResult(&c.Controller, 200, "", c.Data["json"])
					c.Finish()
					return
				}

				c.Data["json"] = map[string]interface{}{
					"error": 10011,
					"text":  "(登陆需要验证码) login failed",
					"url":   "/static/img/acptcha/log.jpg",
					"sign":  resp.CaptchaSign,
				}
				callBackResult(&c.Controller, 200, "", c.Data["json"])
				c.Finish()
				return

				continue
			case client.UnsafeDeviceError:

				// 不安全设备错误
				// fmt.Printf("device lock -> %v\n", resp.VerifyUrl)
				// os.Exit(4)

				c.Data["json"] = map[string]interface{}{
					"error": 10020,
					"text":  "(不安全设备错误) login failed",
					"url":   resp.VerifyUrl,
				}
				callBackResult(&c.Controller, 200, "", c.Data["json"])
				c.Finish()
				return

				continue
			case client.SMSNeededError:

				// 需要SMS错误
				// fmt.Println("device lock enabled, Need SMS Code")
				// fmt.Printf("Send SMS to %s ? (yes)", resp.SMSPhone)
				// t, _ := console.ReadString('\n')
				// t = strings.TrimSpace(t)
				// if t != "yes" {
				// os.Exit(2)
				// }
				// if !bot.Instance.RequestSMS() {
				// logger.Warnf("unable to request SMS Code")
				// log.Println("unable to request SMS Code")
				// os.Exit(2)
				// }
				// logger.Warn("please input SMS Code: ")
				// fmt.Println("please input SMS Code: ")
				// text, _ = console.ReadString('\n')
				// resp, err = bot.Instance.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))

				c.Data["json"] = map[string]interface{}{
					"error": 10030,
					"text":  "(需要短信验证码) login failed",
				}
				callBackResult(&c.Controller, 200, "", c.Data["json"])
				c.Finish()
				return

				continue
			case client.TooManySMSRequestError:
				// 短信请求错误太多
				// fmt.Printf("too many SMS request, please try later.\n")
				// os.Exit(6)

				c.Data["json"] = map[string]interface{}{
					"error": 10040,
					"text":  "(短信请求错误太多) login failed",
				}
				callBackResult(&c.Controller, 200, "", c.Data["json"])
				c.Finish()
				return

				continue
			case client.SMSOrVerifyNeededError:
				// SMS或验证所需的错误

				// fmt.Println("device lock enabled, choose way to verify:")
				// fmt.Println("1. Send SMS Code to ", resp.SMSPhone)
				// fmt.Println("2. Scan QR Code")
				// fmt.Print("input (1,2):")
				// text, _ = console.ReadString('\n')
				// text = strings.TrimSpace(text)
				// switch text {
				// case "1":
				// 	if !Instance.RequestSMS() {
				// 		fmt.Println("unable to request SMS Code")
				// 		os.Exit(2)
				// 	}
				// 	fmt.Print("please input SMS Code: ")
				// 	text, _ = console.ReadString('\n')
				// 	resp, err = Instance.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))
				// 	continue
				// case "2":
				// 	fmt.Printf("device lock -> %v\n", resp.VerifyUrl)
				// 	os.Exit(2)
				// default:
				// 	fmt.Println("invalid input")
				// 	os.Exit(2)
				// }

				c.Data["json"] = map[string]interface{}{
					"error": 10050,
					"text":  "(需要短信验证码或扫描二维码) login failed",
					"url":   resp.VerifyUrl,
				}
				callBackResult(&c.Controller, 200, "", c.Data["json"])
				c.Finish()
				return

				continue
			case client.SliderNeededError:
				// 需要滑动认证
				// fmt.Println("please look at the doc https://github.com/Mrs4s/go-cqhttp/blob/master/docs/slider.md to get ticket")
				// fmt.Printf("open %s to get ticket\n", resp.VerifyUrl)
				// fmt.Printf("等待滑动认证…")
				// resp, err = bot.Instance.SubmitTicket(VerifyTicket)

				c.Data["json"] = map[string]interface{}{
					"error": 10060,
					"text":  "(需要滑动认证) please look at the doc https://github.com/Mrs4s/go-cqhttp/blob/master/docs/slider.md to get ticket",
					"url":   resp.VerifyUrl,
				}
				callBackResult(&c.Controller, 200, "", c.Data["json"])
				c.Finish()
				return

				continue
			case client.OtherLoginError, client.UnknownLoginError:
				// 其他登录错误
				// logger.Fatalf("login failed: %v", resp.ErrorMessage)

				c.Data["json"] = map[string]interface{}{
					"error": 10070,
					"text":  "(其他登陆错误) login failed: " + resp.ErrorMessage,
				}
				callBackResult(&c.Controller, 200, "", c.Data["json"])
				c.Finish()
				return
			}
		} else {
			// 刷新好友列表，群列表
			bot.RefreshList()

			// 将登陆成功的对象加入序列
			bot.Instances[account.Account] = bot.Instance

			// 返回数据
			c.Data["json"] = models.TurnAccountsToMap(account)
			callBackResult(&c.Controller, 200, "QQ 账号登陆成功", c.Data["json"])
			c.Finish()

			// 刷新全部机器人账号信息
			models.RefreshAccountBotInfo()
			return
		}
		break
	}
}

// 获取机器人账号信息
func (c *AccountsController) ApiGetBotInfo() {
	userAssistant(&c.Controller) // 认证

	// 调试 flash
	flash := beego.ReadFromRequest(&c.Controller)
	if n, ok := flash.Data["notice"]; ok {
		fmt.Println("输出: " + n)
	}

	// helper.SearchMovie("黑猫")

	// u_account, _ := c.GetInt64("account")
	// bot := helper.GetBotInfo(u_account)

	// if bot != nil {

	// 	o := orm.NewOrm()

	// 	online := 0
	// 	if bot.Online {
	// 		online = 1
	// 	}

	// 	num, err := o.QueryTable("accounts").Filter("account", u_account).Update(orm.Params{
	// 		"name":   bot.Nickname,
	// 		"avatar": "https://q2.qlogo.cn/headimg_dl?spec=100&dst_uin=" + strconv.FormatInt(bot.Uin, 10),
	// 		"status": online,
	// 	})

	// 	fmt.Println(num)

	// 	if err != nil {
	// 		fmt.Println("账号不存在")
	// 		fmt.Println(strconv.Itoa(1) + "更新失败" + err.Error())

	// 	}

	// } else {
	// 	fmt.Println("bot 不存在")
	// }

	callBackResult(&c.Controller, 200, "", c.Data["json"])
	c.Finish()
}

// 重设机器人账号密码
func (c *AccountsController) ApiUpdateBotPassword() {
	userAssistant(&c.Controller) // 认证

	a_id, err := c.GetInt("id")
	a_password := c.GetString("password")

	if a_id == 0 || a_password == "" {
		callBackResult(&c.Controller, 200, "缺失参数", nil)
		c.Finish()
		return
	}

	account, err := models.GetAccountsById(a_id)
	if err != nil || account == nil {
		// 账号不存在
		callBackResult(&c.Controller, 200, "该账号不存在", nil)
		c.Finish()
		return
	}

	if account, ok := models.AccountReupdatePassword(account, a_password); ok && account != nil {
		c.Data["json"] = models.TurnAccountsToMap(account)
		callBackResult(&c.Controller, 200, "", c.Data["json"])
		c.Finish()
		return
	}

	callBackResult(&c.Controller, 200, "账号密码修改失败，请稍后重试", nil)
	c.Finish()
	return
}

// 下线机器人账号
func (c *AccountsController) ApiLogoutBotAccount() {
	userAssistant(&c.Controller) // 认证

	a_id, err := c.GetInt("id")

	account, err := models.GetAccountsById(a_id)
	if err != nil || account == nil {
		// 账号不存在
		callBackResult(&c.Controller, 200, "该账号不存在", nil)
		c.Finish()
		return
	}

	botObj := bot.Instances[account.Account]

	// fmt.Printf("账号=%+v  昵称=%+v  状态=%+v\n", account.Account, botObj.Nickname)
	if botObj == nil {
		// 账号未登陆
		callBackResult(&c.Controller, 200, "该账号未登陆", nil)
		// fmt.Printf("账号=%+v\n", botObj.Nickname)
		// 刷新机器人在线状态列表数据
		models.RefreshAccountBotInfo()
		c.Finish()
		return
	}

	// 从在线列表删除
	delete(bot.Instances, account.Account)
	// 设置状态为离线
	botObj.QQClient.SetOnlineStatus(client.StatusOffline)
	// 刷新机器人在线状态列表数据
	models.RefreshAccountBotInfo()

	// 返回退出登陆成功结果
	c.Data["json"] = account
	callBackResult(&c.Controller, 200, "", c.Data["json"])
	c.Finish()
	return
}

// 转换登陆结果为 map 数据
func botCallBackToMap(resp *client.LoginResponse) map[string]interface{} {
	switch resp.Error {

	case client.NeedCaptcha:

		err := ioutil.WriteFile("static/img/acptcha/log.jpg", resp.CaptchaImage, os.FileMode(0755))
		if err != nil {
			return map[string]interface{}{
				"error": 10010,
				"text":  "(验证码获取失败) login failed",
			}
		}

		return map[string]interface{}{
			"error": 10011,
			"text":  "(登陆需要验证码) login failed",
			"url":   "/static/img/acptcha/log.jpg",
			"sign":  resp.CaptchaSign,
		}

	case client.UnsafeDeviceError:
		// 不安全设备错误
		return map[string]interface{}{
			"error": 10020,
			"text":  "(不安全设备错误) login failed",
			"url":   resp.VerifyUrl,
		}
	case client.SMSNeededError:

		// 需要SMS错误
		return map[string]interface{}{
			"error": 10030,
			"text":  "(需要短信验证码) login failed",
		}

	case client.TooManySMSRequestError:
		// 短信请求错误太多

		return map[string]interface{}{
			"error": 10040,
			"text":  "(短信请求错误太多) login failed",
		}

	case client.SMSOrVerifyNeededError:
		// SMS或验证所需的错误

		return map[string]interface{}{
			"error": 10050,
			"text":  "(需要短信验证码或扫描二维码) login failed",
			"url":   resp.VerifyUrl,
		}

	case client.SliderNeededError:
		// 需要滑动认证

		return map[string]interface{}{
			"error": 10060,
			"text":  "(需要滑动认证) please look at the doc https://github.com/Mrs4s/go-cqhttp/blob/master/docs/slider.md to get ticket",
			"url":   resp.VerifyUrl,
		}

	case client.OtherLoginError, client.UnknownLoginError:
		// 其他登录错误

		return map[string]interface{}{
			"error": 10070,
			"text":  "(其他登陆错误) login failed: " + resp.ErrorMessage,
		}

	}

	return nil
}
