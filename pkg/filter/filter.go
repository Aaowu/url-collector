package filter

import (
	"strings"
	"url-collector/models"

	mapset "github.com/deckarep/golang-set"
)

var URLFilter *filter = newFilter()

type filter struct {
	uniqueSet mapset.Set
	blackList []string
}

func newFilter() *filter {
	return &filter{
		uniqueSet: mapset.NewSet(),
		blackList: []string{"gov",
			"g3.luciaz.me",
			"www.youtube.com",
			"gitee.com",
			"github.com",
			"stackoverflow.com",
			"developer.aliyun.com",
			"cloud.tencent.com",
			"www.zhihu.com/question",
			"blog.51cto.com",
			"zhidao.baidu.com",
			"www.cnblogs.com",
			"coding.m.imooc.com",
			"weibo.cn",
			"www.taobao.com",
			"www.google.com",
			"go.microsoft.com",
			"facebook.com",
			"blog.csdn.net",
		},
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
