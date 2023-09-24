/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-08-04:20
 * @Description: :)
 * @email: fyfishie@outlook.com
 */

package character

import (
	"github.com/fyfishie/dorapock/store"
	"math"
	"mlar/lib"
	"mlar/utils"
)

var maps = []map[string]lib.RawTrace{}

// it is write specially for pipeline-generator, work before train data generated.
func TraceFactor4PipelineInit(tracespath []string) {
	for _, p := range tracespath {
		traceMap := map[string]lib.RawTrace{}
		traces, err := store.LoadAny[lib.RawTrace](p)
		if err != nil {
			panic(err)
		}
		for _, t := range traces {
			traceMap[t.Ip] = t
		}
		maps = append(maps, traceMap)
	}
}

// it is write specially for pipeline-generator, making program faster and saving more resource
func TraceFactor4City(pairs []lib.Pair) map[string]int {
	sumList := make([]float64, len(pairs))
	for i := 0; i < len(sumList); i++ {
		sumList[i] = 0
	}
	for _, m := range maps {
		for i, p := range pairs {
			f := traceFactor(m[p.IPA], m[p.IPB])
			sumList[i] += f
		}
	}
	l := len(maps)
	res := map[string]int{}

	for i, p := range pairs {
		res[p.ID()] = int(math.Round((sumList[i] / float64(l)) * 100))

	}
	for k, v := range res {
		if !(v < 0) {
			res[k] = -1
		}
	}
	return res
}

// calculate trace_same factor of all pairs input
func TraceFactor(tracepaths []string, pairpath string) map[string]int {
	pairs, err := store.LoadAny[lib.Pair](pairpath)
	if err != nil {
		panic(err)
	}
	sumList := make([]float64, len(pairs))
	for _, p := range tracepaths {
		traceMap := map[string]lib.RawTrace{}
		traces, err := store.LoadAny[lib.RawTrace](p)
		if err != nil {
			panic(err)
		}
		for _, t := range traces {
			traceMap[t.Ip] = t
		}
		for i, p := range pairs {
			sumList[i] += traceFactor(traceMap[p.IPA], traceMap[p.IPB])
		}
	}
	res := map[string]int{}
	for i, p := range pairs {
		res[p.ID()] = int(math.Round((sumList[i] / 10) * 100))
	}
	return res
}

// calculate trace_same factor of one pair
func traceFactor(traceA, traceB lib.RawTrace) float64 {
	maxL := len(traceB.Results)
	if len(traceA.Results) > maxL {
		maxL = len(traceA.Results)
	}
	f := 1.0 - float64(utils.DirDistance(traceA, traceB))/float64(maxL)
	return f
}
