/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-21:16
 * @Description: :)
 * @email: fyfishie@outlook.com
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
// by throwing pairs without same on/off line status away
func HostFilter(pairPath, spingPath, wtPath string) {
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	spingRes, err := store.LoadAny[lib.Sping](spingPath)
	esyerr.AutoPanic(err)
	spingMap := map[string]lib.Sping{}
	for _, s := range spingRes {
		spingMap[s.IP] = s
	}
	wfi, err := os.OpenFile(wtPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	filted := 0
	for _, pair := range pairs {
		onlineNum := 0
		if _, ok := spingMap[pair.IPA]; ok {
			onlineNum++
		}
		if _, ok := spingMap[pair.IPB]; ok {
			onlineNum++
		}
		switch onlineNum {
		case 0:
			bs, _ := json.Marshal(pair)
			wtr.Write(bs)
			wtr.WriteString("\n")
		case 2:
			bs, _ := json.Marshal(pair)
			wtr.Write(bs)
			wtr.WriteString("\n")
		case 1:
			filted++
		}
	}
	fmt.Printf("filted: %v\n", filted)
	esyerr.AutoPanic(wtr.Flush())
}
