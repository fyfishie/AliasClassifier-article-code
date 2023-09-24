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

// calculates scale of ip who has port data in origin ip list
func PortScale(portPath string, pairPath string) {
	ports := loadPorts(portPath)
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	pairNum := 0
	ipMap := map[string]struct{}{}
	for _, pair := range pairs {
		ipMap[pair.IPA] = struct{}{}
		ipMap[pair.IPB] = struct{}{}
		if _, ok := ports[pair.IPA]; ok {
			if _, ok := ports[pair.IPB]; ok {
				pairNum++
			}
		}
	}
	fmt.Printf("pair scale:%v\n", pairNum*10000/len(pairs))
	fmt.Printf("ip scale:%v\n", len(ports)*10000/len(ipMap))
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
