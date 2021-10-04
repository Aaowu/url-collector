package config

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

//CurrentConf 当前配置
var CurrentConf = new(AppConfig)

//AppConfig App配置
type AppConfig struct {
	InputFilePath  string
	OutputFilePath string
	SearchEngine   string
	Keyword        string
	RoutineCount   int
}

//GetReader 将输入的文件或者关键字抽象为一个Reader
func (a *AppConfig) GetReader() (io.Reader, error) {
	if 0 != len(a.InputFilePath) {
		reader, err := os.Open(a.InputFilePath)
		if err != nil {
			log.Println("os.Open failed,err:", err)
			return nil, err
		}
		return reader, nil
	}
	if 0 != len(a.Keyword) {
		reader := strings.NewReader(a.Keyword)
		return reader, nil
	}
	return nil, errors.New("specify -f or -k please")
}

//GetWriter 将输出抽象为一个Writer
func (a *AppConfig) GetWriter() (io.Writer, error) {
	if 0 != len(a.OutputFilePath) {
		dstFile, err := os.OpenFile(a.OutputFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			log.Printf("os.OpenFile() failed,err:%v\n", err)
			return nil, err
		}
		return dstFile, nil
	}
	return os.Stdout, nil
}
