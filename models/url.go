package models

import (
	"crypto/md5"
	"fmt"
	"net/url"

	mapset "github.com/deckarep/golang-set"
)

type URL struct {
	ID  string
	URL *url.URL
}

func NewURL(link string) (urlobj *URL, err error) {
	urlobj = new(URL)
	urlobj.URL, err = url.Parse(link)
	if err != nil {
		return nil, err
	}
	getParamKeys := mapset.NewSet()
	for k := range urlobj.URL.Query() {
		getParamKeys.Add(k)
	}
	text := urlobj.URL.Host + urlobj.URL.Path + "?" + getParamKeys.String()
	urlobj.ID = fmt.Sprintf("%x", md5.Sum([]byte(text)))
	return urlobj, nil
}
