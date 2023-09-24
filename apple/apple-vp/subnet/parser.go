/*
 * @Author: fyfishie
 * @Date: 2023-03-21:08
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-14:21
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package subnet

import (
	"sort"
	"vp/lib"
	"vp/status"
)

type Parser struct {
	//raw data
	ipToDoList []int
	//send result out
	ResultOutChan chan lib.AntiDescriptor
	//if any parallel talent added in, this is ceil of parallel go routine numbers
	routineNum int
	//status monitor
	Eye status.StatusEye
}

func NewParser(ipToDoList []int, routineNum int) *Parser {
	p := Parser{
		ipToDoList:    ipToDoList,
		routineNum:    routineNum,
		ResultOutChan: make(chan lib.AntiDescriptor, 1000),
		Eye:           status.StatusEye{},
	}
	return &p
}

func (p *Parser) WithRun() *Parser {
	go func() {
		//always sort, some one (like me) will must forget it, must!!!
		sort.Ints(p.ipToDoList)
		antiAliasSets := [][]int{}
		netBlockInPutChan := make(chan []int, 10)
		gatherByDivisionOutChan := make(chan [][]int, 100)
		netBlocks := divideIpList(p.ipToDoList)
		g := newGather(netBlockInPutChan, gatherByDivisionOutChan)
		g.gatherRun(p.routineNum)
		go func() {
			for _, netBlock := range netBlocks {
				netBlockInPutChan <- netBlock
			}
		}()

		for i := 0; i < len(netBlocks); i++ {
			antiAliasSet := <-gatherByDivisionOutChan
			antiAliasSets = append(antiAliasSets, antiAliasSet...)
		}
		if len(antiAliasSets) > 0 {
			for _, antiSet := range antiAliasSets {
				antiMap := set2Map(antiSet)
				for ip, set := range antiMap {
					p.ResultOutChan <- lib.AntiDescriptor{Ip: ip, AntiAliasSet: set}
				}
			}
		}
		close(p.ResultOutChan)
	}()
	return p
}

func divideIpList(ipListAll []int) [][]int {
	res := [][]int{}
	segMax := 0
	segMin := 0
	segMin = (ipListAll[0] / 1024) * 1024
	if (ipListAll[len(ipListAll)-1]+1)%1024 == 0 {
		segMax = ipListAll[len(ipListAll)-1]
	} else {
		segMax = ((ipListAll[len(ipListAll)-1])/1024)*1024 + 1023
	}
	leftEdge := segMin
	rightEdge := segMin + 1023
	ipList := []int{leftEdge, rightEdge}
	for i := 0; rightEdge <= segMax; i++ {
		if i == len(ipListAll) {
			res = append(res, ipList)
			break
		}
		if ipListAll[i] >= leftEdge && ipListAll[i] <= rightEdge {
			ipList = append(ipList, ipListAll[i])
		} else {
			if len(ipList) > 2 {
				res = append(res, ipList)
			}
			i--
			leftEdge += 1024
			rightEdge += 1024
			ipList = []int{leftEdge, rightEdge}
		}
	}
	return res
}
func set2Map(set []int) map[int]lib.AntiAliasSet {
	resMap := map[int]lib.AntiAliasSet{}
	for i := 0; i < len(set)-1; i++ {
		resMap[set[i]] = lib.AntiAliasSet{}
		for j := i + 1; j < len(set); j++ {
			resMap[set[i]] = append(resMap[set[i]], set[j])
		}
	}
	return resMap
}
