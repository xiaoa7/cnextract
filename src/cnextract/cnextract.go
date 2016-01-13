//中文数据抽取
//从一篇文本中抽取想要的关系数据
//基于Hanlp中文分词开源项目
package cnextract

import (
	"regexp"
)

const (
	HanlpAPI = "http://192.168.3.14:10800/analyse/2"
)

//
var rules map[string]Rule
var chinesereg *regexp.Regexp

//初始化
func init() {
	rules = make(map[string]Rule)
	chinesereg, _ = regexp.Compile("[\\p{Han}]+")
}
