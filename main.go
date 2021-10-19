package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"url-collector/config"
	"url-collector/pkg/filter"
	"url-collector/pkg/searchengine"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func main() {
	//初始化配置
	config.Init()
	//初始化过滤器
	filter.Init()
	author := cli.Author{
		Name:  "无在无不在",
		Email: "2227627947@qq.com",
	}
	app := &cli.App{
		Name:      "URL-Collector",
		Usage:     "Collect URLs based on dork",
		UsageText: "url-collector",
		Version:   "v0.1",
		Authors:   []*cli.Author{&author},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "input from a file",
				Destination: &config.CurrentConf.InputFilePath,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "specify the output file",
				Destination: &config.CurrentConf.OutputFilePath,
			},
			&cli.StringFlag{
				Name:        "engine",
				Aliases:     []string{"e"},
				Usage:       "specify the search engine(google,bing)",
				Value:       config.DefaultConf.SearchEngine,
				Destination: &config.CurrentConf.SearchEngine,
			},
			&cli.IntFlag{
				Name:        "routine-count",
				Aliases:     []string{"c"},
				Usage:       "specify the count of goroutine",
				Value:       config.DefaultConf.RoutineCount,
				Destination: &config.CurrentConf.RoutineCount,
			},
			&cli.StringFlag{
				Name:        "keyword",
				Aliases:     []string{"k"},
				Usage:       "specify the keyword",
				Destination: &config.CurrentConf.Keyword,
			},
		},
		Action: run,
	}
	//启动app
	if err := app.Run(os.Args); err != nil {
		logrus.Error(err)
	}
}

func run(c *cli.Context) (err error) {
	//1.抽象出一个Reader
	reader, err := config.CurrentConf.GetReader()
	if err != nil {
		cli.ShowAppHelp(c)
		return err
	}
	//2.抽象出一个Writer
	writer, err := config.CurrentConf.GetWriter()
	if err != nil {
		cli.ShowAppHelp(c)
		return err
	}
	//3.创建搜索引擎对象
	var engine *searchengine.SearchEngine
	switch config.CurrentConf.SearchEngine {
	case "google-image":
		engine = searchengine.NewGoogleImage(config.CurrentConf.RoutineCount, reader, writer)
	case "google":
		engine = searchengine.NewGoogle(config.CurrentConf.RoutineCount, reader, writer)
	case "bing":
		engine = searchengine.NewBing(config.CurrentConf.RoutineCount, reader, writer)
	default:
		return errors.New("please specify a search engine,such as google-image,google,bing")
	}
	//4.开始采集
	showConfig()
	engine.Search()
	return nil
}

func showConfig() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"SearchEngine", "BaseURL", "RoutineCount"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)
	data := [][]string{
		{config.CurrentConf.SearchEngine, config.CurrentConf.GetBaseURL(), fmt.Sprintf("%d", config.CurrentConf.RoutineCount)},
	}
	table.AppendBulk(data)
	table.SetCaption(true, "Current Config")
	table.Render()
	fmt.Println("[*] black list:", strings.Join(config.CurrentConf.BlackList, ","))
	fmt.Println("[*] collecting...")
}
