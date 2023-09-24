/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-16:10
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package showup_probability

import (
	"alias_article/lib"
	"fmt"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// calculates scale of ip who has whois data in origin ip list
func WhoisScale(pairPath, whoisPath string) {
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	whoisItems, err := store.LoadAny[lib.WhoisRes](whoisPath)
	esyerr.AutoPanic(err)
	fmt.Printf("len(whoisItems): %v\n", len(whoisItems))
	whoisMap := map[string]lib.WhoisRes{}
	for _, item := range whoisItems {
		whoisMap[item.IP] = item
	}
	no := 0
	all := 0
	mid := 0
	for _, pair := range pairs {
		_, ok1 := whoisMap[pair.IPA]
		_, ok2 := whoisMap[pair.IPB]
		if ok1 && ok2 {
			all++
		} else if (!ok1) && (!ok2) {
			no++
		} else {
			mid++
		}
	}
	fmt.Printf("no: %v\n", no)
	fmt.Printf("all: %v\n", all)
	fmt.Printf("mid: %v\n", mid)
}
