/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-27:08
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package statistic

import (
	"alias_article/lib"
	"alias_article/utils"
	"fmt"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// calculates 'Path Similarity Coefficient' charactor
func PairFactorSta(tracePath, pairPath string, wtpath string, scatterPath string) {
	traces, err := utils.LoadValidTrace(tracePath)
	esyerr.AutoPanic(err)
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	ipTraceMap := map[string]lib.RawTrace{}
	for _, trace := range traces {
		ipTraceMap[trace.Ip] = trace
	}
	factors := map[float64]int{}
	valid := 0
	for _, pair := range pairs {
		if !utils.ValidTracePair(ipTraceMap[pair.IPA], ipTraceMap[pair.IPB]) {
			continue
		}
		valid++
		fac := factor(ipTraceMap[pair.IPA], ipTraceMap[pair.IPB])
		factors[(fac)]++
	}
	WriteAccuDataFromFloat(factors, wtpath, fmt.Sprintf("same factor accu, valid scale: %v\n", (valid*100)/len(pairs)), 0.1)
	utils.WriteBar(factors, scatterPath, 0.1)
}

func factor(traceA, traceB lib.RawTrace) float64 {
	// straceA, straceB := utils.CutOffBeforeEndSame(traceA, traceB)
	// distance := utils.Distance(traceA, traceB)
	s := sameNum(traceA, traceB)
	totalLen := len(traceA.Results) + len(traceB.Results)
	factor := float64(s) / (float64(totalLen) / 2)
	if factor > 1 {
		fmt.Println("?")
	}
	return factor
}

func sameNum(traceA, traceB lib.RawTrace) int {
	ipMap := map[string]struct{}{}
	for _, hop := range traceA.Results {
		ipMap[hop.Ip] = struct{}{}
	}
	res := 0
	for _, hop := range traceB.Results {
		if _, ok := ipMap[hop.Ip]; ok {
			res++
		}
	}
	return res
}
