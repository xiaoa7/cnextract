//基于分词的信息提取
package main

import (
	cn "cnextract"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	HanlpAPI = "http://192.168.3.14:10800/analyse/2"
)

var segret map[string]interface{}
var result cn.HanlpResult

//分词
func DoSegmentation(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	txt := r.FormValue("txt")
	txt = strings.Replace(txt, " ", "", -1)
	vs := make(url.Values)
	vs.Add("op", ",jg,")
	vs.Add("text", txt)
	resp, _ := http.PostForm(HanlpAPI, vs)
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bs, &result)
	w.Write(bs)
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
