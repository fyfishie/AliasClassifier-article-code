/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-28:12
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package showup_probability

import (
	"alias_article/lib"
	"alias_article/utils"
	"fmt"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// calculates scale of ip who has rtt data in origin ip list
func RttScale(pairPath string, smarkPath string) {
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	traces, err := utils.LoadValidTrace(smarkPath)
	esyerr.AutoPanic(err)
	ipMap := map[string]struct{}{}
	for _, pair := range pairs {
		ipMap[pair.IPA] = struct{}{}
		ipMap[pair.IPB] = struct{}{}
	}

	ipTraceMap := map[string]lib.RawTrace{}
	for _, trace := range traces {
		ipTraceMap[trace.Ip] = trace
	}

	pairNum := 0
	for _, pair := range pairs {
		if _, ok := ipTraceMap[pair.IPA]; ok {
			if _, ok := ipTraceMap[pair.IPB]; ok {
				pairNum++
			}
		}
	}

	ipNum := 0
	for ip := range ipMap {
		if _, ok := ipTraceMap[ip]; ok {
			ipNum++
		}
	}
	fmt.Printf("pair scale:%v\n", pairNum*100/len(pairs))
	fmt.Printf("ip scale:%v\n", ipNum*100/len(ipMap))
}
