package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
	beego.GlobalControllerRouter["github.com/DWT/controllers:CdfController"] = append(beego.GlobalControllerRouter["github.com/DWT/controllers:CdfController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"] = append(beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"],
		beego.ControllerComments{
			"Getpage",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"] = append(beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"],
		beego.ControllerComments{
			"Get",
			`/:baseimage`,
			[]string{"get"},
			nil})

}
