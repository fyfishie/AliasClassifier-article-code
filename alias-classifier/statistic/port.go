/*
 * @Author: fyfishie
 * @Date: 2023-04-28:18
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-01:17
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package statistic

import (
	"alias_article/lib"
	"bufio"
	"github.com/fyfishie/dorapock/store"
	"os"
	"strconv"

	"github.com/fyfishie/esyerr"
)

// statistic port similarity charactor
func PortSimilaritySta(scanRdPath string, nodePath string, wtpath string) {
	nodes, _ := store.LoadAny[lib.Node](nodePath)
	scanItems, _ := store.LoadAny[lib.PortScanItem](scanRdPath)
	ipNodeMap := map[string]*lib.Node{}
	for i := 0; i < len(nodes); i++ {
		for _, ip := range nodes[i].IPList {
			ipNodeMap[ip] = &(nodes[i])
		}
	}
	ipPortsMap := map[string][]int{}
	for _, item := range scanItems {
		if _, ok := ipPortsMap[item.IP]; !ok {
			ipPortsMap[item.IP] = []int{}
		}
		for _, p := range item.Ports {
			l := ipPortsMap[item.IP]
			l = append(l, p.Port)
			ipPortsMap[item.IP] = l
		}
	}
	for ip, ports := range ipPortsMap {
		ipNodeMap[ip].OpenPorts = append(ipNodeMap[ip].OpenPorts, ports)
	}
	for i := 0; i < len(nodes); i++ {
		nodes[i].PortAndNum, nodes[i].PortOrNum, nodes[i].PortSimilarity = portSimilarity(nodes[i].OpenPorts)
	}
	writePortsRes(nodes, wtpath)
}

func portSimilarity(lists [][]int) (andNum, orNum, similarity int) {
	andMap := map[int]struct{}{}
	orMap := map[int]struct{}{}
	if len(lists) == 0 {
		return -1, -1, -1
	}
	for _, item := range lists[0] {
		andMap[item] = struct{}{}
		orMap[item] = struct{}{}
	}
	for _, list := range lists[1:] {
		for _, item := range list {
			if _, ok := andMap[item]; !ok {
				delete(andMap, item)
			}
			orMap[item] = struct{}{}
		}
	}
	andSet := []int{}
	for k, _ := range andMap {
		andSet = append(andSet, k)
	}
	return len(andMap), len(orMap), len(andMap) * 100 / len(orMap)
}

func writePortsRes(nodes []lib.Node, wtpath string) {
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	esyerr.AutoPanic(err)
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	wtr.WriteString("node_id,and_num,or_num,port_similarity\n")
	for _, node := range nodes {
		// if node.PortSimilarity == -1 {
		// 	continue
		// }
		// if node.PortOrNum < 3 {
		// 	continue
		// }
		wtr.WriteString(node.NodeID + "," + strconv.Itoa(node.PortAndNum) + "," + strconv.Itoa(node.PortOrNum) + "," + strconv.Itoa(node.PortSimilarity) + "\n")
	}
	esyerr.AutoPanic(wtr.Flush())
}
