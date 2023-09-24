/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-08-03:21
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package filter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"mlar/lib"
	"mlar/utils"
	"os"

	"github.com/fyfishie/dorapock/store"
)

type Filter struct {
	ispPath     string
	pairs       map[string]lib.Pair
	tracePaths  []string
	interfaceIP map[string]struct{}
}

// codes below implements filt function mentioned in mlar
func Filt(parisPath, ispPath string, pairsWtPath, interfaceIPWtpath string, spingsPath, tracesPath []string, runtineNum int, done chan struct{}) {
	f := Filter{ispPath: ispPath, pairs: make(map[string]lib.Pair), tracePaths: tracesPath, interfaceIP: map[string]struct{}{}}
	pairs, err := store.LoadAny[lib.Pair](parisPath)
	if err != nil {
		panic(err)
	}
	pm := map[string]lib.Pair{}
	for _, p := range pairs {
		pm[p.ID()] = p
	}
	f.pairs = pm
	fmt.Printf("%v isp\n", runtineNum)
	f.ispFilt()
	for i, p := range f.tracePaths {
		fmt.Printf("%v trace i: %v\n", runtineNum, i)
		f.traceEditFilt(p)
		f.traceFilt(p)
		f.traceLenFilt(p)
	}
	for i, p := range spingsPath {
		fmt.Printf("%v sping i: %v\n", runtineNum, i)
		f.spingFilt(p)
	}
	wfi, err := os.OpenFile(pairsWtPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	defer wtr.Flush()
	for _, p := range f.pairs {
		bs, _ := json.Marshal(p)
		wtr.Write(bs)
		wtr.WriteString("\n")
	}
	wfi2, err := os.OpenFile(interfaceIPWtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr2 := bufio.NewWriter(wfi2)
	defer wtr2.Flush()
	interfaceip := map[string]struct{}{}
	for _, p := range f.pairs {
		interfaceip[p.IPA] = struct{}{}
		interfaceip[p.IPB] = struct{}{}
	}
	for ip := range interfaceip {
		wtr2.WriteString(ip + "\n")
	}
	wtr.Flush()
	wtr2.Flush()
	done <- struct{}{}
}

// mlar use isp infomation as one accordance of checking alias or anti-alias pairs
func (f *Filter) ispFilt() {
	isps, err := store.LoadAny[lib.ISP](f.ispPath)
	if err != nil {
		panic(err)
	}
	ispMap := map[string]lib.ISP{}
	for _, isp := range isps {
		ispMap[isp.IP] = isp
	}
	newParis := map[string]lib.Pair{}
	for _, pair := range f.pairs {
		if _, ok := ispMap[pair.IPA]; ok {
			if _, ok := ispMap[pair.IPB]; ok {
				a := ispMap[pair.IPA].Isp
				b := ispMap[pair.IPB].Isp
				if a != b {
					continue
				}
			}
		}
		newParis[pair.ID()] = pair
	}
	f.pairs = newParis
	for _, pair := range newParis {
		f.interfaceIP[pair.IPA] = struct{}{}
		f.interfaceIP[pair.IPB] = struct{}{}
	}
}

// mlar thinks that if two ip belongs to the same router, they can not show up at the same route trace
func (f *Filter) traceFilt(tracePath string) {
	traces := utils.ValidTrace(tracePath)
	for _, trace := range traces {
		for _, hop := range trace.Results {
			if _, ok := f.interfaceIP[hop.Ip]; ok {
				for _, hop2 := range trace.Results {
					p := lib.Pair{IPA: hop.Ip, IPB: hop2.Ip}
					delete(f.pairs, p.ID())
				}
			}
		}
	}
	newInterfaceIP := map[string]struct{}{}
	for _, p := range f.pairs {
		newInterfaceIP[p.IPA] = struct{}{}
		newInterfaceIP[p.IPB] = struct{}{}
	}
	f.interfaceIP = newInterfaceIP
}

// mlar thinks that if edit distance of traces of two ip is greater than 2, they are not alias ip pair
func (f *Filter) traceEditFilt(tracePath string) {
	traces := utils.ValidTrace(tracePath)
	traceMap := map[string]lib.RawTrace{}
	for _, t := range traces {
		traceMap[t.Ip] = t
	}
	newPairs := map[string]lib.Pair{}
	for _, p := range f.pairs {
		if _, ok := traceMap[p.IPA]; ok {
			if _, ok := traceMap[p.IPB]; ok {
				if traceEdit(traceMap[p.IPA], traceMap[p.IPB]) > 2 {
					continue
				}
			}
		}
		newPairs[p.ID()] = p
	}
	f.pairs = newPairs
	interfaceIP := map[string]struct{}{}
	for _, p := range newPairs {
		interfaceIP[p.IPA] = struct{}{}
		interfaceIP[p.IPB] = struct{}{}
	}
}

// it calculates trace edit distance for two route trace
func traceEdit(traceA, traceB lib.RawTrace) int {
	hopIPMapA := map[int]string{}
	hopIPMapB := map[int]string{}
	allHop := map[int]struct{}{}
	for _, hop := range traceA.Results {
		hopIPMapA[hop.TTL] = hop.Ip
		allHop[hop.TTL] = struct{}{}
	}
	for _, hop := range traceB.Results {
		hopIPMapB[hop.TTL] = hop.Ip
		allHop[hop.TTL] = struct{}{}
	}
	cnt := 0
	for hop := range allHop {
		_, o1 := hopIPMapA[hop]
		_, o2 := hopIPMapB[hop]
		if o1 != o2 {
			cnt++
		}
	}
	return cnt
}

func (f *Filter) traceLenFilt(tracePath string) {
	traces := utils.ValidTrace(tracePath)
	newPairs := map[string]lib.Pair{}
	traceMap := map[string]lib.RawTrace{}
	for _, t := range traces {
		traceMap[t.Ip] = t
	}
	for _, pair := range f.pairs {
		if _, ok := traceMap[pair.ID()]; ok {
			if _, ok := traceMap[pair.ID()]; ok {
				if traceMap[pair.IPA].Results[len(traceMap[pair.ID()].Results)-1].TTL-traceMap[pair.IPB].Results[len(traceMap[pair.IPB].Results)-1].TTL > 3 {
					continue
				}
			}
		}
		newPairs[pair.ID()] = pair
	}
	f.pairs = newPairs
	newInterfaceIP := map[string]struct{}{}
	for _, pair := range newPairs {
		newInterfaceIP[pair.IPA] = struct{}{}
		newInterfaceIP[pair.IPB] = struct{}{}
	}
	f.interfaceIP = newInterfaceIP
}

// mlar thinks that if two ip are alias ip pair, they should be in tha same status (reachable or unreachable)
func (f *Filter) spingFilt(spingPath string) {
	spings, err := store.LoadAny[lib.Sping](spingPath)
	if err != nil {
		panic(err)
	}
	spingsMap := map[string]lib.Sping{}
	for _, s := range spings {
		spingsMap[s.IP] = s
	}
	newpairs := map[string]lib.Pair{}
	for _, pair := range f.pairs {
		_, o1 := spingsMap[pair.IPA]
		_, o2 := spingsMap[pair.IPB]
		if o1 == o2 {
			newpairs[pair.ID()] = pair
		}
	}
	f.pairs = newpairs
	newInter := map[string]struct{}{}
	for _, p := range f.pairs {
		newInter[p.IPA] = struct{}{}
		newInter[p.IPB] = struct{}{}
	}
	f.interfaceIP = newInter
}
