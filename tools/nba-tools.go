package main

import (
	"bufio"
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

//###################################################
//###############定义全局变量##########################
//###################################################

var logfile *os.File
var logger *log.Logger
var target string
var mode, static_version, hotfix_version, server_id, ip, db string

//###################################################
//###############定义struct##########################
//###################################################
type Configfile struct {
	filestring string
	ips        []string
	dbs        []string
}

//###################################################
//###############重启gmtool##########################
//###################################################

func reload_gm(interface{}) {
	argv2 := []string{"/nba/gm_reload.pl"}
	c2 := exec.Command("perl", argv2...)
	d2, err1 := c2.Output()
	if err1 != nil {
		logger.Println(err1.Error())
		os.Exit(1)
	}
	logger.Println(d2)

}

func reload_cdn() {

	if target == "cn" {

	}
	if target == "cn2015" {
		url_i := "http://download-nba2015-c.mobage.cn:8400/1/ios/StaticData/StaticData_" + static_version + "_0.unity3d"
		url_a := "http://download-nba2015-c.mobage.cn:8400/1/android/StaticData/StaticData_" + static_version + "_0.unity3d"

		argv2 := []string{"/nba/cdn.py", url_i}
		argv3 := []string{"/nba/cdn.py", url_a}
		c2 := exec.Command("python", argv2...)
		c3 := exec.Command("python", argv3...)
		logger.Println("Start flush cdn ios for  ", static_version)
		printlog(c2)
		logger.Println("flush cdn ios Done")
		logger.Println("Start flush cdn android for  ", static_version)
		printlog(c3)
		logger.Println("flush cdn android Done")

	}
	if target == "tw" {

	}

	argv2 := []string{"/nba/cdn.py", "-s", "ca", "-w", "1", "-t", "web", "-o", "restart", "-a", "yes"}
	c2 := exec.Command("python", argv2...)

	logger.Println("Start restart calculation ")
	printlog(c2)
	logger.Println("restart calculation server Done")
}

//###################################################
//###############打印os/exec调用系统命令的返回##########
//###################################################

func printlog(c2 *exec.Cmd) {
	stdout, err := c2.StdoutPipe()
	if err != nil {
		logger.Println(err.Error())
		os.Exit(1)
	}
	if err := c2.Start(); err != nil {
		logger.Println("Command  error:", err.Error())
		os.Exit(1)
	}
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		logger.Println(in.Text())
	}
	if err := in.Err(); err != nil {
		logger.Println("Err:", err.Error())
		os.Exit(1)
	}
}

//###################################################
//###############为配置文件新增行######################
//###################################################

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

//###################################################
//###############批量重启服务##########################
//###################################################

