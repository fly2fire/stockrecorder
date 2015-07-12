package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nzai/stockrecorder/config"
	"github.com/nzai/stockrecorder/task"
)

const (
	configFileName           = "config.ini"
	configLogSection         = "path"
	configLogKey             = "logpath"
	configLogDefaultFileName = "main.log"
)

func main() {
	//	当前目录
	root := filepath.Dir(os.Args[0])
	filename := filepath.Join(root, configFileName)

	//	读取配置文件
	err := config.SetConfigFile(filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	//	打开日志文件
	file, err := openLogFile()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	//	设置日志输出文件
	log.SetOutput(file)

	//	启动任务
	err = task.StartTasks()
	if err != nil {
		log.Fatal(err)
		return
	}

	//	阻塞，一直运行
	channel := make(chan int)
	<-channel
}

//	打开日志文件
func openLogFile() (*os.File, error) {
	//	日志文件路径
	logPath := config.GetString(configLogSection, configLogKey, configLogDefaultFileName)
	logDir := filepath.Dir(logPath)
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(logDir, 0x644)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0x644)
}
