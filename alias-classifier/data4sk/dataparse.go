/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-27:09
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package data4sk

import (
	"alias_article/lib"
	"alias_article/utils"
	"math"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/fyfishie/ipop"
)

// calculates 'Difference Value of Path Length' of one pair
func difLength(traceA, traceB lib.RawTrace) (abs int) {
	traceA, traceB = utils.CutOffBeforeEndSame(traceA, traceB)
	tA := traceA.Results[len(traceA.Results)-1].TTL - traceA.Results[0].TTL
	tB := traceB.Results[len(traceB.Results)-1].TTL - traceB.Results[0].TTL
	abs = utils.Abs(tA - tB)
	if abs >= 4 {
		abs = 4
	}
	return abs
}

// calculates 'Difference Value of Path Direction' of one pair
func difDirect(traceA, traceB lib.RawTrace) (abs int) {
	sTraceA, sTraceB := utils.CutOffBeforeEndSame(traceA, traceB)
	dis := utils.Distance(sTraceA, sTraceB)
	//maxLen := utils.Max(len(traceA.Results), len(traceB.Results))
	g := dis / 2
	if g >= 8 {
		g = 8
	}
	return g
}

// calculates 'Path Similarity Coefficient' of one pair
func factor(traceA, traceB lib.RawTrace) int {
	s := sameNum(traceA, traceB)
	totalLen := len(traceA.Results) + len(traceB.Results)
	factor := float64(s) / (float64(totalLen) / 2) * 10
	if factor >= 6 {
		factor = 6
	}
	return int(factor)
}

// calculates 'Difference Value of Relative round-trip time' of one pair
func rttGap(traceA, traceB lib.RawTrace) int {
	traceA, traceB = utils.CutOffBeforeEndSame(traceA, traceB)
	rttA := (traceA.Results[len(traceA.Results)-1].Rtt - traceA.Results[0].Rtt) / 1000000
	rttB := (traceB.Results[len(traceB.Results)-1].Rtt - traceB.Results[0].Rtt) / 1000000
	g := utils.Abs(rttB - rttA)
	if g >= 180 {
		return 180
	}
	return g
}

// calculates 'Top-Level Domain Name Consistency' of one pair
func topConsist(a, b string) bool {
	ssA := strings.Split(a, ".")
	ssB := strings.Split(b, ".")
	return ssA[len(ssA)-1] == ssB[len(ssB)-1]
}

// calculates 'Second-Level Domain Name Consistency' of one pair
func secondConsist(a, b string) (valid bool, consist bool) {
	ssA := strings.Split(a, ".")
	ssB := strings.Split(b, ".")
	if len(ssA) < 2 && len(ssB) < 2 {
		return false, false
	}
	if len(ssA) < 2 || len(ssB) < 2 {
		return true, false
	}
	return true, ssA[len(ssA)-2] == ssB[len(ssB)-2]
}

// calculates 'Third-Level Domain Name Consistency' of one pair
func subConsist(a, b string) (valid, consist bool) {
	ssA := strings.Split(a, ".")
	ssB := strings.Split(b, ".")
	if len(ssA) < 3 && len(ssB) < 3 {
		return false, false
	}
	if len(ssA) < 3 || len(ssB) < 3 {
		return true, false
	}
	return true, ssA[len(ssA)-3] == ssB[len(ssB)-3]
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

// calculates 'Character Edit Distance of the Domain Name' of one pair
func domainDistance(da, db string) int {
	da = reverseDomain(da)
	db = reverseDomain(db)
	d := levenshtein.ComputeDistance(da, db)
	if d >= 33 {
		return 33
	}
	if d <= 0 {
		return 1
	}
	return d
}
func reverseDomain(s string) string {
	ss := strings.Split(s, ".")
	res := ""
	for i := len(ss) - 1; i > -1; i-- {
		res = res + ss[i]
	}
	return res
}

// calculates 'Spatial Distance of the IP pair' of one pair
func SecGap(a, b string) int {
	ia := ipop.String2Int(a)
	ib := ipop.String2Int(b)
	gap := utils.Abs(ia - ib)
	lg := math.Log2(float64(gap))
	g := int(math.Round(lg))
	if g <= 3 {
		return 3
	}
	if g >= 30 {
		return 30
	}
	return g
}
