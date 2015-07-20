package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var logfile *os.File
var logger *log.Logger
var target string
var mode, static_version, hotfix_version, server_id, ip, db string

//var ips,dbs []string

type Configfile struct {
	filestring string
	ips        []string
	dbs        []string
}

func (filedata *Configfile) addnew_line(path, server_id string) {
	logger.Println("===The new server id is", server_id)
	logger.Println("===The Gip of this server is ", ip)
	server_idint, err := strconv.Atoi(server_id)
	if err != nil {
		logger.Println(err.Error())
		os.Exit(1)
	}
	server_id1 := server_idint - 4
	server_id2 := strconv.Itoa(server_id1)

	var a1, a2 string
	argv1 := []string{"-n", "2p", path}
	a := exec.Command("sed", argv1...)
	argv2 := []string{"-n", "2,10000p", path}
	b := exec.Command("sed", argv2...)
	argv3 := []string{"-n", "1p", path}
	c := exec.Command("sed", argv3...)
	d1, _ := a.Output()
	d2, _ := b.Output()
	d3, _ := c.Output()

	c1 := string(d1)
	c2 := string(d2)
	c3 := string(d3)

	filedata1, _ := regexp.Compile("nba_redis[0-9]+")
	n1 := "nba_redis" + server_id
	a1 = filedata1.ReplaceAllString(c1, n1)

	filedata2, _ := regexp.Compile("{_id:[0-9]*")
	n2 := "{_id:" + server_id
	a2 = filedata2.ReplaceAllString(a1, n2)

	filedata3, _ := regexp.Compile("game:[^]]*]")
	lenofip := len(filedata.ips)
	n3 := "game:["
	for k, v := range filedata.ips {

		if k != lenofip-1 {
			n3 = n3 + "{h:'" + v + "',p:8601}," + "{h:'" + v + "',p:8602}," + "{h:'" + v + "',p:8603}," + "{h:'" + v + "',p:8604},"
		} else if k == lenofip-1 {
			n3 = n3 + "{h:'" + v + "',p:8601}," + "{h:'" + v + "',p:8602}," + "{h:'" + v + "',p:8603}," + "{h:'" + v + "',p:8604}"

		}
	}
	n4 := n3 + "]"
	a3 := filedata3.ReplaceAllString(a2, n4)

	filedata5, _ := regexp.Compile(",n:'[^'']*'")
	n5 := ",n:'" + "公测" + server_id2 + "区[00ff00](新)[-]'"
	a5 := filedata5.ReplaceAllString(a3, n5)

	filedata6, _ := regexp.Compile(",v:[0-9]*,w:[0-9]*,")
	n6 := ",v:0,w:1,"
	a6 := filedata6.ReplaceAllString(a5, n6)

	if filedata.dbs == nil {
		filedata.filestring = c3 + a6 + c2
		logger.Println(filedata.filestring)
	} else {
		filedata7, _ := regexp.Compile("game_db:[^]]*]")
		lenofdb := len(filedata.dbs)
		n7 := "game_db:["
		for k, v := range filedata.dbs {

			if k != lenofdb-1 {
				n7 = n7 + "{h:'" + v + "',p:27017},"
			} else if k == lenofdb-1 {
				n7 = n7 + "{h:'" + v + "',p:27017}"

			}
		}
		n8 := n7 + "]"
		a7 := filedata7.ReplaceAllString(a6, n8)

		filedata.filestring = c3 + a7 + c2
		logger.Println(filedata.filestring)
	}

}

//批量重启服务器
func reload_instance(reload_mode string) {
	if reload_mode == "gs" {
		argv1 := []string{"/nba/nba.pl", "-s", "gs", "-w", "1", "-t", "web", "-o", "restart", "-a", "yes"}
		c1 := exec.Command("perl", argv1...)
		d1, err := c1.Output()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		logger.Println(d1)
	}
	if reload_mode == "ca" {
		argv2 := []string{"/nba/nba.pl", "-s", "ca", "-w", "1", "-t", "web", "-o", "restart", "-a", "yes"}
		c2 := exec.Command("perl", argv2...)
		d2, err1 := c2.Output()
		if err1 != nil {
			logger.Println(err1.Error())
			os.Exit(1)
		}
		logger.Println(d2)
	}
	if reload_mode == "login" {
		argv3 := []string{"/nba/nba.pl", "--host", "nba_login1", "nba_login2", "nba_login3", "-t", "login", "-p", "8200", "-o", "restart"}
		c3 := exec.Command("perl", argv3...)
		d3, err3 := c3.Output()
		if err3 != nil {
			logger.Println(err3.Error())
			os.Exit(1)
		}
		logger.Println(d3)
		argv4 := []string{"/nba/nba.pl", "--host", "nba_login1", "nba_login2", "nba_login3", "-t", "login", "-p", "8100", "-o", "restart"}
		c4 := exec.Command("perl", argv4...)
		d4, err4 := c4.Output()
		if err4 != nil {
			logger.Println(err4.Error())
			os.Exit(1)
		}
		logger.Println(d4)
	}
}

