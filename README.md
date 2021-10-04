# URL-Collector url采集器

## Usage
```shell
NAME:
   URL-Collector - Collect URLs based on dork

USAGE:
   url-collector

VERSION:
   v0.1

AUTHOR:
   无在无不在 <2227627947@qq.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value           input from a file
   --output value, -o value         specify the output file
   --engine value, -e value         specify the search engine (default: "google")
   --routine-count value, -c value  specify the count of goroutine (default: 5)
   --keyword value, -k value        specify the keyword
   --help, -h                       show help (default: false)
   --version, -v                    print the version (default: false)
```

## Demo 

```shell
url-collector -k ".php?id=" 
```

![avatar](https://img-blog.csdnimg.cn/20211005030802669.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBA5peg5Zyo5peg5LiN5Zyo,size_20,color_FFFFFF,t_70,g_se,x_16)