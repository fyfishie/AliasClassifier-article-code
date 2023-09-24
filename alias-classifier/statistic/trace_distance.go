/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-25:11
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package statistic

import (
	"alias_article/lib"
	"alias_article/utils"
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/fyfishie/esyerr"
)

func TraceDistanceSta(tracePath, wtpath string) {
	normal, scale := TraceSameFactorMake(tracePath, wtpath)
	WriteAccuDataFrom(normal, wtpath, "distance accu\n")
	WriteAccuDataFromFloat(scale, wtpath, "distance scale sccu\n", 0.1)
}

// calculates 'trace distance of two route trace to the same target ip' charactor
func TraceSameFactorMake(traceDirPath string, wtpath string) (normal map[int]int, scale map[float64]int) {
	fs, err := os.ReadDir(traceDirPath)
	esyerr.AutoPanic(err)
	ipTracesMap := map[string][]lib.RawTrace{}
	cnt := 0
	for _, f := range fs {
		if cnt == 10 {
			break
		}
		cnt++
		traces, err := utils.LoadValidTrace(traceDirPath + "/" + f.Name())
		esyerr.AutoPanic(err)
		for _, trace := range traces {
			if _, ok := ipTracesMap[trace.Ip]; !ok {
				ipTracesMap[trace.Ip] = []lib.RawTrace{}
			}
			ipTracesMap[trace.Ip] = append(ipTracesMap[trace.Ip], trace)
		}
	}
	normalRes := map[int]int{}
	scaleRes := map[float64]int{}
	count := 0
	for _, traces := range ipTracesMap {
		if count == 10000 {
			break
		}
		count++
		if len(traces) < 2 {
			continue
		}
		normal := avgDistance(traces)
		normalRes[int(math.Round(normal))]++
		max := 0
		for _, trace := range traces {
			if len(trace.Results) > max {
				max = len(trace.Results)
			}
		}
		scale := float64(math.Round(normal)) / float64(max)
		if scale > 100 {
			fmt.Println("?")
		}
		scaleRes[scale]++
	}
	return normalRes, scaleRes
}
func avgDistance(traces []lib.RawTrace) float64 {
	total := 0
	for i := 0; i < len(traces); i++ {
		for j := i + 1; j < len(traces); j++ {
			traceA := traces[i]
			traceB := traces[j]
			dis := utils.Distance(traceA, traceB)
			total += dis
		}
	}
	length := len(traces)
	pairNum := length * (length - 1) / 2
	return float64(total) / float64(pairNum)
}
func WriteDistance(factors map[string]float64, wtpath string, head string) {
	staMap := map[int]int{}
	for _, factor := range factors {
		staMap[int(factor)]++
	}
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	wtr.WriteString(head)
	max := 0
	for left, _ := range staMap {
		if left > max {
			max = left
		}
	}
	cntList := []int{}
	for i := 0; i < max+1; i++ {
		if cnt, ok := staMap[i]; ok {
			cntList = append(cntList, cnt)
			continue
		}
		cntList = append(cntList, 0)
	}
	total := 0
	for i := 0; i < len(cntList); i++ {
		total += cntList[i]
		wtr.WriteString(strconv.Itoa(i) + "-" + strconv.Itoa(i+1) + "," + strconv.Itoa(total*100/len(factors)) + "\n")

	}
	esyerr.AutoPanic(wtr.Flush())
}