//批量更新代码
func pull_code() {
	argv := []string{"/nba/full.git.zl", "node", "pull"}
	c := exec.Command("perl", argv...)
	d, err := c.Output()
	logger.Println("=== start pull code on node server ")
	if err != nil {
		logger.Println(err.Error())
		os.Exit(1)
	}
	logger.Println(d)

	argv1 := []string{"/nba/full.git.zl", "redis", "pull"}
	c1 := exec.Command("perl", argv1...)
	d1, err1 := c1.Output()
	logger.Println("=== start pull code on redis server ")
	if err1 != nil {
		logger.Println(err1.Error())
		os.Exit(1)
	}
	logger.Println(d1)

	argv2 := []string{"/nba/full.git.zl", "login", "pull", "8100"}
	c2 := exec.Command("perl", argv2...)
	d2, err2 := c2.Output()
	logger.Println("=== start pull code on login 8100 ")
	if err2 != nil {
		logger.Println(err2.Error())
		os.Exit(1)
	}
	logger.Println(d2)

	argv3 := []string{"/nba/full.git.zl", "login", "pull"}
	c3 := exec.Command("perl", argv3...)
	d3, err3 := c3.Output()
	logger.Println("=== start pull code on login 8200 ")
	if err3 != nil {
		logger.Println(err3.Error())
		os.Exit(1)
	}
	logger.Println(d3)

}

//拉静态数据
func push_s(static_version string) {
	argv := []string{"/nba/scp_staticdata.pl", "-n", static_version}
	c := exec.Command("perl", argv...)
	d, err := c.Output()
	logger.Println("=== start scp staticdata file")
	if err != nil {
		logger.Println(err.Error())
		os.Exit(1)
	}
	logger.Println(d)

}

//拉hotfix文件
func push_h(hotfix_version string) {
	argv := []string{"/nba/scp_hotfix.pl", "-n", hotfix_version}
	c := exec.Command("perl", argv...)
	d, err := c.Output()
	logger.Println("=== start scp hotfix file")
	if err != nil {
		logger.Println(err.Error())
		os.Exit(1)
	}
	logger.Println(d)
}

//切换目录
func chdir(dir string) {
	logger.Println(os.Chdir(dir))
	logger.Println(os.Getwd())
}

//git pull代码
func git_pull() {
	chdir("/data/nba/nba_game_server/")
	argv := []string{"pull"}
	c := exec.Command("git", argv...)
	d, err := c.Output()
	logger.Println("===start git pull")
	if err != nil {
		logger.Println(err.Error())
		os.Exit(1)
	}
	logger.Println(d)

}

//git 提交代码（修改config之后）
func git_push(path string) {
	chdir("/data/nba/nba_game_server/")
	argv1 := []string{"add", path}
	argv2 := []string{"commit", "-m", "update configfile"}
	argv3 := []string{"push"}
	logger.Println("===start git add")
	c1 := exec.Command("git", argv1...)
	d1, err1 := c1.Output()
	if err1 != nil {
		logger.Println(err1.Error())
		os.Exit(1)
	}
	logger.Println(d1)

	logger.Println("===start git commit")
	c2 := exec.Command("git", argv2...)
	d2, err2 := c2.Output()
	if err2 != nil {
		logger.Println(err2.Error())
		os.Exit(1)
	}
	logger.Println(d2)

	logger.Println("===start git push")
	c3 := exec.Command("git", argv3...)
	d3, err3 := c3.Output()
	if err3 != nil {
		logger.Println(err3.Error())
		os.Exit(1)
	}
	logger.Println(d3)

}

//删除文件
func removefle(path string) {
	logger.Println("==start remove ", path)
	err1 := os.Remove(path)
	if err1 != nil {
		logger.Println(err1.Error())
		os.Exit(1)
	}
	logger.Println("===remove ", path)
}

//复制配置文件
func (filedata *Configfile) copyconfig(filepath, filepath_bk string) *Configfile {
	//复制配置文件
	logger.Println("====start copy configfile")
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Println("read", filepath, "error")
		os.Exit(1)
	}
	//filedata := &Configfile{string(data)}
	filedata.filestring = string(data)
	logger.Println(filedata.filestring)

	//删除old bk file

	//removefle(filepath_bk)

	//写新的 bk file
	newfile, err3 := os.OpenFile(filepath_bk, os.O_RDWR|os.O_CREATE, 0666)
	defer newfile.Close()
	if err3 != nil {
		logger.Println(err3.Error())
		os.Exit(1)
	}
	_, err4 := newfile.WriteString(filedata.filestring)
	if err4 != nil {
		logger.Println(err4.Error())
		os.Exit(1)
	}
	logger.Println("===Copy config file done")
	return filedata
}

