/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-08-03:20
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package character

import (
	"bufio"
	"math"
	"mlar/lib"
	"mlar/utils"
	"os"
)

var ipTracesMap = map[string][]int{}
var interfaceipMap = map[string]struct{}{}

// it is write specially for pipeline-generator, work before train data generated.
func RttSameFactor4PipelineInit(tracePaths []string, allInterPath string) {
	rfi, err := os.Open(allInterPath)
	if err != nil {
		panic(err.Error())
	}
	defer rfi.Close()
	scaner := bufio.NewScanner(rfi)
	for scaner.Scan() {
		ip := scaner.Text()
		interfaceipMap[ip] = struct{}{}
		ipTracesMap[ip] = []int{}
	}
	for _, p := range tracePaths {
		traces := utils.ValidTrace(p)
		tm := map[string]lib.RawTrace{}
		for _, t := range traces {
			tm[t.Ip] = t
		}
		for ip := range interfaceipMap {
			if t, ok := tm[ip]; ok {
				ipTracesMap[ip] = append(ipTracesMap[ip], (t.Results[len(t.Results)-1].Rtt / 1000000))
			} else {
				ipTracesMap[ip] = append(ipTracesMap[ip], -1)
			}
		}
	}
}

// it is write specially for pipeline-generator, making program faster and saving more resource
func RttSameFactor4City(pairs []lib.Pair) map[string]int {
	res := map[string]int{}
	for _, p := range pairs {
		t := rttfactor(ipTracesMap[p.IPA], ipTracesMap[p.IPB], 6)
		if t < 0 {
			t = -1
		}
		res[p.ID()] = t
	}
	return res
}

// it calculate rtt_same factor for one pair according to traces of the pair
func rttfactor(rttsA, rttsB []int, traceNum int) int {
	if len(rttsA) != traceNum || len(rttsB) != traceNum {
		return -1
	}
	total := 0.0
	for i := 0; i < traceNum; i++ {
		a := rttsA[i]
		b := rttsB[i]
		if a == -1 || b == -1 || a > 1000 || b > 1000 {
			continue
		}
		total += math.Pow(float64(a)-float64(b), 2)
	}
	t := math.Sqrt(total)
	it := int(math.Round(t))
	return it
}
