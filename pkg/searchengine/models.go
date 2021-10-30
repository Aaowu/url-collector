package searchengine

import (
	"context"
	"io"
	"regexp"
	"sync"
	"url-collector/pkg/alg"

	mapset "github.com/deckarep/golang-set"
)

//BaseConfig 基础配置
type BaseConfig struct {
	FetchCount   int
	Format       string
	DorkReader   io.Reader
	ResultWriter io.Writer
}

//SearchEngineConfig 搜索引擎配置
type SearchEngineConfig struct {
	BaseConfig
	baseURL    string
	userAgent  string
	nextPageRe *regexp.Regexp
}

//SearchEngine 搜索引擎
type SearchEngine struct {
	SearchEngineConfig
	atagRe          *regexp.Regexp
	progress        *alg.Progress
	dorkCh          chan string
	resultCh        chan string
	saverWg         sync.WaitGroup
	dorkWg          sync.WaitGroup
	fetcherWg       sync.WaitGroup
	ctx             context.Context
	cancel          context.CancelFunc
	FinishedDorkSet mapset.Set
}
