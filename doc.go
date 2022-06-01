package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/liumingmin/goutils/utils"
)

func main() {
	var sb strings.Builder

	prefix, _ := ioutil.ReadFile(filepath.Join(utils.GetCurrPath(), "BASE-README.MD"))
	sb.WriteString(string(prefix))

	genDocByTestFile(utils.GetCurrPath(), 1, &sb)
	ioutil.WriteFile("README.md", []byte(sb.String()), 0666)

	if err := exec.Command("cmd", "/c", "markdown-toc --maxdepth 3 -i README.md").Run(); err != nil {
		fmt.Println(err)
	}
}

var moduleCnName = map[string]string{
	"cache":                   "缓存模块",
	"mem_cache_test.go":       "内存缓存",
	"rds_cache_test.go":       "Redis缓存",
	"conf":                    "yaml配置模块",
	"container":               "容器模块",
	"bitmap_test.go":          "比特位表",
	"const_hash_test.go":      "一致性HASH",
	"lighttimer_test.go":      "轻量级计时器",
	"db":                      "数据库",
	"elasticsearch":           "ES搜索引擎",
	"es6":                     "ES6版本API",
	"es7":                     "ES7版本API",
	"kafka":                   "kafka消息队列",
	"mongo":                   "mongo数据库",
	"redis":                   "go-redis",
	"log":                     "zap日志库",
	"middleware":              "中间件",
	"captcha_test.go":         "验证码模块",
	"limit_conn_test.go":      "限连接模块",
	"limit_req_test.go":       "限流模块",
	"service_handler_test.go": "service封装器",
	"thumb_image_test.go":     "缩略图",
	"net":                     "网络库",
	"httpx":                   "兼容http1.x和2.0的httpclient",
	"packet":                  "tcp包model",
	"proxy":                   "ssh proxy",
	"serverx":                 "兼容http1.x和2.0的http server",
	"utils":                   "通用工具库",
	"buffer_invoker":          "异步调用",
	"cbk":                     "熔断器",
	"csv":                     "CSV文件解析为MDB内存表",
	"distlock":                "分布式锁",
	"fsm":                     "有限状态机",
	"hc":                      "httpclient工具",
	"ismtp":                   "邮件工具",
	"safego":                  "安全的go协程",
	"ws":                      "websocket客户端和服务端库",
	"docgen":                  "文档自动生成",
	"crc16_test.go":           "crc16算法",
	"descartes_test.go":       "笛卡尔组合",
	"list_test.go":            "Redis List工具库",
	"zset_test.go":            "Redis ZSet工具库",
	"mq_test.go":              "Redis PubSub工具库",
	"lock_test.go":            "Redis 锁工具库",
	"tags_test.go":            "结构体TAG生成器",
	"snowflake_test.go":       "雪花ID生成器",
}

//dir := filepath.Dir(filePath)
func genDocByTestFile(dir string, level int, sb *strings.Builder) {
	files, _ := ioutil.ReadDir(dir)

	nextLevel := level + 1

	for _, file := range files {
		if file.IsDir() {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}
			genDirLevel(file.Name(), nextLevel, sb)
			genDocByTestFile(filepath.Join(dir, file.Name()), nextLevel, sb)
			continue
		}

		if strings.HasSuffix(file.Name(), "_test.go") {
			codeFilePath := dir + "/" + file.Name()
			bs, err := ioutil.ReadFile(codeFilePath)
			if err != nil {
				continue
			}
			content := string(bs)

			genDirLevel(file.Name(), nextLevel, sb)
			parseTestCode(nextLevel, content, sb)
		}
	}
}

func genDirLevel(dirName string, level int, sb *strings.Builder) {
	prefixSymbol := ""
	for i := 0; i < level; i++ {
		prefixSymbol += "#"
	}

	dirCnName, ok := moduleCnName[dirName]
	if ok {
		dirCnName = dirName + " " + dirCnName
	} else {
		dirCnName = dirName
	}
	fmt.Println(dirName)
	sb.WriteString(fmt.Sprintf("%s %s\n", prefixSymbol, dirCnName))
}

func parseTestCode(level int, content string, sb *strings.Builder) {
	reg, _ := regexp.Compile(`(?U)func (?P<fname>.*)\(t \*testing\.T\) *\{(?P<body>(.|\n)*)\n\}`)
	match := reg.FindAllStringSubmatch(content, -1)

	for _, item := range match {
		genDirLevel(item[1], level+1, sb)

		sb.WriteString("```go\n")
		sb.WriteString(item[2] + "\n")
		sb.WriteString("```\n")
	}
}