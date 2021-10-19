package searchengine

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
	"url-collector/config"
	"url-collector/models"
	"url-collector/pkg/alg"
	"url-collector/pkg/filter"

	"github.com/asmcos/requests"
	mapset "github.com/deckarep/golang-set"
)

//SearchEngine 搜索引擎
type SearchEngine struct {
	dorkReader      io.Reader
	resultWriter    io.Writer
	fetchCount      int
	baseURL         string
	RawQueryParam   string
	userAgent       string
	nextPageRe      *regexp.Regexp
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

func newSearchEngine(baseURL string, nextPageRe *regexp.Regexp, fetchCount int, dorkReader io.Reader, resultWriter io.Writer, userAgent string, rawQueryParam string) *SearchEngine {
	progress := alg.NewProgress()
	ctx, cancel := context.WithCancel(context.Background())
	atagRe := regexp.MustCompile(`<a[^>]+href="(http[^>"]+)"[^>]+>`)

	return &SearchEngine{
		baseURL:         baseURL,
		nextPageRe:      nextPageRe,
		atagRe:          atagRe,
		dorkReader:      dorkReader,
		resultWriter:    resultWriter,
		dorkCh:          make(chan string, 10240),
		resultCh:        make(chan string, 1024),
		ctx:             ctx,
		cancel:          cancel,
		progress:        progress,
		userAgent:       userAgent,
		RawQueryParam:   rawQueryParam,
		FinishedDorkSet: mapset.NewSet(),
	}
}

//NewBing Bing搜索
func NewBing(fetchCount int, dorkReader io.Reader, resultWriter io.Writer) *SearchEngine {
	nextPageRe := regexp.MustCompile(`<a[^>]+href="(/search\?q=[^>]+)"[^>]+>`)
	userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"
	rawQueryParam := ""
	return newSearchEngine(config.CurrentConf.GetBaseURL(), nextPageRe, fetchCount, dorkReader, resultWriter, userAgent, rawQueryParam)
}

//NewGoogleImage Goolge镜像搜索
func NewGoogleImage(fetchCount int, dorkReader io.Reader, resultWriter io.Writer) *SearchEngine {
	nextPageRe := regexp.MustCompile(`<a href="(/search\?q=[^>]+)" id="pnnext"[^>]+>`)
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:92.0) Gecko/20100101 Firefox/92.0"
	return newSearchEngine(config.CurrentConf.GetBaseURL(), nextPageRe, fetchCount, dorkReader, resultWriter, userAgent, "")
}

//NewGoogle Google搜索
func NewGoogle(fetchCount int, dorkReader io.Reader, resultWriter io.Writer) *SearchEngine {
	nextPageRe := regexp.MustCompile(`<a href="(/search\?q=[^>]+)" id="pnnext"[^>]+>`)
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:92.0) Gecko/20100101 Firefox/92.0"
	return newSearchEngine(config.CurrentConf.GetBaseURL(), nextPageRe, fetchCount, dorkReader, resultWriter, userAgent, "")
}

func (s *SearchEngine) save() {
	go func() {
		s.saverWg.Add(1)
		for result := range s.resultCh {
			time.Sleep(time.Millisecond * 100)
			fmt.Fprintln(s.resultWriter, result)
		}
		s.saverWg.Done()
	}()
}

func (s *SearchEngine) fetch() {
	for i := 0; i < config.CurrentConf.RoutineCount; i++ {
		go func() {
			s.fetcherWg.Add(1)
		LOOP:
			for {
				select {
				case <-s.ctx.Done():
					break LOOP
				case dork := <-s.dorkCh:
					//1.发送请求
					headers := requests.Header{
						"User-Agent": s.userAgent,
					}
					resp, err := requests.Get(dork, headers)
					if err != nil {
						log.Printf("requests.Get(%s) failed,err:%v", dork, err)
						continue
					}
					//2.解析响应 寻找指定文件的URL
					if strings.Contains(resp.Text(), "需要验证您是否来自浙江大学") {
						if err := s.postAnswer(); err != nil {
							log.Println("s.postAnswer failed,err:", err)
							return
						}
						s.dorkCh <- dork
						continue
					}
					matches := s.atagRe.FindAllStringSubmatch(resp.Text(), -1)
					for _, match := range matches {
						link := match[1]
						URL, err := models.NewURL(link)
						if err != nil {
							continue
						}
						//1.过滤重复的url
						if filter.URLFilter.IsDuplicate(URL) {
							continue
						}
						//2.过去黑名单中的url 例如gov
						if filter.URLFilter.IsInBlackList(link) {
							continue
						}
						s.resultCh <- link
					}
					//3.寻找“下一页URL”
					u, err := url.Parse(dork)
					if err != nil {
						log.Println("url.Parse failed,err:", err)
						return
					}
					keyword := u.Query().Get("q")
					nextPageURLs := make([]string, 0)
					matches = s.nextPageRe.FindAllStringSubmatch(resp.Text(), -1)
					for _, match := range matches {
						nextPageURL := s.baseURL + match[1]
						temp, err := models.NewURL(nextPageURL)
						if err != nil {
							log.Println("models.NewURL failed,err:", err)
							continue
						}
						if temp.URL.Query().Get("q") != keyword && temp.URL.Query().Get("q") != url.QueryEscape(keyword) {
							continue
						}
						nextPageURLs = append(nextPageURLs, nextPageURL)
					}

					if len(nextPageURLs) > 0 {
						for i := range nextPageURLs {
							s.dorkCh <- nextPageURLs[i]
						}
					} else {
						if s.FinishedDorkSet.Contains(keyword) {
							continue
						}
						s.FinishedDorkSet.Add(keyword)
						s.dorkWg.Done() //针对该dork的搜索任务完成
						s.progress.AddFinished()
					}
				}
			}
			s.fetcherWg.Done()
		}()
	}
}

//Search 开始搜索
func (s *SearchEngine) Search() {
	//定时显示进度
	s.progress.Show(s.ctx)
	//保存结果  (消费者:快)
	s.save()
	//发送请求 （生产者:慢）
	s.fetch()
	//从reader中读取dork
	scanner := bufio.NewScanner(s.dorkReader)
	for scanner.Scan() {
		keyword := strings.TrimSpace(scanner.Text())
		request := fmt.Sprintf("%s/search?q=%s&%s", s.baseURL, url.QueryEscape(keyword), s.RawQueryParam)
		s.dorkCh <- request
		s.dorkWg.Add(1)
		s.progress.AddTotal()
	}
	//等待各部门结束工作
	s.wait()
}

func (s *SearchEngine) wait() {
	//因为dork是有限的，所以等待所有dork搜索完成
	s.dorkWg.Wait()
	//通知fetcher和progress结束工作
	s.cancel()
	s.fetcherWg.Wait()
	//关闭resultCh
	close(s.resultCh)
	//等待saver结束工作
	s.saverWg.Wait()
	fmt.Println("\n搜索完成")
}

//提交答案
func (s *SearchEngine) postAnswer() error {
	headers := requests.Header{
		"User-Agent": s.userAgent,
	}
	data := requests.Datas{
		"0":      "心灵之约",
		"1":      "水朝夕",
		"2":      "csxy@123",
		"origin": "aHR0cHM6Ly9nLmx1Y2lhei5tZS8=",
	}
	resp, err := requests.Post("https://g.luciaz.me/ip_ban_verify_page?origin=aHR0cHM6Ly9nLmx1Y2lhei5tZS8=", data, headers)
	if err != nil {
		return err
	}
	if strings.Contains(resp.Text(), "Wrong answer") {
		return errors.New("Wrong answer")
	}
	return nil
}
