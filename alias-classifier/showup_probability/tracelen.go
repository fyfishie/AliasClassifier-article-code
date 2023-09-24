/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-17:16
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

// calculates scale of ip who has route trace length data in origin ip list
func TraceLen(path string) {
	traces, err := store.LoadAny[lib.RawTrace](path)
	esyerr.AutoPanic(err)
	m := map[int]int{}
	total := 0
	for _, trace := range traces {
		if len(trace.Results) == 0 {
			continue
		}
		if trace.Ip == trace.Results[len(trace.Results)-1].Ip {
			total++
			m[len(trace.Results)]++
		}
	}
	fmt.Printf("total: %v\n", total)
	fmt.Printf("len(traces): %v\n", len(traces))
	fmt.Println(m)
}
