package cnextract

import (
	"log"
	"strings"
	"time"
)

type SegWord struct {
	Frequency int
	Offset    int
	Nature    string
	Word      string
}
type SegResult struct {
	TakeTime int
	SegWords []*SegWord
}

//提取
func Extract(content string) (ret map[string]string) {
	ret = make(map[string]string)
	//先提取中文
	t1 := time.Now()
	//TODO
	result := SegResult{}

	log.Printf("分词耗时:%.6f秒\n", time.Now().Sub(t1).Seconds())
	t1 = time.Now()
	//TODO 2中查找模式，现在用第一种
	//开始识别,迭代所有规则
	resultlen := len(result.SegWords)
	for name, v := range rules {
		//识别一个规则文件
		//log.Println("识别规则", name, v)
		retarr := []string{}
	L1:
		for _, v1 := range v.Subrules {
			//识别一条规则
			//L2:
			for index, v2 := range result.SegWords {
				//验证一条规则
				ok := false
				//按词性校验
				if len(v1.StartWordsNature) > 0 {
					for _, v3 := range v1.StartWordsNature {
						if v3 == v2.Nature {
							ok = true
							break
						}
					}
					if !ok { //当前词匹配开始词不对
						continue
					}
				}
				//按词识别
				ok = false
				for _, v3 := range v1.StartWords {
					if v3.IsMust && v3.Txt == v2.Word {
						ok = true
						break
					} else if v3.IsPrefix && strings.HasPrefix(v2.Word, v3.Txt) {
						ok = true
						break
					} else if v3.IsSuffix && strings.HasSuffix(v2.Word, v3.Txt) {
						ok = true
						break
					}
				}
				if !ok {
					continue
				}
				//向后找满足条件的词
				ok = false
			L3:
				for i := index + v1.BreakWordsLen[0]; i+index < resultlen && i <= index+v1.BreakWordsLen[1]; i++ {
					tmpword := result.SegWords[i]
					check := false
					//验证选择词的词性
					if len(v1.ExtractWordNature) > 0 {
						for _, v4 := range v1.ExtractWordNature {
							//log.Println(name, v4, tmpword.Nature, tmpword.Word)
							if v4 == tmpword.Nature {
								//log.Println("找到了")
								check = true
								break
							}
						}
						if !check {
							//log.Println("提取词词性不符")
							continue L3
						}
					}
					//TODO 验证词类型,用正则验证即可，先不做了
					//验证找到的词，是字母、数字、英文、中文
					//
					//验证长度
					if len(tmpword.Word) >= v1.ExtractWordLen[0] && len(tmpword.Word) <= v1.ExtractWordLen[1] {
						//找到词
						if v1.ExtractWortMaxLen == 0 {
							retarr = append(retarr, tmpword.Word)
						} else {
							//找结束词
							findstr := []string{tmpword.Word}
						L5:
							for j := i + 1; j <= v1.ExtractWortMaxLen+i; j++ {
								tmpword2 := result.SegWords[j]
								findstr = append(findstr, tmpword2.Word)
								//验证结束词
								ok2 := false
								if len(v1.EndStopWordsNature) > 0 {
									//log.Println(v1.EndStopWordsNature)
									for _, v5 := range v1.EndStopWordsNature {
										if tmpword2.Nature == v5 {
											ok2 = true
										}
									}
									if !ok2 {
										continue L5
									}
								}
								ok2 = false
								if len(v1.EndStopWords) > 0 {
									//验证
									for _, v5 := range v1.EndStopWords {
										//log.Println(v5.Txt, tmpword2.Word)
										if v5.IsMust && v5.Txt == tmpword2.Word {
											ok2 = true
											break
										} else if v5.IsPrefix && strings.HasPrefix(tmpword2.Word, v5.Txt) {
											ok2 = true
											break
										} else if v5.IsSuffix && strings.HasSuffix(tmpword2.Word, v5.Txt) {
											ok2 = true
											break
										}
									}
								} else {
									ok2 = true
								}
								if !ok2 {
									continue L5
								} else {
									break L5
								}
							}
							retarr = append(retarr, strings.Join(findstr[:len(findstr)-1], ""))
						}
						break L1
					}
				}
				//
			}
		}
		//
		ret[name] = strings.Join(retarr, ",")
	}
	log.Printf("提取耗时:%.6f秒\n", time.Now().Sub(t1).Seconds())
	return
}