//批量重启服务器
func reload_instance(reload_mode string) {
	if reload_mode == "gs" {
		argv1 := []string{"/nba/nba.pl", "-s", "gs", "-w", "1", "-t", "web", "-o", "restart", "-a", "yes"}
		c1 := exec.Command("perl", argv1...)
		logger.Println("Start restart node server")

		/*stdout, err := c1.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c1.Start(); err != nil {
			logger.Println("Command  error:", err.Error())
			os.Exit(1)
		}
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			logger.Println(in.Text())
		}
		if err := in.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
		*/
		printlog(c1)
		logger.Println("restart node server Done")

	}
	if reload_mode == "ca" {
		argv2 := []string{"/nba/nba.pl", "-s", "ca", "-w", "1", "-t", "web", "-o", "restart", "-a", "yes"}
		c2 := exec.Command("perl", argv2...)
		logger.Println("Start restart calculation ")

		/*stdout, err := c2.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c2.Start(); err != nil {
			logger.Println("Command  error:", err.Error())
			os.Exit(1)
		}
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			logger.Println(in.Text())
		}
		if err := in.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
		*/
		printlog(c2)
		logger.Println("restart calculation server Done")

	}
	if reload_mode == "login" {
		argv3 := []string{"/nba/nba.pl", "--host", "nba_login1", "nba_login2", "nba_login3", "-t", "login", "-p", "8200", "-o", "restart"}
		c3 := exec.Command("perl", argv3...)
		logger.Println("Start restart CN 8200 login")
		/*
			stdout, err := c3.StdoutPipe()
			if err != nil {
				logger.Println(err.Error())
				os.Exit(1)
			}
			if err := c3.Start(); err != nil {
				logger.Println("Command  error:", err.Error())
				os.Exit(1)
			}
			in := bufio.NewScanner(stdout)
			for in.Scan() {
				logger.Println(in.Text())
			}
			if err := in.Err(); err != nil {
				logger.Println("Err:", err.Error())
				os.Exit(1)
			}
		*/
		printlog(c3)
		logger.Println("restart login CN 8200 Done")

		argv4 := []string{"/nba/nba.pl", "--host", "nba_login1", "nba_login2", "nba_login3", "-t", "login", "-p", "8100", "-o", "restart"}
		c4 := exec.Command("perl", argv4...)
		logger.Println("Start restart CN 8100 login")
		/*
			stdout1, err := c4.StdoutPipe()
			if err != nil {
				logger.Println(err.Error())
				os.Exit(1)
			}
			if err := c4.Start(); err != nil {
				logger.Println("Command  error:", err.Error())
				os.Exit(1)
			}
			in1 := bufio.NewScanner(stdout1)
			for in1.Scan() {
				logger.Println(in1.Text())
			}
			if err := in1.Err(); err != nil {
				logger.Println("Err:", err.Error())
				os.Exit(1)
			}
		*/
		printlog(c4)
		logger.Println("restart login CN 8100 Done")

	}

	if reload_mode == "login2015" {
		argv3 := []string{"/nba/nba.pl", "--host", "nba2015_login1", "nba2015_login2", "-t", "login", "-p", "8200", "-o", "restart"}
		c3 := exec.Command("perl", argv3...)
		logger.Println("Start restart 2015 8200 login")
		/*
			stdout, err := c3.StdoutPipe()
			if err != nil {
				logger.Println(err.Error())
				os.Exit(1)
			}
			if err := c3.Start(); err != nil {
				logger.Println("Command  error:", err.Error())
				os.Exit(1)
			}
			in := bufio.NewScanner(stdout)
			for in.Scan() {
				logger.Println(in.Text())
			}
			if err := in.Err(); err != nil {
				logger.Println("Err:", err.Error())
				os.Exit(1)
			}
		*/
		printlog(c3)
		logger.Println("restart login 2015 8200 Done")

		argv4 := []string{"/nba/nba.pl", "--host", "nba2015_login1", "nba2015_login2", "-t", "login", "-p", "8100", "-o", "restart"}
		c4 := exec.Command("perl", argv4...)
		logger.Println("Start restart 2015 8100 login")

		/*
			stdout1, err := c4.StdoutPipe()
			if err != nil {
				logger.Println(err.Error())
				os.Exit(1)
			}
			if err := c4.Start(); err != nil {
				logger.Println("Command  error:", err.Error())
				os.Exit(1)
			}
			in1 := bufio.NewScanner(stdout1)
			for in1.Scan() {
				logger.Println(in1.Text())
			}
			if err := in1.Err(); err != nil {
				logger.Println("Err:", err.Error())
				os.Exit(1)
			}
		*/
		printlog(c4)
		logger.Println("restart login 2015 8100 Done")

	}
	if reload_mode == "logintw" {
		argv3 := []string{"/nba/nba.pl", "--host", "twnba_login1", "twnba_login2", "-t", "login", "-p", "8200", "-o", "restart"}
		c3 := exec.Command("perl", argv3...)
		logger.Println("Start restart TW 8200 login")
		/*
			stdout, err := c3.StdoutPipe()
			if err != nil {
				logger.Println(err.Error())
				os.Exit(1)
			}
			if err := c3.Start(); err != nil {
				logger.Println("Command  error:", err.Error())
				os.Exit(1)
			}
			in := bufio.NewScanner(stdout)
			for in.Scan() {
				logger.Println(in.Text())
			}
			if err := in.Err(); err != nil {
				logger.Println("Err:", err.Error())
				os.Exit(1)
			}
		*/
		printlog(c3)
		logger.Println("restart login TW 8200 Done")

		argv4 := []string{"/nba/nba.pl", "--host", "twnba_login1", "twnba_login2", "-t", "login", "-p", "8100", "-o", "restart"}
		c4 := exec.Command("perl", argv4...)
		logger.Println("Start restart TW 8100 login")
		/*
			stdout1, err := c4.StdoutPipe()
			if err != nil {
				logger.Println(err.Error())
				os.Exit(1)
			}
			if err := c4.Start(); err != nil {
				logger.Println("Command  error:", err.Error())
				os.Exit(1)
			}
			in1 := bufio.NewScanner(stdout1)
			for in1.Scan() {
				logger.Println(in1.Text())
			}
			if err := in1.Err(); err != nil {
				logger.Println("Err:", err.Error())
				os.Exit(1)
			}
		*/
		printlog(c4)
		logger.Println("restart login TW 8100 Done")

	}
}

