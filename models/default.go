/*
 * @Author: Bin
 * @Date: 2021-07-07
 * @FilePath: /ferry_ship/models/default.go
 */
package models

import (
	"ferry_ship/bot"
	"ferry_ship/helper"
	"fmt"
	"sync"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

var singleOrmInstance orm.Ormer
var once sync.Once

func GetSharedOrmer() orm.Ormer {
	once.Do(func() {
		singleOrmInstance = orm.NewOrm()
	})
	return singleOrmInstance
}

func addUser(name string, password string) {
	user, err := AddUsers(&Users{Name: name, Password: helper.StringToMd5(password)})
	if err == nil {
		fmt.Println("用户创建成功: " + user.Name)
	}
}

func init() {
	tag := "【Model.go】"

	driver, _ := helper.Env("dbdriver")
	username, _ := helper.Env("dbusername")
	password, _ := helper.Env("dbpassword")
	database, _ := helper.Env("dbdatabase")
	host, _ := helper.Env("dbhost")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	connection_url := helper.GetConnectionURL(username, password, host, database)

	// fmt.Println(tag + "连接URL是: " + connection_url)

	orm.RegisterDataBase("default", driver, connection_url)

	fmt.Println(tag + "注册数据模型")
	orm.RegisterModel(
		new(Configs),     // 配置
		new(Users),       // 用户
		new(Accounts),    // QQ账号
		new(Rules),       // 回复规则
		new(Apis),        // 请求接口
		new(RulesToApis), // 规则和接口多对多关系表
	)

	// 第二个参数为 true 则强制重新建表
	orm.RunSyncdb("default", false, true)

	// 添加默认用户 admin
	d_user, d_u_err := GetUserById(1)
	if d_u_err != nil || d_user == nil {
		addUser("admin", "admin")
		fmt.Println(tag + "注册默认用户 admin 成功！")
	}

	// 刷新全部机器人账号信息
	RefreshAccountBotInfo()

	// 机器人加载搜索电影组件
	bot.RegisterModule(helper.MovieInstance)

}
