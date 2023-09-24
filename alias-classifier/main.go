/*
 * @Author: fyfishie
 * @Date: 2023-04-18:07
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-08-04:15
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package main

import (
	"alias_article/data4sk"
	"alias_article/filter"
	"alias_article/lib"
	"alias_article/parts"
	"alias_article/showup_probability"
	"alias_article/statistic"
	"alias_article/time_cost"
	"alias_article/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/fyfishie/dorapock/store"
)

func main() {
	sta()
	dataGen()
	statistic.IpDistribute()
	showup()
	statistic.MakeIPIDArray("./db/origin/detect/ipid/origin/midar-iff.nodes.cleaned.ipid.withoutlocal.10packs.json", "./db/origin/detect/ipid/dewarp/midar-iff.nodes.cleaned.ipid.withoutlocal.10packs.json")
	findBest()
	argsSort()
	data4sk.Merge2("./db/database/multisk/merged012.csv", "./db/database/multisk/sk_data3.csv", "./db/merged0123.csv", "2")
	oobPairs()
	dataGen()
	parts.Part("./db/database/city/shanghai/ipv4.txt", "")
	parts.Trace()
	rfc_pipe()
	makeNB()
}

// use this top function to generate charactor vectors for machine learning
func dataGen() {
	g := data4sk.DataGener{
		Alias:     true,
		WriteHead: true,
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
		SmarkPath:  "./db/database/oob/5/trace.json",
		SpingPath:  "./db/database/oob/5/sping.json",
		PairsPath:  "./db/database/oob/5/test_paris.json",
		DomainPath: "./db/database/oob/5/rdns.json",
		SavePath:   "./db/database/oob/5/train_sk-sp.csv",
	}
	g.Run()
	g = data4sk.DataGener{
		Alias:     false,
		WriteHead: false,
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
		SmarkPath:  "./db/database/oob/5/trace.json",
		SpingPath:  "./db/database/oob/5/sping.json",
		PairsPath:  "./db/database/oob/5/test_uparis.json",
		DomainPath: "./db/database/oob/5/rdns.json",
		SavePath:   "./db/database/oob/5/train_sk-sp.csv",
	}
	g.Run()
}

// use this top function to make alias pairs and anti-alias pairs from itdk data
func makePair() {
	utils.MakeUnPair("./db/origin/itdk/midar-iff.nodes.cleaned", "./db/pairs/unpairs-250.json")
	//statistic.PairDiffSta("./db/database/distect/36.103.226.176/36.103.226.176/allres/smark/all_4_1.json","./db/pairs/all_pairs-host-port-whois.json",)
	utils.AllPair("./db/origin/itdk/midar-iff.nodes.cleaned", "./db/pairs/all_pairs.json")
	filter.HostFilter("./db/pairs/all_pairs.json", "./db/origin/detect/sping/london/resultmidar-iff.nodes.cleaned.ip.londonmachine.spingres.0", "./db/pairs/all_pairs-host.json")
	filter.PortsFilt("./db/origin/detect/port/midar-iff.nodes.cleaned.ip.localmachine.json", "./db/pairs/all_pairs-host.json", "./db/pairs/all_pairs-host-port.json")
	filter.WhoisFilter("./db/pairs/all_pairs-host-port.json", "./db/pairs/all_pairs-host-port-whois.json", "./db/origin/detect/midar-iff.nodes.cleaned.whoisinfo")
}

// use this top function to make statistic data which is used to draw statistic pictures
func sta() {
	statistic.TraceDistanceSta("./db/database/distect/36.103.226.176/36.103.226.176/allres/smark/all_4_1.json", "./db/distect/")
	utils.TrimIPFromPair("./db/pairs/all_pairs-host-port-whois.json")
	pairDiffSta()
	pairFactor()
	rttGapSta()
	ttlGapSta()
	domainConsistSta()
	domainDistance()
	netSec()
}

// use this top function to calculates scale of some meta information
func showup() {
	showup_probability.PortScale("./db/origin/detect/port/midar-iff.nodes.cleaned.ip.localmachine.json", "./db/pairs/all_pairs-host-port-whois.json")
	showup_probability.IPID("./db/origin/detect/ipid/dewarp/midar-iff.nodes.cleaned.ipid.withoutlocal.10packs.json", "./db/pairs/all_pairs-host-port-whois.json")
	showup_probability.DNSScale("./db/origin/detect/itdk-run-20220224-dns-names.txt", "./db/pairs/all_pairs-host-port-whois.json")
	showup_probability.TTLScale("./db/origin/detect/sping/london/resultmidar-iff.nodes.cleaned.ip.londonmachine.spingres.0", "./db/pairs/all_pairs-host-port-whois.json")
	showup_probability.RttScale("./db/pairs/all_pairs-host-port-whois.json", "./db/origin/detect/smark/london/midar-iff.nodes.cleaned.ip.londonmachine.smarkres.0")
}

func pairDiffSta() {
	statistic.PairDiffSta("./db/database/distect/36.103.226.176/36.103.226.176/allres/smark/all_4_1.json",
		"./db/pairs/all_pairs-host-port-whois.json",
		"./db/result/pair/diff/length1.csv",
		"./db/result/pair/diff/length_scale1.csv",
		"./db/result/pair/diff/direct1.csv",
		"./db/result/pair/diff/direct_scale1.csv")
	statistic.PairDiffSta("./db/database/distect/36.103.226.176/36.103.226.176/allres/smark/all_4_1.json",
		"./db/pairs/unpairs-250.json",
		"./db/result/upair/diff/length1.csv",
		"./db/result/upair/diff/length_scale1.csv",
		"./db/result/upair/diff/direct1.csv",
		"./db/result/upair/diff/direct_scale1.csv")
}

func pairFactor() {
	statistic.PairFactorSta("./db/origin/detect/smark/london/midar-iff.nodes.cleaned.ip.londonmachine.smarkres.0", "./db/pairs/all_pairs-host-port-whois.json", "./db/result/pair/factor/london.csv", "./db/result/pair/factor/scatter.csv")
	statistic.PairFactorSta("./db/origin/detect/smark/london/midar-iff.nodes.cleaned.ip.londonmachine.smarkres.0", "./db/pairs/unpairs-250.json", "./db/result/upair/factor/london.csv", "./db/result/upair/factor/scatter.csv")
}

func rttGapSta() {
	statistic.RttGapSta("./db/origin/detect/smark/london/midar-iff.nodes.cleaned.ip.londonmachine.smarkres.0", "./db/pairs/all_pairs-host-port-whois.json", "./db/result/pair/rtt/london.csv")
	statistic.RttGapSta("./db/origin/detect/smark/london/midar-iff.nodes.cleaned.ip.londonmachine.smarkres.0", "./db/pairs/unpairs-250.json", "./db/result/upair/rtt/london.csv")
}

func ttlGapSta() {
	statistic.TTLScale("./db/origin/detect/sping/london/resultmidar-iff.nodes.cleaned.ip.londonmachine.spingres.0", "./db/pairs/all_pairs-host-port-whois.json", "./db/result/pair/ttl/london.csv")
	statistic.TTLScale("./db/origin/detect/sping/london/resultmidar-iff.nodes.cleaned.ip.londonmachine.spingres.0", "./db/pairs/unpairs-250.json", "./db/result/upair/ttl/london.csv")
}

func domainConsistSta() {
	statistic.DomainConsistenceSta("./db/origin/detect/itdk-run-20220224-dns-names.txt", "./db/pairs/all_pairs-host-port-whois.json", "./db/result/pair/domain/consist.txt")
	statistic.DomainConsistenceSta("./db/origin/detect/itdk-run-20220224-dns-names.txt", "./db/pairs/unpairs-250.json", "./db/result/upair/domain/consist.txt")
}

func domainDistance() {
	statistic.DomainDistanceSta("./db/origin/detect/itdk-run-20220224-dns-names.txt", "./db/pairs/all_pairs-host-port-whois.json", "./db/result/pair/domain/distance.txt")
	statistic.DomainDistanceSta("./db/origin/detect/itdk-run-20220224-dns-names.txt", "./db/pairs/unpairs-250.json", "./db/result/upair/domain/distance.txt")
}

func netSec() {
	statistic.NetSecSta("./db/pairs/all_pairs-host-port-whois.json", "./db/result/pair/netsec/netsec.csv")
	statistic.NetSecSta("./db/pairs/unpairs-250.json", "./db/result/upair/netsec/netsec.csv")
}

// use this top function to search best parameters from file which records trainning result
func findBest() {
	rfi, err := os.Open("./db/database/oob/1/dct-args.json")
	if err != nil {
		panic(err.Error())
	}
	defer rfi.Close()
	scaner := bufio.NewScanner(rfi)
	args := lib.MArgs{}
	for scaner.Scan() {
		t := scaner.Text()
		scaner.Scan()
		arg := lib.MArg{}
		err := json.Unmarshal([]byte(t), &arg)
		if err != nil {
			continue
		}
		args = append(args, arg)
	}
	sort.Sort(args)
	t := args[len(args)-1]
	fmt.Println(t)

}

// use this functions to sort parameters according specified field from file which records trainning result
func argsSort() {
	args, err := store.LoadAny[lib.MArg]("./db/machinelearn/args.json")
	if err != nil {
		panic(err)
	}
	var Margs lib.MArgs
	Margs = args
	sort.Sort(Margs)
	wfi, err := os.OpenFile("./db/machinelearn/args_sorted.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	for i, j := 0, len(Margs)-1; i < j; i, j = i+1, j-1 {
		Margs[i], Margs[j] = Margs[j], Margs[i]
	}
	for _, arg := range Margs {
		bs, _ := json.Marshal(arg)
		wtr.Write(bs)
		wtr.WriteString("\n")
	}
	wtr.Flush()
}

// make alias pairs from sublist of origin ip list
func oobPairs() {
	paris, err := store.LoadAny[lib.Pair]("./db/pairs/all_pairs-host-port-whois.json")
	rfi, err := os.Open("./db/database/oob/1/test.txt")
	if err != nil {
		panic(err.Error())
	}
	defer rfi.Close()
	parisMap := map[string]lib.Pair{}
	for _, p := range paris {
		parisMap[p.ID()] = p
	}
	scaner := bufio.NewScanner(rfi)
	ips := map[string]struct{}{}
	for scaner.Scan() {
		ips[scaner.Text()] = struct{}{}
	}
	wfi, err := os.OpenFile("./db/database/oob/1/test_paris.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	defer wtr.Flush()
	for _, pair := range paris {
		if _, ok := ips[pair.IPA]; ok {
			if _, ok := ips[pair.IPB]; ok {
				bs, _ := json.Marshal(pair)
				wtr.Write(bs)
				wtr.WriteString("\n")
			}
		}
	}
}

// running this function to generate charactor vectors block by block and send them to be predicted
// then we pick time-cost data from result
func rfc_pipe() {
	lwfi, err := os.OpenFile("./db/database/log/shanghai_200_clissifer_pipeline.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	//for _, num := range Num_w[:] {
	//	fmt.Printf("num: %v\n", num)
	start := time.Now().Unix()
	ipPath := "./db/database/city/shanghai/200.ipv4"
	allIP := utils.ReadALLLine(ipPath)
	tracepath := "./db/database/city/shanghai/200.trace"
	spingPath := "./db/database/city/shanghai/200.sping_by_qingdao"
	domainPath := "./db/database/city/shanghai/200_rdns"
	fmt.Println("cut slide")
	groups, _ := utils.Slide_windows(spingPath, allIP)
	tCnt := 0
	wfi, err := os.OpenFile("./db/database/city/shanghai/aliasC/alias.pair", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	wtr := bufio.NewWriter(wfi)
	for gi := 0; gi < len(groups); gi++ {
		group := groups[gi]
		fmt.Printf("group num;%v\n", gi)
		tCnt += time_cost.RfcPipeline(tracepath, spingPath, domainPath, group, wtr)
		lwfi.WriteString(fmt.Sprintf("group %v done time, %v\n", gi, time.Now().Unix()-start))
	}
	wtr.Flush()
	wfi.Close()
	//tCnt += UnCheck(un)
	end := time.Now().Unix()
	lwfi.WriteString(fmt.Sprintf("city %v,time(second) %v ,pairs cnt = %v\n", "shanghai", end-start, tCnt))
	//}
	lwfi.Close()
}

// use this top function to make charactor vectors for NB machine learning
func makeNB() {
	for i := 0; i < 5; i++ {
		rdpath := "./db/database/oob/" + strconv.Itoa(i+1) + "/test_sk.csv"
		wtpath := "./db/database/oob/" + strconv.Itoa(i+1) + "/alias/nb_test_sk.csv"
		data4sk.MakeNBNLBData(rdpath, wtpath)
	}
}
