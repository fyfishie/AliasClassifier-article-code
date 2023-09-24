/*
 * @Author: fyfishie
 * @Date: 2023-05-30:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-21:16
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package filter

import (
	"alias_article/lib"
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// used to generate Ground-real datasets:
// by throwing pairs without same ports listening away
func PortsFilt(portPath, pairPath, wtPath string) {
	portsMap := loadPorts(portPath)
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	wfi, err := os.OpenFile(wtPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	filted := 0
	for _, pair := range pairs {
		if _, ok := portsMap[pair.IPA]; ok {
			if _, ok := portsMap[pair.IPB]; ok {
				portsA := portsMap[pair.IPA]
				portsB := portsMap[pair.IPB]
				if portsEqual(portsA, portsB) {
					bs, _ := json.Marshal(pair)
					wtr.Write(bs)
					wtr.WriteString("\n")
					continue
				} else {
					filted++
				}
			}
		}
		bs, _ := json.Marshal(pair)
		wtr.Write(bs)
		wtr.WriteString("\n")
	}
	fmt.Printf("filted: %v\n", filted)
	esyerr.AutoPanic(wtr.Flush())
}

func portsEqual(a, b []int) bool {
	am := map[int]struct{}{}
	bm := map[int]struct{}{}
	for _, ai := range a {
		am[ai] = struct{}{}
	}
	for _, bi := range b {
		bm[bi] = struct{}{}
	}
	if len(am) != len(bm) {
		return false
	}
	for k := range am {
		if _, ok := bm[k]; !ok {
			return false
		}
	}
	return true
}

func loadPorts(portPath string) map[string][]int {
	ports, err := store.LoadAny[lib.PortScanItem](portPath)
	esyerr.AutoPanic(err)
	res := map[string][]int{}
	for _, portScanItem := range ports {
		if len(portScanItem.Ports) != 0 {
			for _, port := range portScanItem.Ports {
				if _, ok := res[portScanItem.IP]; !ok {
					res[portScanItem.IP] = []int{}
				}
				res[portScanItem.IP] = append(res[portScanItem.IP], port.Port)
			}
		}
	}
	return res
}
