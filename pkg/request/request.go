package request

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
	"url-collector/config"

	"github.com/sirupsen/logrus"
)

var client *http.Client

func Init() error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if config.CurrentConf.Proxy != "" {
		proxy, err := url.Parse(config.CurrentConf.Proxy)
		if err != nil {
			logrus.Error("url.Parse failed,err:", err)
			return err
		}
		tr.Proxy = http.ProxyURL(proxy)
	}
	//跳过证书验证
	client = &http.Client{
		Transport: tr,
	}
	return nil
}

//Get 发送get请求
func Get(url string, headers map[string]string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("http.NewRequest failed,err:", err)
		return nil, err
	}
	for k, v := range headers {
		request.Header.Set(k, v)
		if v == "genIP()" {
			request.Header.Set(k, genIP())
		}
	}

	return client.Do(request)
}

//Post 发送Post请求
func Post(target string, data, headers map[string]string) (*http.Response, error) {
	params := make(url.Values)
	for k, v := range data {
		params.Set(k, v)
	}
	request, err := http.NewRequest(http.MethodPost, target, strings.NewReader(params.Encode()))
	if err != nil {
		log.Println("http.NewRequest failed,err:", err)
		return nil, err
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	return client.Do(request)
}

//genIP 产生随机的IP地址
func genIP() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