//检测词是否匹配
func checkWord(word string, ws []Word) bool {
	for _, v := range ws {
		if v.IsMust && v.Txt == word {
			return true
		} else if v.IsPrefix && strings.HasPrefix(word, v.Txt) {
			return true
		} else if v.IsSuffix && strings.HasSuffix(word, v.Txt) {
			return true
		} else if v.IsMiddle && strings.Index(word, v.Txt) > -1 {
			return true
		}
	}
	return false
}

//判断字符串，是否在一个数组中
func checkInArr(s1 string, s2 []string) bool {
	for _, v := range s2 {
		if s1 == v {
			return true
		}
	}
	return false
}

//提取
func ExtractTest(v Rule, result SegResult) (ret map[string]string) {
	ret = make(map[string]string)
	//先提取中文
	resultlen := len(result.SegWords)
	retarr := []string{}
L1:
	for _, v1 := range v.Subrules {
		//识别一条规则
	L2:
		for index, v2 := range result.SegWords {
			//验证一条规则
			ok := false
			//按词性校验开始词
			if len(v1.StartWordsNature) > 0 && !checkInArr(v2.Nature, v1.StartWordsNature) {
				//log.Println("开始词词性校验", v1.StartWordsNature, v2.Nature)
				continue L2
			}
			//按词识别
			ok = checkWord(v2.Word, v1.StartWords)
			if !ok {
				//log.Println("开始词校验", v1.StartWords, v2.Word)
				continue
			}
			log.Println("开始词", v2.Word, "验证成功", v1.BreakWordsLen, index, resultlen, v1.ExtractWordNature)
			//向后找满足条件的词
			ok = false
		L3:
			for i := index + v1.BreakWordsLen[0]; i < resultlen && i <= index+v1.BreakWordsLen[1]; i++ {
				tmpword := result.SegWords[i]
				//验证选择词的词性
				if len(v1.ExtractWordNature) > 0 && !checkInArr(tmpword.Nature, v1.ExtractWordNature) {
					continue L3
				}
				//TODO 验证词类型,用正则验证即可，先不做了
				//验证找到的词，是字母、数字、英文、中文
				//
				log.Println(tmpword.Word, len(tmpword.Word), v1.ExtractWordLen[0], v1.ExtractWordLen[1], v1.ExtractWortMaxLen)
				//验证长度
				if len(tmpword.Word) >= v1.ExtractWordLen[0] && len(tmpword.Word) <= v1.ExtractWordLen[1] {
					log.Println("找到词", tmpword.Word)
					//找到词
					if v1.ExtractWortMaxLen == 0 {
						retarr = append(retarr, tmpword.Word)
					} else {
						//找结束词
						findstr := []string{tmpword.Word}
					L5:
						for j := i + 1; j <= v1.ExtractWortMaxLen+i; j++ {
							tmpword2 := result.SegWords[j]
							findstr = append(findstr, tmpword2.Word)
							//验证结束词
							ok2 := false
							if len(v1.EndStopWordsNature) > 0 && !checkInArr(tmpword2.Nature, v1.EndStopWordsNature) {
								continue L5
							}
							ok2 = false
							if len(v1.EndStopWords) > 0 {
								ok2 = checkWord(tmpword2.Word, v1.EndStopWords)
							} else {
								ok2 = true
							}
							if !ok2 {
								continue L5
							} else {
								break L5
							}
						}
						retarr = append(retarr, strings.Join(findstr[:len(findstr)-1], ""))
					}
					break L1
				}
			}
			//
		}
	}
	//
	ret["测试规则"] = strings.Join(retarr, ",")
	return
}
