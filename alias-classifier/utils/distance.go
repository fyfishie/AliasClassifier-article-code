package utils

import (
	"alias_article/lib"

	"github.com/agnivade/levenshtein"
)

// calculates edit distance of two route trace
func Distance(traceA, traceB lib.RawTrace) int {
	allChar := map[string]bool{}
	for k, _ := range AllChar {
		allChar[k] = false
	}
	ip2AsciiMap := map[string]string{}
	strA := ""
	strB := ""
	for _, hop := range traceA.Results {
		if c, ok := ip2AsciiMap[hop.Ip]; ok {
			strA = strA + c
			continue
		}

		for k, v := range allChar {
			if !v {
				ip2AsciiMap[hop.Ip] = k
				strA = strA + k
				allChar[k] = true
				break
			}
		}
	}
	for _, hop := range traceB.Results {
		if c, ok := ip2AsciiMap[hop.Ip]; ok {
			strB = strB + c
			continue
		}
		for k, v := range allChar {
			if !v {
				ip2AsciiMap[hop.Ip] = k
				strB = strB + k
				allChar[k] = true
				break
			}
		}
	}
	return levenshtein.ComputeDistance(strA, strB)
}

var AllChar = map[string]bool{
	"0": false,
	"1": false,
	"2": false,
	"3": false,
	"4": false,
	"5": false,
	"6": false,
	"7": false,
	"8": false,
	"9": false,
	"a": false,
	"b": false,
	"c": false,
	"d": false,
	"e": false,
	"f": false,
	"g": false,
	"h": false,
	"i": false,
	"j": false,
	"k": false,
	"l": false,
	"m": false,
	"n": false,
	"o": false,
	"p": false,
	"q": false,
	"r": false,
	"s": false,
	"t": false,
	"u": false,
	"v": false,
	"w": false,
	"x": false,
	"y": false,
	"z": false,
}
