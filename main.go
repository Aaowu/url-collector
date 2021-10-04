package main

import (
	"errors"
	"os"
	"url-collector/config"
	"url-collector/pkg/searchengine"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func main() {

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

	err := app.Run(os.Args)
	if err != nil {
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
	case "google":
		engine = searchengine.NewGoogle(config.CurrentConf.RoutineCount, reader, writer)
		break
	case "bing":
		engine = searchengine.NewBing(config.CurrentConf.RoutineCount, reader, writer)
		break
	default:
		return errors.New("pls specify a search engine,such as google,bing")
	}
	//4.开始采集
	engine.Search()
	return nil
}
