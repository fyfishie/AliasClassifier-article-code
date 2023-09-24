/*
 * @Author: fyfishie
 * @Date: 2023-04-26:16
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-28:17
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

// calculates 'Difference Value of Path Length'
// and 'Difference Value of Path Direction'
// then write them out
func PairDiffSta(smarkPath, pairsPath, lenDifPath, lenDifScalePath, dirDifPath, dirDifScalePath string) {
	traces, err := utils.LoadValidTrace(smarkPath)
	esyerr.AutoPanic(err)
	pairs, err := store.LoadAny[lib.Pair](pairsPath)
	esyerr.AutoPanic(err)
	ipTraceMap := map[string]lib.RawTrace{}
	for _, trace := range traces {
		ipTraceMap[trace.Ip] = trace
	}
	difLengthMap := map[int]int{}
	difDirectMap := map[int]int{}
	difLengthScaleMap := map[int]int{}
	difDirectScaleMap := map[int]int{}
	valid := 0
	total := 0
	for _, pair := range pairs {
		total++
		if _, ok := ipTraceMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := ipTraceMap[pair.IPB]; !ok {
			continue
		}
		if !utils.ValidTracePair(ipTraceMap[pair.IPA], ipTraceMap[pair.IPB]) {
			continue
		}
		valid++
		lengthAbs, lengthScale := difLength(ipTraceMap[pair.IPA], ipTraceMap[pair.IPB])
		directAbs, directScale := difDirect(ipTraceMap[pair.IPA], ipTraceMap[pair.IPB])
		difLengthMap[lengthAbs]++
		difDirectMap[directAbs]++
		difLengthScaleMap[lengthScale]++
		difDirectScaleMap[directScale]++
	}
	WriteAccuDataFrom(difLengthMap, lenDifPath, fmt.Sprintf("valid pair scale: %v\n", valid*100/total))
	WriteAccuDataFrom(difLengthScaleMap, lenDifScalePath, fmt.Sprintf("valid pair scale: %v\n", valid*100/total))
	WriteAccuDataFrom(difDirectMap, dirDifPath, fmt.Sprintf("valid pair scale: %v\n", valid*100/total))
	WriteAccuDataFrom(difDirectScaleMap, dirDifScalePath, fmt.Sprintf("valid pair scale: %v\n", valid*100/total))
}

func difLength(traceA, traceB lib.RawTrace) (abs int, scale int) {
	traceA, traceB = utils.CutOffBeforeEndSame(traceA, traceB)
	tA := traceA.Results[len(traceA.Results)-1].TTL - traceA.Results[0].TTL
	tB := traceB.Results[len(traceB.Results)-1].TTL - traceB.Results[0].TTL
	abs = utils.Abs(tA - tB)
	return abs, -1
}

func difDirect(traceA, traceB lib.RawTrace) (abs int, scale int) {
	sTraceA, sTraceB := utils.CutOffBeforeEndSame(traceA, traceB)
	dis := utils.Distance(sTraceA, sTraceB)
	maxLen := utils.Max(len(traceA.Results), len(traceB.Results))
	return dis, dis * 100 / maxLen
}

func AllFit(smarkPath, pairsPath string) {
	traces, err := utils.LoadValidTrace(smarkPath)
	esyerr.AutoPanic(err)
	pairs, err := store.LoadAny[lib.Pair](pairsPath)
	esyerr.AutoPanic(err)
	ipTraceMap := map[string]lib.RawTrace{}
	for _, trace := range traces {
		ipTraceMap[trace.Ip] = trace
	}
	fit := 0
	total := 0
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
		total++
		lengthAbs, _ := difLength(ipTraceMap[pair.IPA], ipTraceMap[pair.IPB])
		directAbs, _ := difDirect(ipTraceMap[pair.IPA], ipTraceMap[pair.IPB])
		if lengthAbs < 3 {
			if directAbs <= 2 {
				fit++
			}
		}
	}
	fmt.Printf("fit scale:%v\n", fit*100/total)
}
