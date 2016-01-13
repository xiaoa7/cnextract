package cnextract

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	//加载配置
	LoadRules("../rules")
	//
	fi, err := os.Open("../1.txt")
	if err != nil {
		log.Println(err.Error())
	}
	defer fi.Close()
	bs, _ := ioutil.ReadAll(fi)
	ret := Extract(HanlpAPI, string(bs))
	log.Println(ret)
}