//###################################################
//###############批量更新服务器代码#####################
//###################################################

func pull_code() {

	//###################################################
	//###############开始为node服务器更新代码###############
	//###################################################
	argv := []string{"/nba/full.git.zl", "node", "pull"}
	c := exec.Command("perl", argv...)

	//只输出到屏幕不记录log的简单方法
	/*
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if err != nil {
			logger.Println(err.Error())
			return
		}
	*/
	//输出到logger
	logger.Println("Start full code on node servers")
	/*
		stdout, err := c.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c.Start(); err != nil {
			logger.Println("Command Start error:", err.Error())
			os.Exit(1)
		}
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			logger.Println(in.Text())
		}
		if err := in.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
	*/
	printlog(c)
	logger.Println("pull code on node servers Done")

	//###################################################
	//###############开始为redis服务器更新代码##############
	//###################################################
	logger.Println("Start full code on redis servers")
	argv1 := []string{"/nba/full.git.zl", "redis", "pull"}
	c1 := exec.Command("perl", argv1...)
	/*
		stdout1, err := c1.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c1.Start(); err != nil {
			logger.Println("Command Start error:", err.Error())
			os.Exit(1)
		}
		in1 := bufio.NewScanner(stdout1)
		for in1.Scan() {
			logger.Println(in1.Text())
		}
		if err := in1.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
	*/
	printlog(c1)
	logger.Println("pull code on redis servers Done")

	//###################################################
	//###############开始为login服务器更新代码##############
	//###################################################
	logger.Println("Start full code on login servers")
	argv2 := []string{"/nba/full.git.zl", "login", "pull", "8100"}
	c2 := exec.Command("perl", argv2...)
	/*
		stdout2, err := c2.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c2.Start(); err != nil {
			logger.Println("Command Start error:", err.Error())
			os.Exit(1)
		}
		in2 := bufio.NewScanner(stdout2)
		for in2.Scan() {
			logger.Println(in2.Text())
		}
		if err := in2.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
	*/
	printlog(c2)
	logger.Println("pull code on login 8100 Done")

	argv3 := []string{"/nba/full.git.zl", "login", "pull"}
	c3 := exec.Command("perl", argv3...)

	/*
		stdout3, err := c3.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c3.Start(); err != nil {
			logger.Println("Command Start error:", err.Error())
			os.Exit(1)
		}
		in3 := bufio.NewScanner(stdout3)
		for in3.Scan() {
			logger.Println(in3.Text())
		}
		if err := in3.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
	*/
	printlog(c3)
	logger.Println("pull code on login 8200 Done")

}

//###################################################
//###############拉去静态数据#########################
//###################################################

//拉静态数据
func pull_s(static_version string) {
	argv := []string{"/nba/scp_staticdata.pl", "-n", static_version}
	c := exec.Command("perl", argv...)
	logger.Println("Start scp statiscdata")
	/*
		stdout, err := c.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c.Start(); err != nil {
			logger.Println("Command  error:", err.Error())
			os.Exit(1)
		}
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			logger.Println(in.Text())
		}
		if err := in.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
	*/
	printlog(c)
	logger.Println("pull static file Done")
}

