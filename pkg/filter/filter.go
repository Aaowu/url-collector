package filter

import (
	"strings"
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
func (s *filter) IsDuplicate(u *models.URL) bool {
	if s.uniqueSet.Contains(u.ID) {
		return true
	} else {
		s.uniqueSet.Add(u.ID)
		return false
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
