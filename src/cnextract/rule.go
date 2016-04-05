package cnextract

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	//要提取的词类型，后续要做正则验证
	WORD_TYPE_NOLIMIT = iota //不显，此时要用词性限制
	WORD_TYPE_DIGITAL        //纯数字
	WORD_TYPE_EN             //英文
	WORD_TYPE_CN             //中文
)

//单个词
type Word struct {
	Txt      string
	IsMust   bool //绝对匹配
	IsPrefix bool //前缀匹配
	IsSuffix bool //后缀匹配
	IsMiddle bool //中间匹配
}

//单行规则
type SubRule struct {
	StartWords         []Word   //开始词
	StartWordsNature   []string //开始词要求词性
	BreakWordsLen      [2]int   //间隔词数范围
	ExtractWordType    int      //提取词类型
	ExtractWordNature  []string //提取词词性
	ExtractWordLen     [2]int   //提取词长度
	ExtractWortMaxLen  int      //提取词最多抓几个词语
	EndStopWords       []Word   //结束词标志
	EndStopWordsNature []string //结束词词性

}

//文件规则
type Rule struct {
	Name     string
	Txt      string //规则的文本文件
	Subrules []SubRule
}

//从一个目录中加载规则
func LoadRules(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		filename := info.Name()
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(filename, ".rul") {
			return nil
		}
		fi, _ := os.Open(path)
		bs, _ := ioutil.ReadAll(fi)
		fi.Close()
		r := Rule{}
		r.Name = filename[:strings.LastIndex(filename, ".")]
		r.Txt = string(bs)
		tmp := strings.Split(r.Txt, "\r\n")
		//TODO 解析每一行规则
		for _, v := range tmp {
			r.ParseRule(v)
		}
		rules[r.Name] = r
		return nil
	})

}

//解析词
func parseWord(src string) []Word {
	tmp2 := strings.Split(src, "|")
	words := []Word{}
	for _, v := range tmp2 {
		var w Word
		if strings.HasPrefix(v, "*") && strings.HasSuffix(v, "*") {
			w = Word{Txt: v[1 : len(v)-1], IsMiddle: true}
		} else if strings.HasPrefix(v, "*") {
			w = Word{Txt: v[1:], IsPrefix: true}
		} else if strings.HasSuffix(v, "*") {
			w = Word{Txt: v[:len(v)-1], IsSuffix: true}
		} else {
			w = Word{Txt: v, IsMust: true}
		}
		words = append(words, w)
	}
	return words
}

//解析取值范围
func parseRange(src string) [2]int {
	ret := [2]int{}
	tmp2 := strings.Split(src, "-")
	if len(tmp2) == 2 {
		start, _ := strconv.Atoi(tmp2[0])
		end, _ := strconv.Atoi(tmp2[1])
		ret[0] = start
		ret[1] = end
	}
	return ret
}

//解析规则
func (r *Rule) ParseRule(line string) {
	tmp := strings.Split(line, " ")
	if len(tmp) != 6 && len(tmp) != 9 {
		return
	}
	subrule := SubRule{}
	//处理开始词
	subrule.StartWords = parseWord(tmp[0])
	//处理开始次词性
	if tmp[1] != "*" {
		subrule.StartWordsNature = strings.Split(tmp[1], "|")
	}
	//处理间隔长度
	subrule.BreakWordsLen = parseRange(tmp[2])
	//处理要求词类型
	subrule.ExtractWordType, _ = strconv.Atoi(tmp[3])
	//处理要求词词性
	if tmp[4] != "*" {
		subrule.ExtractWordNature = strings.Split(tmp[4], "|")
	}
	//处理要求次长度
	subrule.ExtractWordLen = parseRange(tmp[5])
	if len(tmp) == 9 {
		//提取的最大长度
		subrule.ExtractWortMaxLen, _ = strconv.Atoi(tmp[6])
		//结束停词
		if tmp[7] != "*" {
			subrule.EndStopWords = parseWord(tmp[7])
		}
		//结束停词，词性
		if tmp[8] != "*" {
			subrule.EndStopWordsNature = strings.Split(tmp[8], "|")
		}
	}
	r.Subrules = append(r.Subrules, subrule)
}
