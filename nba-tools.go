package main

import (
	"flag"
	//"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	//"os/exec"
	"regexp"
	//"strings"
)

var logfile *os.File
var logger *log.Logger
var target string
var mode, static_version, hotfix_version string

type Configfile struct {
	filestring string
}

func copyconfig(filepath string) *Configfile {
	//复制配置文件
	logger.Println("==== start copy configfile ====")
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Println("read", filepath, "error")
		os.Exit(1)
	}
	filedata := &Configfile{string(data)}
	logger.Println(filedata)

	newfile, err3 := os.OpenFile("./server_config_CN_PROD.bk.js", os.O_RDWR|os.O_CREATE, 0666)
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
	return filedata
}

func (filedata *Configfile) writeconfig(filepath string) {
	//更新配置文件
	logger.Println("===start write new config===")
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
	logger.Println("===write complete===")

}

func (filedata *Configfile) changes_v() {
	//正则匹配修改s_v
	logger.Println("===start change s_v===")
	filedata1, _ := regexp.Compile("version:'[0-9]*.[0-9]*")
	n1 := "version:'0." + static_version
	filedata.filestring = filedata1.ReplaceAllString(filedata.filestring, n1)
	logger.Println(filedata.filestring)

}

func (filedata *Configfile) changeh_v() {
	//正则匹配修改h_v
	logger.Println("===start change h_v===")
	filedata3, _ := regexp.Compile("curHotFixVersion:[0-9]*")
	n2 := "curHotFixVersion:" + hotfix_version
	filedata.filestring = filedata3.ReplaceAllString(filedata.filestring, n2)
	logger.Println(filedata.filestring)
}

func main() {

	var err1 error
	var multi_logfile []io.Writer
	var filepath string

	flag.StringVar(&mode, "mode", "", "The mode of this process")
	flag.StringVar(&hotfix_version, "h_v", "", "The version number of hotfix data")
	flag.StringVar(&static_version, "s_v", "", "The version number of statistics data")
	flag.StringVar(&target, "target", "", "Whether cn or tw")
	flag.Parse()

	if target == "cn" {
		//filepath := "/data/nba/nba_game_server/app/config_data_cn/server_config_CN_PROD.js"
		filepath = "./server_config_CN_PROD.js"
	} else if target == "tw" {
		filepath = "/data/nba/nba_game_server/app/config_data_tw/server_config_TW_PROD.js"
	}
	//建立日志文件，并初始化日志文件
	logfile, err1 = os.OpenFile("./edc_log.log", os.O_RDWR|os.O_CREATE, 0666)
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
		os.Exit(1)
	}

	if target == "" {
		logger.Println("The target is error pls check the input")
		os.Exit(1)
	}

	if mode != "" && static_version != "" && hotfix_version == "" && target != "" {

		filedata := copyconfig(filepath)
		filedata.changes_v()
		filedata.writeconfig(filepath)

	}
	if mode != "" && static_version == "" && hotfix_version != "" && target != "" {
		filedata := copyconfig(filepath)
		filedata.changeh_v()
		filedata.writeconfig(filepath)

	}
	if mode != "" && static_version != "" && hotfix_version != "" && target != "" {
		filedata := copyconfig(filepath)
		filedata.changeh_v()
		filedata.changes_v()
		filedata.writeconfig(filepath)

	}

}
