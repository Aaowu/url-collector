package models

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type URL struct {
	ID          string
	QueryParams []string
	URL         *url.URL
}

func NewURL(link string) (urlobj *URL, err error) {
	urlobj = new(URL)
	urlobj.QueryParams = make([]string, 0)
	urlobj.URL, err = url.Parse(link)
	if err != nil {
		return nil, err
	}
	for k := range urlobj.URL.Query() {
		urlobj.QueryParams = append(urlobj.QueryParams, k)
	}
	//排序，确保顺序一致，防止出现  a&b&c != b&a&c 的情况
	sort.Strings(urlobj.QueryParams)
	text := urlobj.URL.Host + urlobj.URL.Path + "?" + strings.Join(urlobj.QueryParams, "&")
	urlobj.ID = fmt.Sprintf("%x", md5.Sum([]byte(text)))
	return urlobj, nil
}
