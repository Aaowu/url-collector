package filter

import (
	"bufio"
	"net/http"
	"os"
	"strings"
	"time"
	"url-collector/config"
	"url-collector/models"

	mapset "github.com/deckarep/golang-set"
	"github.com/sirupsen/logrus"
)

var URLFilter *filter

type filter struct {
	uniqueSet mapset.Set
	blackList []string
}

//Init 初始化
func Init() error {
	URLFilter = &filter{
		uniqueSet: mapset.NewSet(),
		blackList: config.CurrentConf.BlackList,
	}
	//从OutputFile中读取已经采集的url
	if config.CurrentConf.OutputFilePath != "" {
		file, err := os.OpenFile(config.CurrentConf.OutputFilePath, os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			logrus.Error("os.Open failed,err:", err)
			return err
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			URL, err := models.NewURL(strings.TrimSpace(scanner.Text()))
			if err != nil {
				logrus.Error("models.NewURL failed,err:", err)
				continue
			}
			if !URLFilter.uniqueSet.Contains(URL.ID) {
				URLFilter.uniqueSet.Add(URL.ID)
			}
		}
	}
	return nil
}

//IsDuplicate 是否重复
func (s *filter) IsDuplicate(url string) (bool, error) {
	URL, err := models.NewURL(url)
	if err != nil {
		return true, err
	}
	if s.uniqueSet.Contains(URL.ID) {
		return true, nil
	} else {
		s.uniqueSet.Add(URL.ID)
		return false, nil
	}

}

//IsInBlackList 是否在黑名单中
func (s *filter) IsInBlackList(link string) bool {
	for i := range s.blackList {
		if strings.Contains(link, s.blackList[i]) {
			return true
		}
	}
	return false
}

func (s *filter) CheckRedirect(link string) (url string, err error) {
	list := []string{
		"www.baidu.com/link?url=",
		"www.baidu.com/baidu.php?url=",
	}
	for i := range list {
		if strings.Contains(link, list[i]) {
			c := http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
				Timeout: 30 * time.Second,
			}
			req, err := http.NewRequest(http.MethodGet, link, nil)
			if err != nil {
				return "", err
			}
			resp, err := c.Do(req)
			if err != nil {
				return "", err
			}
			return resp.Header.Get("location"), nil
		}
	}
	return link, nil
}
