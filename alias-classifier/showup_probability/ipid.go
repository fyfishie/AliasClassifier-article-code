/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-28:11
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

// calculates scale of ip who has ipid data in origin ip list
func IPID(ipidPath, pairPath string) {
	ipids, err := store.LoadAny[lib.IPAID](ipidPath)
	esyerr.AutoPanic(err)
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	ipidHas := map[string]struct{}{}
	ipMap := map[string]struct{}{}
	paidNum := 0
	for _, pair := range pairs {
		ipMap[pair.IPA] = struct{}{}
		ipMap[pair.IPB] = struct{}{}
	}
	for _, ipid := range ipids {
		ipidHas[ipid.IP] = struct{}{}
	}
	for _, pair := range pairs {
		if _, ok := ipidHas[pair.IPA]; ok {
			if _, ok := ipidHas[pair.IPB]; ok {
				paidNum++
			}
		}
	}
	fmt.Printf("pair scale:%v\n", paidNum*100/len(pairs))
	fmt.Printf("ip scale:%v\n", len(ipidHas)*100/len(ipMap))
}
