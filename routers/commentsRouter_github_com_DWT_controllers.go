package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
	beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"] = append(beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"],
		beego.ControllerComments{
			"Getpage",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"] = append(beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"],
		beego.ControllerComments{
			"Check",
			`/checkimage`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"] = append(beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"],
		beego.ControllerComments{
			"Break",
			`/break`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"] = append(beego.GlobalControllerRouter["github.com/DWT/controllers:TeminalController"],
		beego.ControllerComments{
			"Get",
			`/:baseimage`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["github.com/DWT/controllers:CdfController"] = append(beego.GlobalControllerRouter["github.com/DWT/controllers:CdfController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

}
