/*
 * @Author: fyfishie
 * @Date: 2023-04-28:15
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-22:16
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package statistic

import (
	"alias_article/lib"
	"alias_article/utils"
	"fmt"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// statistic 'Difference Value of Relative round-trip time' charactor
func RttGapSta(tracePath string, pairPath string, wtPath string) {
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	traces, err := utils.LoadValidTrace(tracePath)
	esyerr.AutoPanic(err)
	ipTraceMap := map[string]lib.RawTrace{}
	for _, t := range traces {
		ipTraceMap[t.Ip] = t
	}
	validCnt := 0
	rttGapMap := map[int]int{}
	for _, pair := range pairs {
		if _, ok := ipTraceMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := ipTraceMap[pair.IPB]; !ok {
			continue
		}
		if !utils.ValidTracePair(ipTraceMap[pair.IPA], ipTraceMap[pair.IPB]) {
			continue
		}
		validCnt++
		rttGapMap[rttGap(ipTraceMap[pair.IPA], ipTraceMap[pair.IPB])]++
	}
	WriteAccuDataFrom(rttGapMap, wtPath, fmt.Sprintf("valid scale: %v\n", validCnt*100/len(pairs)))
}

func rttGap(traceA, traceB lib.RawTrace) int {
	traceA, traceB = utils.CutOffBeforeEndSame(traceA, traceB)
	rttA := (traceA.Results[len(traceA.Results)-1].Rtt - traceA.Results[0].Rtt) / 1000000
	rttB := (traceB.Results[len(traceB.Results)-1].Rtt - traceB.Results[0].Rtt) / 1000000
	return utils.Abs(rttB - rttA)
}
