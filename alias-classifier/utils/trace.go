/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-25:14
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package utils

import (
	"alias_article/lib"

	"github.com/fyfishie/dorapock/store"
)

func CutOffBeforeEndSame(traceA, traceB lib.RawTrace) (a, b lib.RawTrace) {
	l := LastSameIP(traceA, traceB)
	if l == "" {
		return traceA, traceB
	}
	sa := lib.RawTrace{
		Ip:      traceA.Ip,
		Results: []lib.Hop{},
	}
	sb := lib.RawTrace{
		Ip:      traceB.Ip,
		Results: []lib.Hop{},
	}
	sw := false
	for _, hop := range traceA.Results {
		if hop.Ip == l {
			if sw {
				sa.Results = []lib.Hop{}
			}
			sw = true
		}
		if sw {
			sa.Results = append(sa.Results, hop)
		}
	}
	sw = false
	for _, hop := range traceB.Results {
		if hop.Ip == l {
			if sw {
				sb.Results = []lib.Hop{}
			}
			sw = true
		}
		if sw {
			sb.Results = append(sb.Results, hop)
		}
	}
	return sa, sb
}

func LastSameIP(traceA, traceB lib.RawTrace) string {
	mapA := map[string]struct{}{}
	for _, hop := range traceA.Results {
		mapA[hop.Ip] = struct{}{}
	}
	for i := len(traceB.Results) - 1; i > -1; i-- {
		if _, ok := mapA[traceB.Results[i].Ip]; ok {
			return traceB.Results[i].Ip
		}
	}
	return ""
}

func BranchTailNum(traceA, traceB lib.RawTrace) int {
	same := LastSameIP(traceA, traceB)
	res := 0
	for i := len(traceA.Results) - 1; i >= 0; i-- {
		if traceA.Results[i].Ip == same {
			break
		}
		res++
	}
	for i := len(traceB.Results) - 1; i >= 0; i-- {
		if traceB.Results[i].Ip == same {
			break
		}
		res++
	}
	return res
}

func LoadValidTrace(path string) ([]lib.RawTrace, error) {
	traces, err := store.LoadAny[lib.RawTrace](path)
	if err != nil {
		return nil, err
	}
	res := []lib.RawTrace{}
	for _, trace := range traces {
		if len(trace.Results) != 0 {
			if trace.Ip == trace.Results[len(trace.Results)-1].Ip {
				if !circleTrace(trace) {
					res = append(res, trace)
				}
			}
		}
	}
	return res, nil
}
func circleTrace(trace lib.RawTrace) bool {
	m := map[string]struct{}{}
	for _, hop := range trace.Results {
		if _, ok := m[hop.Ip]; ok {
			return true
		}
		m[hop.Ip] = struct{}{}
	}
	return false
}
