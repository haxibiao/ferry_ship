package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

func (c *AdminController) Get() {
	_, _, is := isLogin(c)

	if !is {
		// token 获取失败或失效，或用户被禁用将跳转登陆
		c.Redirect("/login", 301)
		c.Finish()
	}

	c.TplName = "admin.tpl"
}

// AdminController operations for Admin
type AdminController struct {
	beego.Controller
}

func (c *AdminController) LoginPage() {
	_, _, is := isLogin(c)
	if is {
		c.Redirect("/admin", 301)
	} else {
		c.TplName = "login.tpl"
	}
}
