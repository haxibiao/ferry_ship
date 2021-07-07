package main

import (
	_ "ferry_ship/models"
	_ "ferry_ship/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}