//###################################################
//###############拉去静态数据file######################
//###################################################
//拉hotfix文件
func pull_h(hotfix_version string) {
	argv := []string{"/nba/scp_hotfix.pl", "-n", hotfix_version}
	c := exec.Command("perl", argv...)
	logger.Println("Start scp hotfix file")
	/*
		stdout, err := c.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c.Start(); err != nil {
			logger.Println("Command  error:", err.Error())
			os.Exit(1)
		}
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			logger.Println(in.Text())
		}
		if err := in.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
	*/
	printlog(c)
	logger.Println("pull hotfix file Done")
}

//###################################################
//###############切换目录#############################
//###################################################

//切换目录
func chdir(dir string) {
	logger.Println(os.Chdir(dir))
	logger.Println(os.Getwd())
}

//###################################################
//###############本地git pull 代码####################
//###################################################

//git pull代码
func git_pull(filepath string) {
	chdir(filepath)
	argv := []string{"pull"}
	c := exec.Command("git", argv...)
	logger.Println("Start pull code localhost")
	/*
		stdout, err := c.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c.Start(); err != nil {
			logger.Println("Command  error:", err.Error())
			os.Exit(1)
		}
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			logger.Println(in.Text())
		}
		if err := in.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
	*/
	printlog(c)
	logger.Println("local git pull Done")
}

//###################################################
//###############本地git push 代码####################
//###################################################

