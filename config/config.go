package config

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

//CurrentConf 当前配置
var CurrentConf AppConfig

//AppConfig App配置
type AppConfig struct {
	RoutineCount   int               `mapstructure:"routine_count"`
	BlackList      []string          `mapstructure:"black_list"`
	BaseURL        map[string]string `mapstructure:"base_url"`
	InputFilePath  string            `mapstructure:"input_file_path"`
	OutputFilePath string            `mapstructure:"ouput_file_path"`
	SearchEngine   string            `mapstructure:"search_engine"`
	Keyword        string            `mapstructure:"keyword"`
}

//Init 初始化配置
func Init(filePath string) error {
	CurrentConf.BaseURL = DefaultConf.BaseURL
	CurrentConf.BlackList = DefaultConf.BlackList
	if len(filePath) == 0 {
		return nil
	}
	//指定配置文件
	viper.SetConfigFile(filePath)
	//读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("viper.ReadInConfig() failed,err:%v\n", err)
		return err
	}
	//反序列化配置信息
	if err := viper.Unmarshal(&CurrentConf); err != nil {
		fmt.Printf("viper.Unmarshal(&CurrentConf) failed,err:%v\n", err)
		return err
	}
	return nil
}

//GetBaseURL 获取搜索引擎对应的baseURL
func (a *AppConfig) GetBaseURL() string {
	return a.BaseURL[a.SearchEngine]
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
