package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	//"strconv"
	"strings"
)

var lastcpmmand string
var ptyfile *os.File

// Operations about Users
type TeminalController struct {
	beego.Controller
}

var containerid = "null"
var baseimage = "null"

var wsmap_term = make(map[string]*websocket.Conn)

// @Title render main page
// @Description : start the websocket connection
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [get]
func (o *TeminalController) Getpage() {
	fmt.Println("test")
	t, err := template.ParseFiles("views/terminal.html")
	if err != nil {
		panic(err)
	}
	t.Execute(o.Ctx.ResponseWriter, nil)
	//some err with template render in beego??
	//o.TplNames = "views/terminal.html"
}

// @Title check if the image is exist
// @Description : check the image info
// @router /checkimage [post]
func (o *TeminalController) Check() {
	imagename := (o.Ctx.Request.Form["imagename"])[0]
	fmt.Println("check result:", imagename)
	//get images/(name)/json
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	_, err := client.InspectImage(imagename)
	if err != nil {
		fmt.Println(err.Error())
		o.Ctx.WriteString("the image not exist locallly")
		return
	}
	o.Ctx.WriteString("ok")

}

// send break signal
// ps -ax |grep -v grep |grep "bee run"| awk '{print $1} |head -1'
// or pgrep

// @Title send the break signal
// @Description : send the break signal
// @router /break [get]
// problems in sending ctrl+c
func (o *TeminalController) Break() {
	fmt.Println(containerid)
	if containerid == "null" {
		fmt.Println("contianer not exist")
		return
	}

	//get the pid of /bin/bash
	//command := `ps -ef|grep -v grep |grep "sudo docker run "`
	file, err := os.Open("/var/lib/docker/containers/")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(file.Name())
	//find the correct name by regex
	tmpregxp := containerid + ".*"
	fmt.Println("regx:", tmpregxp)
	//regxp := "a66b9fb984.*"
	names, _ := file.Readdirnames(0)
	fullimagename := "null"
	for _, name := range names {
		m, _ := regexp.MatchString(tmpregxp, name)
		if m {
			fullimagename = name
			break
		}
	}
	fmt.Println(fullimagename)

	fullpath := "/var/lib/docker/containers/" + fullimagename + "/config.json"

	fmt.Println("filepath:", fullpath)

	dat, err := ioutil.ReadFile(fullpath)
	if err != nil {
		fmt.Println(err.Error())
	}

	js, err := simplejson.NewJson(dat)

	if err != nil {
		fmt.Println(err.Error())
	}

	var configmap = make(map[string]interface{})
	configmap, _ = js.Map()
	//fmt.Println(configmap)
	containerpid := configmap["State"].(map[string]interface{})["Pid"]
	fmt.Println(reflect.TypeOf(containerpid))
	fmt.Println(containerpid)
	containerpid_int, _ := containerpid.(json.Number).Int64()

	containerpid_str := fmt.Sprintf("%d", containerpid_int)

	termcommand := `kill 2 $(ps -ef|grep -v grep |grep "` + lastcpmmand + `" |grep ` + containerpid_str + `  |awk '{print $2}')`
	system(termcommand)
}

// @Title testterm
// @Description : start the websocket connection
// @Param	body		body 	models.User	true		"body for user content"
// @router /:baseimage [get]
func (o *TeminalController) Get() {
	baseimage = o.GetString(":baseimage")
	//if baseimage == "" {
	//	http.Error(o.Ctx.ResponseWriter, "null image id", 400)
	//	return
	//}

	endpoint := o.Ctx.Request.RemoteAddr

	url := strings.Split(endpoint, ":")[0]
	fmt.Println(url)
	ws, err := websocket.Upgrade(o.Ctx.ResponseWriter, o.Ctx.Request, nil, 1024, 1024)
	wsmap_term[url] = ws
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(o.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}
	o.Ctx.WriteString("connection ok")

	//start the pty
	// ubuntu:latest
	//c := exec.Command("docker", "run", "-it", "ubuntu:latest", "/bin/bash")
	c := exec.Command("docker", "run", "-it", baseimage, "/bin/bash")
	//c := exec.Command("/bin/bash")
	//	pty.Open()
	f, err := pty.Start(c)
	ptyfile = f
	if err != nil {
		panic(err)
	}
	//pipeReader, pipeWriter := io.Pipe()
	wsm := wsmap_term[url]
	go func() {

		for {
			_, p, err := wsm.ReadMessage()

			if err != nil {
				panic(err)
			}

			//write the command into the terminal

			fmt.Println("input command:", string(p))

			//store the last command
			lastcpmmand = string(p)

			p = append(p, 10)
			io.Copy(f, strings.NewReader(string(p)))

		}
	}()
	//it's ok to redirect the output to the stdout
	//io.Copy(os.Stdout, f)
	getid := false
	go func() {
		//attention the position to create the newreader
		r := bufio.NewReader(f)
		for {

			line, _, err := r.ReadLine()
			if err != nil {
				break
			}

			if strings.Contains(string(line), "@") {
				if getid == false {
					str1 := strings.Split(string(line), "@")[1]
					str2 := strings.Split(str1, ":")[0]
					containerid = str2
					getid = true
				}
				continue
			}

			wsm.WriteMessage(websocket.TextMessage, line)

		}
	}()

}