//git 提交代码（修改config之后）
func git_push(path string, filepath1 string) {
	chdir(filepath1)
	argv1 := []string{"add", path}
	argv2 := []string{"commit", "-m", "update configfile"}
	argv3 := []string{"push"}
	logger.Println("start git add")
	c1 := exec.Command("git", argv1...)
	/*
		stdout, err := c1.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c1.Start(); err != nil {
			logger.Println("Command  error:", err.Error())
			os.Exit(1)
		}
		in := bufio.NewScanner(stdout)
		for in.Scan() {
			logger.Println(in.Text())
		}
		if err := in.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
	*/
	printlog(c1)
	logger.Println("Git add done")

	logger.Println("start git commit")
	c2 := exec.Command("git", argv2...)

	/*
		stdout2, err := c2.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c2.Start(); err != nil {
			logger.Println("Command  error:", err.Error())
			os.Exit(1)
		}
		in2 := bufio.NewScanner(stdout2)
		for in2.Scan() {
			logger.Println(in2.Text())
		}
		if err := in2.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
	*/
	printlog(c2)
	logger.Println("Git commit Done")

	logger.Println("start git push")
	c3 := exec.Command("git", argv3...)
	/*
		stdout3, err := c3.StdoutPipe()
		if err != nil {
			logger.Println(err.Error())
			os.Exit(1)
		}
		if err := c3.Start(); err != nil {
			logger.Println("Command  error:", err.Error())
			os.Exit(1)
		}
		in3 := bufio.NewScanner(stdout3)
		for in3.Scan() {
			logger.Println(in3.Text())
		}
		if err := in3.Err(); err != nil {
			logger.Println("Err:", err.Error())
			os.Exit(1)
		}
	*/
	printlog(c3)
	logger.Println("git push Done")

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

//###################################################
//###############备份配置文件##########################
//###################################################

//复制配置文件
func (filedata *Configfile) copyconfig(filepath, filepath_bk string) *Configfile {
	//复制配置文件
	logger.Println("====start copy configfile")
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Println("read", filepath, "error")
		os.Exit(1)
	}

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

//###################################################
//###############更新配置文件##########################
//###################################################
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

//###################################################
//###############修改静态数据版本号####################
//###################################################

func (filedata *Configfile) changes_v() {
	//正则匹配修改s_v
	logger.Println("===start change s_v")
	filedata1, _ := regexp.Compile("version:'[0-9]*.[0-9]*")
	n1 := "version:'0." + static_version
	filedata.filestring = filedata1.ReplaceAllString(filedata.filestring, n1)
	logger.Println(filedata.filestring)

}

//###################################################
//###############修改热修版本号#######################
//###################################################

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
	var filepath, filepath_bk, codepath string

	flag.StringVar(&mode, "mode", "", "The mode of this process")
	flag.StringVar(&hotfix_version, "h_v", "", "The version number of hotfix data")
	flag.StringVar(&static_version, "s_v", "", "The version number of statistics data")
	flag.StringVar(&target, "target", "", "Whether cn or cn2015 or tw")
	flag.StringVar(&server_id, "server_id", "", "The new server id")
	flag.StringVar(&ip, "ip", "", "The new server gips like 1.1.1.1,2.2.2.2,3.3.3.3")
	flag.StringVar(&db, "db", "", "The new server game_db like mongod1,mongod2,mongod3,if do not input then copy the last line of file")
	flag.Parse()

	//定义filedata
	filedata := &Configfile{"", nil, nil}

	timestamp := time.Now().Format("20060102150405")
	if target == "cn" || target == "cn2015" {
		filepath = "/data/zhou/nba_game_server/app/config_data_cn/server_config_CN_PROD.js"
		filepath_bk = "/tmp/server_config_CN_PROD.bk." + timestamp + ".js"
		//filepath = "/Users/zhou.liyang/Desktop/server_config_CN_PROD.js"
		//filepath_bk = "/Users/zhou.liyang/Desktop/server_config_CN_PROD.bk." + timestamp + ".js"
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
	logpath := "/tmp/edc_log.log." + timestamp
	logfile, err1 = os.OpenFile(logpath, os.O_RDWR|os.O_CREATE, 0666)
	defer logfile.Close()
	if err1 != nil {
		logger.Println(err1.Error())
		os.Exit(1)
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
	if target == "cn" || target == "cn2015" {
		codepath = "/data/zhou/nba_game_server"

	}
	if target == "tw" {
		codepath = "/data/nba/nba_game_server"
	}

	//===Mode为1
	//只更新静态数据版本
	if mode == "1" && static_version != "" && hotfix_version == "" && target != "" {
		git_pull(codepath)
		filedata := filedata.copyconfig(filepath, filepath_bk)
		filedata.changes_v()
		filedata.writeconfig(filepath)
		git_push(filepath, codepath)
		pull_s(static_version)
		pull_code()
		if target == "cn2015" {
			reload_instance("login2015")
		}
		if target == "tw" {
			reload_instance("logintw")
		}
		if target == "cn" {
			reload_instance("login")
		}
		reload_instance("gs")

	}
	//只更新hotfix 版本
	if mode == "1" && static_version == "" && hotfix_version != "" && target != "" {
		git_pull(codepath)
		filedata := filedata.copyconfig(filepath, filepath_bk)
		filedata.changeh_v()
		filedata.writeconfig(filepath)
		git_push(filepath, codepath)
		pull_h(hotfix_version)
		pull_code()
		if target == "cn2015" {
			reload_instance("login2015")
		}
		if target == "tw" {
			reload_instance("logintw")
		}
		if target == "cn" {
			reload_instance("login")
		}
		reload_instance("gs")

	}
	//同时更新静态数据与hotfix
	if mode == "1" && static_version != "" && hotfix_version != "" && target != "" {
		git_pull(codepath)
		filedata := filedata.copyconfig(filepath, filepath_bk)
		filedata.changeh_v()
		filedata.changes_v()
		filedata.writeconfig(filepath)
		git_push(filepath, codepath)
		pull_h(hotfix_version)
		pull_code()
		if target == "cn2015" {
			reload_instance("login2015")
		}
		if target == "tw" {
			reload_instance("logintw")
		}
		if target == "cn" {
			reload_instance("login")
		}

		reload_instance("gs")

	}
	//===Mode为2
	if mode == "2" && ip == "" {
		logger.Println("Pls input the ip of new line")
	}
	if mode == "2" && ip != "" && server_id != "" && target != "" {
		git_pull(codepath)
		filedata := filedata.copyconfig(filepath, filepath_bk)
		filedata.addnew_line(filepath, server_id)
		filedata.writeconfig(filepath)
	}

}