func (filedata *Configfile) writeconfig(filepath string) {
	//更新配置文件
	//删除老的配置文件
	//removefle(filepath)

	logger.Println("===start write new config")
	configfile, err5 := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)
	defer configfile.Close()
	if err5 != nil {
		logger.Println(err5.Error())
		os.Exit(1)
	}
	_, err6 := configfile.WriteString(filedata.filestring)
	if err6 != nil {
		logger.Println(err6.Error())
		os.Exit(1)
	}
	logger.Println("===write complete")

}

func (filedata *Configfile) changes_v() {
	//正则匹配修改s_v
	logger.Println("===start change s_v")
	filedata1, _ := regexp.Compile("version:'[0-9]*.[0-9]*")
	n1 := "version:'0." + static_version
	filedata.filestring = filedata1.ReplaceAllString(filedata.filestring, n1)
	logger.Println(filedata.filestring)

}

func (filedata *Configfile) changeh_v() {
	//正则匹配修改h_v
	logger.Println("===start change h_v")
	filedata3, _ := regexp.Compile("curHotFixVersion:[0-9]*")
	n2 := "curHotFixVersion:" + hotfix_version
	filedata.filestring = filedata3.ReplaceAllString(filedata.filestring, n2)
	logger.Println(filedata.filestring)
}

func main() {

	var err1 error
	var multi_logfile []io.Writer
	var filepath, filepath_bk string

	flag.StringVar(&mode, "mode", "", "The mode of this process")
	flag.StringVar(&hotfix_version, "h_v", "", "The version number of hotfix data")
	flag.StringVar(&static_version, "s_v", "", "The version number of statistics data")
	flag.StringVar(&target, "target", "", "Whether cn or tw")
	flag.StringVar(&server_id, "server_id", "", "The new server id")
	flag.StringVar(&ip, "ip", "", "The new server gips like 1.1.1.1,2.2.2.2,3.3.3.3")
	flag.StringVar(&db, "db", "", "The new server game_db like mongod1,mongod2,mongod3")
	flag.Parse()

	//定义filedata
	filedata := &Configfile{"", nil, nil}
	//timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	if target == "cn" {
		//filepath := "/data/nba/nba_game_server/app/config_data_cn/server_config_CN_PROD.js"
		//filepath_bk:="/tmp/server_config_CN_PROD.bk.js"
		filepath = "/Users/zhou.liyang/Desktop/server_config_CN_PROD.js"
		filepath_bk = "/Users/zhou.liyang/Desktop/server_config_CN_PROD.bk." + timestamp + ".js"
	} else if target == "tw" {

		filepath = "/data/nba/nba_game_server/app/config_data_tw/server_config_TW_PROD.js"
		filepath_bk = "/tmp/server_config_TW_PROD.bk." + timestamp + ".js"
	}
	if ip != "" {
		filedata.ips = strings.Split(ip, ",")

	}
	if db != "" {
		filedata.dbs = strings.Split(db, ",")

	}

	//建立日志文件，并初始化日志文件
	logfile, err1 = os.OpenFile("/tmp/edc_log.log", os.O_RDWR|os.O_CREATE, 0666)
	defer logfile.Close()
	if err1 != nil {
		logger.Println(err1.Error())
		os.Exit(-1)
	}
	multi_logfile = []io.Writer{
		logfile,
		os.Stdout,
	}
	logfiles := io.MultiWriter(multi_logfile...)
	logger = log.New(logfiles, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("=====job start.=====")
	if mode == "" {
		logger.Println("The mode is error pls check the input")
		logger.Println("mode :1 change statictics file version or hotfix version")
		logger.Println("mode :2 add new server line")
		os.Exit(1)
	}
	logger.Println("=== mode ", mode, " ===")

	if target == "" {

		logger.Println("The target is error pls check the input")
		os.Exit(1)
	}

	//===Mode为1
	//只更新静态数据版本
	if mode == "1" && static_version != "" && hotfix_version == "" && target != "" {

		filedata := filedata.copyconfig(filepath, filepath_bk)
		filedata.changes_v()
		filedata.writeconfig(filepath)

	}
	//只更新hotfix 版本
	if mode == "1" && static_version == "" && hotfix_version != "" && target != "" {
		logger.Println("123")
		filedata := filedata.copyconfig(filepath, filepath_bk)
		filedata.changeh_v()
		filedata.writeconfig(filepath)

	}
	//同时更新静态数据与hotfix
	if mode == "1" && static_version != "" && hotfix_version != "" && target != "" {
		filedata := filedata.copyconfig(filepath, filepath_bk)
		filedata.changeh_v()
		filedata.changes_v()
		filedata.writeconfig(filepath)

	}
	//===Mode为2
	if mode == "2" && ip == "" {
		logger.Println("Pls input the ip of new line")
	}
	if mode == "2" && ip != "" && server_id != "" && target != "" {
		filedata := filedata.copyconfig(filepath, filepath_bk)
		filedata.addnew_line(filepath, server_id)
		filedata.writeconfig(filepath)
	}

}
