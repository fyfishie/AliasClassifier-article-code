package data4sk

import (
	"alias_article/utils"
	"bufio"
	"os"
	"strconv"
	"strings"
)

// len_diff,dir_diff,same_factor,rtt_gap,reply_ttl_gap,top_domain_consist,second_domain_consist,sub_domain_consist,domain_distance,ip_distance,alias
var edge = []int{1, 3, 0, 25, 3, 0, 0, 0, 16, 24}

// make charactor vectors for nb machine learning
func MakeNBData(rdPath, wtPath string) {
	rfi, err := os.OpenFile(rdPath, os.O_RDONLY, 0000)
	if err != nil {
		panic(err)
	}
	scaner := bufio.NewScanner(rfi)
	scaner.Scan()
	head := scaner.Text()
	wfi, err := os.OpenFile(wtPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	wfi.WriteString(head + "\n")
	for scaner.Scan() {
		t := scaner.Text()
		ss := strings.Split(t, ",")
		nt := ""
		for _, s := range ss[:len(ss)-1] {
			it := utils.MustInt(s)
			ns := strconv.Itoa(it + 1)
			nt = nt + ns + ","
		}
		nt = nt + ss[len(ss)-1]
		wfi.WriteString(nt + "\n")
	}
	rfi.Close()
	wfi.Close()
}
func MakeNBNLBData(rdPath, wtPath string) {
	rfi, err := os.OpenFile(rdPath, os.O_RDONLY, 0000)
	if err != nil {
		panic(err)
	}
	scaner := bufio.NewScanner(rfi)
	scaner.Scan()
	head := scaner.Text()
	wfi, err := os.OpenFile(wtPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	wfi.WriteString(head + "\n")
	for scaner.Scan() {
		t := scaner.Text()
		ss := strings.Split(t, ",")
		if !count0(ss) {
			continue
		}
		nt := ""
		for index, s := range ss[:len(ss)-1] {
			it := utils.MustInt(s)
			if it == -1 {
				nt += "3,"
			} else {
				if it <= edge[index] {
					nt += "0,"
				} else {
					nt += "1,"
				}
			}
		}
		nt = nt + ss[len(ss)-1]
		wfi.WriteString(nt + "\n")
	}
	rfi.Close()
	wfi.Close()
}

func count0(ss []string) bool {
	l := len(ss)
	cnt := 0
	for _, s := range ss {
		if s == "-1" {
			cnt++
			if cnt*2 > l {
				return false
			}
		}
	}
	return true
}
