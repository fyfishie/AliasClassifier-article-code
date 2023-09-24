/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-08-05:07
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package main

import (
	"bufio"
	"fmt"
	"mlar/character"
	"mlar/lib"
	"mlar/pipeline"
	"os"
	"strconv"
	"time"

	"github.com/fyfishie/dorapock/store"
)

var tracePaths [5][]string
var spingPaths [5][]string
var pairsPaths [5]string
var upairsPaths [5]string
var uinterPaths [5]string
var interfaceIPPaths [5]string

func main() {
	parts()
	oobDataGen()
}

// generate test vectors in pipeline way for time cost statistic
func parts() {
	wfi, err := os.OpenFile("./database/log/shanghai_parts_mlar_pipeline-sp.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	for _, num := range lib.Num_w[1:] {
		start := time.Now().Unix()
		I := pipeline.Inter{}
		I.FiltInit("./database/parts/"+strconv.Itoa(num)+"/traces/trace_1", "./database/parts/"+strconv.Itoa(num)+"/isp", "./database/parts/"+strconv.Itoa(num)+"/spings/sping_1")
		cnt, minute, pre_time := I.PipeLineMlar("./database/parts/"+strconv.Itoa(num)+"/traces", "./database/parts/"+strconv.Itoa(num)+"/ip", "./database/parts/"+strconv.Itoa(num)+"/alias.pairs", "./database/parts/"+strconv.Itoa(num)+"/ip.whois")
		if minute != -1 {
			wfi.WriteString(fmt.Sprintf("city: %v, cnt: %v, minute used: %v, pre_time:%v\n", "shanghai", cnt, minute, pre_time))
		}
		end := time.Now().Unix()
		wfi.WriteString(fmt.Sprintf("%v,%v\n", "parts", end-start))
		wfi.WriteString("num split ================================\n")
	}
	wfi.Close()
}

// generate train data or test data for machine learning
func oobDataGen() {
	wfi, err := os.OpenFile("./database/log/mlar-int-2-factor-test-gen.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	defer wtr.Flush()
	for i := 0; i < 5; i++ {
		start := time.Now().Unix()
		mlar_data_gen(pairsPaths[i], tracePaths[i], "./database/oob/"+strconv.Itoa(i+1)+"/mlar/train_data-2-factor-int.csv", true, true)
		mlar_data_gen(upairsPaths[i], tracePaths[i], "./database/oob/"+strconv.Itoa(i+1)+"/mlar/train_data-2-factor-int.csv", false, false)
		end := time.Now().Unix()
		wtr.WriteString(fmt.Sprintf("%v:%v\n", i+1, end-start))
	}
}
func init() {
	for i := 1; i < 6; i++ {
		oobinit(i)
	}
}
func oobinit(oobNum int) {
	distectDir := "./database/distect/"
	fs, err := os.ReadDir(distectDir)
	if err != nil {
		panic(err)
	}
	tracepaths := []string{}
	for _, f := range fs {
		if f.Name() == "36.138.22.160" || f.Name() == "newark" || f.Name() == "qingdao" || f.Name() == "london.text" {
			continue
		}
		path := "./database/distect/" + f.Name() + "/allres/smark/all_4_" + strconv.Itoa(oobNum) + ".json"
		tracepaths = append(tracepaths, path)
	}
	spingpaths := []string{}
	for _, f := range fs {
		path := "./database/distect/" + f.Name() + "/allres/sping/all_4_" + strconv.Itoa(oobNum) + ".json"
		spingpaths = append(spingpaths, path)
	}
	tracePaths[oobNum-1] = tracepaths
	spingPaths[oobNum-1] = spingpaths
	pairp := "./database/oob/" + strconv.Itoa(oobNum) + "/mlar/pairs.json"
	pairsPaths[oobNum-1] = pairp
	interp := "./database/oob/" + strconv.Itoa(oobNum) + "/mlar/interfaceIP.json"
	upairsPaths[oobNum-1] = "./database/oob/" + strconv.Itoa(oobNum) + "/train_upairs.json"
	uinterPaths[oobNum-1] = "./database/oob/" + strconv.Itoa(oobNum) + "/mlar/train_uinter.json"
	interfaceIPPaths[oobNum-1] = interp
}

func mlar_data_gen(pairsPath string, tracePaths []string, wtpath string, head bool, pair bool) {
	g := character.NewGener(tracePaths, pairsPath, wtpath)
	rttSames, traceSames := g.Run()

	//nres3 := character.Whois("./database/itdk/midar-iff.nodes.cleaned.ip.whois", upairsPaths[oobNum-1])
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	defer wtr.Flush()
	if head {
		wtr.WriteString("rtt_same,trace_same,alias\n")
	}
	pairs, err := store.LoadAny[lib.Pair](pairsPath)
	if err != nil {
		panic(err)
	}
	for _, p := range pairs {
		if rttSames[p.ID()] == -1 || traceSames[p.ID()] == -1 {
			continue
		}
		if pair {
			wtr.WriteString(fmt.Sprintf("%d,%d,%v\n", rttSames[p.ID()], traceSames[p.ID()], 1))
		} else {
			wtr.WriteString(fmt.Sprintf("%d,%d,%v\n", rttSames[p.ID()], traceSames[p.ID()], 0))
		}
	}
	wtr.Flush()
}
