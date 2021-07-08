package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

func (c *AdminController) Get() {
	_, _, is := isLogin(c)

	flash := beego.NewFlash()
	if _, ok := flash.Data["notice"]; ok {
		c.TplName = "admin.tpl"
		return
	}

	if !is {
		// token 获取失败或失效，或用户被禁用将跳转登陆
		c.Redirect("/login", 301)
		c.Finish()
		return
	}

	c.TplName = "admin.tpl"
}

// AdminController operations for Admin
type AdminController struct {
	beego.Controller
}

func (c *AdminController) LoginPage() {
	token, _, is := isLogin(c)
	if is {

		flash := beego.NewFlash()
		flash.Notice(token)
		flash.Store(&c.Controller)

		c.Redirect("/admin", 301)
		return
	}

	c.TplName = "login.tpl"

}
