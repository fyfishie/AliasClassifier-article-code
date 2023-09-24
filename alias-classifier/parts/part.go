package parts

import (
	"alias_article/lib"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fyfishie/dorapock/store"
	"golang.org/x/exp/rand"
	"os"
	"strconv"
)

var vpList = map[string]string{
	"1": "./db/database/distect/36.138.22.160/36.138.22.160.1/",
	"2": "./db/database/distect/36.140.40.41/36.140.40.41.1/",
	"3": "./db/database/distect/36.140.14.122/36.140.14.122.1/",
	"4": "./db/database/distect/london.tect/london.1/",
	"5": "./db/database/distect/newark/newark.1/",
	"6": "./db/database/distect/qingdao/qingdao.1/",
}

//var Num_w = []int{1, 10, 50, 100, 200}

var Num_w = []int{1, 10, 50, 100, 200, 500, 800, 1000}

// select some ip randomly
func Part(interPath string, wtDir string) {
	allIP := []string{}
	rfi, err := os.OpenFile(interPath, os.O_RDONLY, 0000)
	if err != nil {
		panic(err)
	}
	scaner := bufio.NewScanner(rfi)
	defer rfi.Close()
	for scaner.Scan() {
		allIP = append(allIP, scaner.Text())
	}
	rand.Seed(114514)
	l := len(allIP)
	for _, num := range Num_w {
		os.Mkdir("./db/database/shanghai_parts/"+strconv.Itoa(num), 0777)
		wfi, err := os.OpenFile("./db/database/shanghai_parts/"+strconv.Itoa(num)+"/ip", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		if err != nil {
			panic(err)
		}
		intersMap := map[string]struct{}{}
		for i := 0; i < num*10000; i++ {
			for {
				ip := allIP[rand.Intn(l)]
				if _, ok := intersMap[ip]; ok {
					continue
				}
				intersMap[ip] = struct{}{}
				break
			}
		}
		wtr := bufio.NewWriter(wfi)
		for ip := range intersMap {
			wtr.WriteString(ip + "\n")
		}
		wtr.Flush()
		wfi.Close()
	}
}
func Trace() {
	for i, num := range Num_w {
		fmt.Printf("start %v\n", i)
		dir := "./db/database/shanghai_parts/" + strconv.Itoa(num)
		ipPath := dir + "/ip"
		done := make(chan bool)
		for index, vp := range vpList {
			origin := vp + "shanghai_smark_4.json"
			wtPath := dir + "/trace_" + index

			go trace(ipPath, origin, wtPath, index, done)
		}
		for i := 0; i < len(vpList); i++ {
			<-done
		}
	}
}

// extract trace data from local data, to reduce time cost of experiments
func trace(ipPath string, originTracePath string, wtPath string, id string, done chan bool) {
	allIP := map[string]struct{}{}
	rfi, err := os.OpenFile(ipPath, os.O_RDONLY, 0000)
	if err != nil {
		panic(err)
	}
	scaner := bufio.NewScanner(rfi)
	defer rfi.Close()
	for scaner.Scan() {
		allIP[scaner.Text()] = struct{}{}
	}
	ldr := store.NewLoader[lib.RawTrace](originTracePath)
	err = ldr.Open()
	if err != nil {
		panic(err)
	}
	wfi, err := os.OpenFile(wtPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	defer wtr.Flush()
	cnt := 0
	for ldr.HasNext() {
		cnt++
		if cnt%1000000 == 0 {
			fmt.Printf("%v:%v\n", id, cnt)
		}
		t := ldr.Next()
		if _, ok := allIP[t.Ip]; ok {
			if len(t.Results) != 0 {
				if t.Ip == t.Results[len(t.Results)-1].Ip {
					bs, _ := json.Marshal(t)
					wtr.Write(bs)
					wtr.WriteString("\n")
				}
			}

		}
	}
	done <- true
}
