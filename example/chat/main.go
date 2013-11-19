package main

import (
	"github.com/linphy/beego"
	"github.com/linphy/beego/example/chat/controllers"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.WSController{})
	beego.Run()
}
