package pipeline

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fyfishie/dorapock/store"
	"io"
	"mlar/character"
	"mlar/lib"
	"mlar/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

var vpList = map[string]string{
	"1":       "./database/distect/36.138.22.160/36.138.22.160.1/",
	"2":       "./database/distect/36.140.40.41/36.140.40.41.1/",
	"3":       "./database/distect/36.140.14.122/36.140.14.122.1/",
	"london":  "./database/distect/london.detect/london.1/",
	"newark":  "./database/distect/newark/newark.1/",
	"qingdao": "./database/distect/qingdao/qingdao.1/",
}

type Inter struct {
	ipTraceMap  map[string]lib.RawTrace
	reachableIP map[string]struct{}
	ipIspMap    map[string]string
}

func (I *Inter) FiltInit(tracePath string, ispPath string, spingPath string) {
	traces, err := store.LoadAny[lib.RawTrace](tracePath)
	if err != nil {
		panic(err)
	}
	im := map[string]lib.RawTrace{}
	for _, t := range traces {
		im[t.Ip] = t
	}
	I.ipTraceMap = im
	ream := map[string]struct{}{}
	spings, err := store.LoadAny[lib.Sping](spingPath)
	if err != nil {
		panic(err)
	}
	for _, s := range spings {
		ream[s.IP] = struct{}{}
	}
	isps, err := store.LoadAny[lib.ISP](ispPath)
	if err != nil {
		panic(err)
	}
	ipIspMap := map[string]string{}
	for _, is := range isps {
		ipIspMap[is.IP] = is.Isp
	}

	return
}
func (I *Inter) Filt(ipA, ipB string) bool {
	_, o1 := I.reachableIP[ipA]
	_, o2 := I.reachableIP[ipB]
	if o1 != o2 {
		return false
	}
	ispA, o1 := I.ipIspMap[ipA]
	ispB, o2 := I.ipIspMap[ipB]
	if o1 && o2 {
		if ispA != ispB {
			return false
		}
	}
	if tA, ok := I.ipTraceMap[ipA]; ok {
		if tB, ok := I.ipTraceMap[ipB]; ok {
			if utils.LenDiff(tA, tB) > 3 {
				return false
			}
			if utils.DirDistance(tA, tB) > 2 {
				return false
			}
		}
	}
	return true
}

// it generate character vectors in pipeline way and predicts them block by block
func (I *Inter) PipeLineMlar(traceDir string, ipPath string, wtpath string, whoisPath string) (cnt int, used int, pre_time int) {
	start := time.Now().Unix()
	wfi, err := os.OpenFile(wtpath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	wtr := bufio.NewWriter(wfi)
	defer wfi.Close()
	tracesPath := []string{}
	for i := 0; i < 6; i++ {
		tracesPath = append(tracesPath, traceDir+"/trace_"+strconv.Itoa(i+1))
	}
	fmt.Println("rtt same init")
	character.RttSameFactor4PipelineInit(tracesPath, ipPath)
	fmt.Println("trace init")
	character.TraceFactor4PipelineInit(tracesPath)
	allIP := []string{}
	fi, err := os.OpenFile(ipPath, os.O_RDONLY, 0000)
	if err != nil {
		panic(err)
	}
	scaner := bufio.NewScanner(fi)
	for scaner.Scan() {
		allIP = append(allIP, scaner.Text())
	}
	pairs := []lib.Pair{}
	total := len(allIP)
	for i := 0; i < len(allIP); i++ {
		now := time.Now().Unix()
		if now-start > lib.Hour*10 {
			minute := int((now - start) / 60)
			totalTime := (total * total) / ((2*total - i) * i) * minute
			return cnt, minute, totalTime
		}
		for j := i + 1; j < len(allIP); j++ {
			if !I.Filt(allIP[i], allIP[j]) {
				continue
			}
			pair := lib.Pair{IPA: allIP[i], IPB: allIP[j]}
			pairs = append(pairs, pair)
			if len(pairs) > 10000 {
				fmt.Println(fmt.Sprintf("%v/%v", i, len(allIP)))
				mid := time.Now().Unix()
				minute := int((mid - start) / 60)
				if minute != 0 && i != 0 {
					totalTime := (total * total) / ((2*total - i) * i) * minute
					fmt.Println(fmt.Sprintf("city %v,minute %v", "shanghai", minute))
					fmt.Println(fmt.Sprintf("total time,%v\n============\n", totalTime))
				}
				f1 := character.RttSameFactor4City(pairs)
				f2 := character.TraceFactor4City(pairs)
				//f3 := character.Whois4City(pairs)
				data := [][2]int{}
				for _, p := range pairs {
					data = append(data, [2]int{f1[p.ID()], f2[p.ID()]})
				}
				res := predict(data)
				for i, r := range res {
					if r {
						p := pairs[i]
						bs, _ := json.Marshal(p)
						wtr.Write(bs)
						wtr.WriteString("\n")
					}
				}
				pairs = []lib.Pair{}
			}
		}
	}
	wtr.Flush()
	return cnt, -1, -1
}

// send characters to server and return predict result
func predict(data [][2]int) []bool {
	fmt.Println("predict")
	bs, _ := json.Marshal(data)
	resp, err := http.Post("http://localhost:8982/", "application/json", bytes.NewReader(bs))
	fmt.Println("resp received")
	if err != nil {
		panic(err)
	}
	ires := []int{}
	rbs, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(rbs, &ires)
	if err != nil {
		panic(err)
	}
	res := []bool{}
	for _, in := range ires {
		if in == 0 {
			res = append(res, false)
		} else {
			res = append(res, true)
		}
	}
	return res
}
