package controllers

import (
	"ferry_ship/helper"
	"ferry_ship/models"
	"fmt"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
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

	if u_auto == 1 || user.Name == "" {

	}

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

// 认证机器人账号登陆滑块
func (c *AccountsController) ApiBotVerifyTicket() {

	u_ticket := c.GetString("ticket")

	helper.VerifyTicket = u_ticket

	callBackResult(&c.Controller, 200, "", c.Data["json"])
	c.Finish()
}

func (c *AccountsController) ApiGetBotInfo() {
	flash := beego.ReadFromRequest(&c.Controller)

	if n, ok := flash.Data["notice"]; ok {
		fmt.Println("输出: " + n)
	}

	u_account, _ := c.GetInt64("account")
	bot := helper.GetBotInfo(u_account)

	if bot != nil {

		o := orm.NewOrm()

		online := 0
		if bot.Online {
			online = 1
		}

		num, err := o.QueryTable("accounts").Filter("account", u_account).Update(orm.Params{
			"name":   bot.Nickname,
			"avatar": "https://q2.qlogo.cn/headimg_dl?spec=100&dst_uin=" + strconv.FormatInt(bot.Uin, 10),
			"status": online,
		})

		fmt.Println(num)

		if err != nil {
			fmt.Println("账号不存在")
			fmt.Println(strconv.Itoa(1) + "更新失败" + err.Error())

		}

	} else {
		fmt.Println("bot 不存在")
	}

	callBackResult(&c.Controller, 200, "", c.Data["json"])
	c.Finish()
}
