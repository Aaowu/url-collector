package config

//DefaultConf 默认配置
var DefaultConf AppConfig = AppConfig{
	RoutineCount: 5,
	SearchEngine: "google-image",
	BaseURL: map[string]string{
		"google":       "https://www.google.com",
		"google-image": "https://g.luciaz.me",
		"bing":         "https://cn.bing.com",
	},
}
