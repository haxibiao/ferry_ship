package routers

import (
	"ferry_ship/controllers"
	apis "ferry_ship/controllers/apis"

	web "github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.MainController{})

	web.Router("/admin", &controllers.AdminController{})
	web.Router("/login", &controllers.AdminController{}, "get:LoginPage")

	api := web.NewNamespace("/api",
		// 控制台登陆
		web.NSRouter("/login", &apis.UsersController{}, "get:ApiLogin"),

		// 用户相关 API
		web.NSNamespace("/user",
			// 获取用户信息
			web.NSRouter("/me", &apis.UsersController{}, "get:ApiGetMe"),
			// web.NSRouter("/create", &apis.UsersController{}, "post:ApiCreateUser"),
			web.NSRouter("/upstatus", &apis.UsersController{}, "post:ApiUpStatusUser"),
			web.NSRouter("/update", &apis.UsersController{}, "post:ApiUpdateUser"),
			web.NSRouter("/list", &apis.UsersController{}, "get:ApiUserList"),
		),

		// QQ账号相关 API
		web.NSNamespace("/account",
			// 获取用户信息
			web.NSRouter("/add", &apis.AccountsController{}, "post:ApiAddBotAccount"),
			web.NSRouter("/delete", &apis.AccountsController{}, "post:ApiDeleteBotAccount"),
			web.NSRouter("/getinfo", &apis.AccountsController{}, "get:ApiGetBotInfo"),
			web.NSRouter("/list", &apis.AccountsController{}, "get:ApiGetAllAccount"),
			web.NSRouter("/bot/login", &apis.AccountsController{}, "post:ApiLoginBotAccount"),
			web.NSRouter("/bot/logout", &apis.AccountsController{}, "post:ApiLogoutBotAccount"),
			web.NSRouter("/bot/ticket", &apis.AccountsController{}, "post:ApiBotVerifyTicket"),
			web.NSRouter("/bot/repassword", &apis.AccountsController{}, "post:ApiUpdateBotPassword"),
		),

		// 系统配置相关 API
		web.NSNamespace("/system",
			// 获取用户信息
			web.NSRouter("/message_template/save", &apis.SystemController{}, "post:ApiSaveMessageTemplate"),
			web.NSRouter("/message_template/get", &apis.SystemController{}, "get:ApiGetMessageTemplate"),
		),
	)

	web.AddNamespace(api)
}
