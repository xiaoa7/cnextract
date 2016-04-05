//基于分词的信息提取
package main

import (
	cn "cnextract"
	"encoding/json"
	"github.com/huichen/sego"
	"net/http"
	"strings"
)

var segret map[string]interface{}
var result cn.SegResult
var segmenter sego.Segmenter

//加载分词字典
func init() {
	segmenter.LoadDictionary("dictionary.txt")
}

//分词
func DoSegmentation(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	txt := r.FormValue("txt")
	txt = strings.Replace(txt, " ", "", -1)
	//
	segments := segmenter.Segment([]byte(txt))
	//显示用
	result = cn.SegResult{}
	result.TakeTime = 0
	result.SegWords = []*cn.SegWord{}
	ret := ""
	for _, v := range segments {
		result.SegWords = append(result.SegWords, &cn.SegWord{
			Frequency: v.Token().Frequency(),
			Offset:    v.Start(),
			Nature:    v.Token().Pos(),
			Word:      v.Token().Text(),
		})
		ret += "[" + v.Token().Pos() + "]" + v.Token().Text() + " "
	}
	w.Write([]byte(ret))
}

//抽取
func DoExtract(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rulestr := r.FormValue("rule")
	rule := cn.Rule{}
	rule.ParseRule(rulestr)
	ret := cn.ExtractTest(rule, result)
	bs, _ := json.Marshal(ret)
	w.Write(bs)
}
func main() {
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.HandleFunc("/seg", DoSegmentation)
	http.HandleFunc("/extract", DoExtract)
	http.ListenAndServe(":80", nil)
}
