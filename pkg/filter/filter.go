package filter

import (
	"net/http"
	"strings"
	"time"
	"url-collector/config"
	"url-collector/models"

	mapset "github.com/deckarep/golang-set"
)

var URLFilter *filter

type filter struct {
	uniqueSet mapset.Set
	blackList []string
}

//Init 初始化
func Init() {
	URLFilter = &filter{
		uniqueSet: mapset.NewSet(),
		blackList: config.CurrentConf.BlackList,
	}
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
