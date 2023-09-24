package time_cost

import (
	"alias_article/data4sk"
	"alias_article/lib"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fyfishie/dorapock/store"
	"io"
	"net/http"
	"time"
)

// this function is write for experiment of time cost, where the AliasClassifier was asked to
// parse a mass of ip and we'll pick up the time cost result
func RfcPipeline(tracePath, spingpath, domainPath string, allIP []string, wtr *bufio.Writer) int {
	spings, err := store.LoadAny[lib.Sping](spingpath)
	if err != nil {
		panic(err)
	}
	spingExist := map[string]struct{}{}
	for _, s := range spings {
		spingExist[s.IP] = struct{}{}
	}
	g := data4sk.DataGener{
		Field: data4sk.FieldChoice{
			LengthDiff:          true,
			DirectDiff:          true,
			TraceSameFactor:     true,
			RttGap:              true,
			ReplyTtl:            true,
			TopDomainConsist:    true,
			SecondDomainConsist: true,
			SubDomainConsist:    true,
			IPDistance:          true,
			DomainDistance:      true,
		},
		SmarkPath:  tracePath,
		SpingPath:  spingpath,
		DomainPath: domainPath,
	}
	aliasCnt := 0
	fmt.Println("load data")
	g.LoadData()
	buf := [][]int{}
	start := time.Now().Unix()
	fmt.Println("start pair")
	total := len(allIP)
	for i := 0; i < len(allIP); i++ {
		for j := i + 1; j < len(allIP); j++ {
			if _, ok := spingExist[allIP[i]]; ok {
				if _, ok := spingExist[allIP[j]]; ok {
					t := g.OnePair(allIP[i], allIP[j])
					pairRecord := []lib.Pair{}
					pairRecord = append(pairRecord, lib.Pair{IPA: allIP[i], IPB: allIP[j]})
					buf = append(buf, t)
					if len(buf) > 40000 {
						fmt.Printf("i:%v\n", i)
						mid := time.Now().Unix()
						minute := int((mid - start) / 60)
						fmt.Printf("time used (minute): %v\n", minute)
						if (2*total-i)*i != 0 {
							totalTime := (total * total) / ((2*total - i) * i) * minute
							fmt.Printf("pre total minute: %v\n", totalTime)
						}
						fmt.Println("==========================\n")
						indexs := predict(buf)
						for _, index := range indexs {
							pair := pairRecord[index]
							bs, _ := json.Marshal(pair)
							wtr.Write(bs)
							wtr.WriteString("\n")
						}
						buf = [][]int{}
						pairRecord = []lib.Pair{}
					}
				}
			}
		}
	}
	return aliasCnt
}

// send one block of charactor vectors and receives predict result
func predict(data [][]int) []int {
	bs, _ := json.Marshal(data)
	resp, err := http.Post("http://127.0.0.1:8981/", "application/json", bytes.NewReader(bs))
	if err != nil {
		panic(err)
	}
	rbs, _ := io.ReadAll(resp.Body)
	ires := []int{}
	err = json.Unmarshal(rbs, &ires)
	res := []int{}
	for index, ir := range ires {
		if ir == 1 {
			res = append(res, index)
		}
	}
	return res
}
