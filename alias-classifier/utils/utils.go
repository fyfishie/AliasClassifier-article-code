/*
 * @Author: fyfishie
 * @Date: 2023-04-18:07
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-21:16
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package utils

import (
	"alias_article/lib"
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

func MustInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func PanicOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func Mode(list []int) int {
	m := map[int]int{}
	for _, i := range list {
		m[i]++
	}
	maxNum := 0
	maxCnt := 0
	for num, cnt := range m {
		if cnt > maxCnt {
			maxNum = num
			maxCnt = cnt
		}
	}
	return maxNum
}

func ModeCnt(list []int) int {
	m := map[int]int{}
	for _, i := range list {
		m[i]++
	}
	maxCnt := 0
	for _, cnt := range m {
		if cnt > maxCnt {
			maxCnt = cnt
		}
	}
	return maxCnt

}

func MaxGap(list []int) int {
	max := 0
	min := math.MaxInt
	for _, i := range list {
		if i < min {
			min = i
		}
		if i > max {
			max = i
		}
	}
	return max - min
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func ValidTracePair(traceA, traceB lib.RawTrace) bool {
	if len(traceA.Results) == 0 || len(traceB.Results) == 0 {
		return false
	}
	if traceA.Ip != traceA.Results[len(traceA.Results)-1].Ip {
		return false
	}
	return traceB.Ip == traceB.Results[len(traceB.Results)-1].Ip
}

func TrimIPFromPair(pairPath string) {
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	m := map[string]struct{}{}
	for _, p := range pairs {
		m[p.IPA] = struct{}{}
		m[p.IPB] = struct{}{}
	}
	fmt.Printf("len(m): %v\n", len(m))
}
func whoisInfoEqule(a, b lib.WhoIsInfoInstance) bool {
	e := true
	e = e && a.Domain.UpdatedDate == b.Domain.UpdatedDate
	e = e && a.Domain.CreatedDate == b.Domain.CreatedDate
	e = e && a.Registrant.Organization == b.Registrant.Organization
	e = e && a.Registrant.Name == b.Registrant.Name
	e = e && a.Registrant.Country == b.Registrant.Country
	e = e && a.Registrant.City == b.Registrant.City
	e = e && a.Registrant.Email == b.Registrant.Email
	return e
}

func ReverseDomain(s string) string {
	ss := strings.Split(s, ".")
	res := ""
	for i := len(ss) - 1; i > -1; i-- {
		res = res + ss[i]
	}
	return res
}

func WriteBar(data map[float64]int, path string, gap float64) {
	wfi, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	var start float64 = 0
	for start = 0; start <= 1; start += gap {
		high := 0
		for k, v := range data {
			if k >= start && k < start+gap {
				high += v
			}
		}
		wtr.WriteString(fmt.Sprintf("%.2f,%d\n", (start + gap/2), high))
	}
	wtr.Flush()
}
