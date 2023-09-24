package subnet

import (
	"vp/lib"
	"vp/queue"
)

type gather struct {
	InputChan  chan []int
	OutPutChan chan [][]int
}

func (g *gather) gatherRun(num int) {
	for i := 0; i < num; i++ {
		go func() {
			for block := range g.InputChan {
				antiAliasSubNet := subNetGatherByDivision(block[2:], block[0], block[1])
				g.OutPutChan <- antiAliasSubNet
			}
		}()
	}
}

func newGather(inputChan chan ([]int), outPutChan chan [][]int) gather {
	g := gather{
		InputChan:  inputChan,
		OutPutChan: outPutChan,
	}
	return g
}

func subNetGatherByDivision(ipList []int, left int, right int) [][]int {
	res := [][]int{}
	q := queue.NewQueue()
	pair := lib.IpPair{
		Left:  left,
		Right: right,
	}
	q.Push(pair)
	for {
		if q.Length == 0 {
			break
		}
		p := q.Pop()
		num, edge, leftIndexInList, rightIndexInList := getIpNum(p.Left, p.Right, ipList)
		if num >= (p.Right-p.Left+1)/2 && !edge {
			iplist := make([]int, num)
			j := 0
			for i := leftIndexInList; i < rightIndexInList+1; i++ {
				iplist[j] = (ipList)[i]
				j++
			}
			res = append(res, iplist)
		} else {
			if p.Left+3 == p.Right {
				continue
			}
			lmid := (p.Left + p.Right - 1) / 2
			rmid := lmid + 1
			lpair := lib.IpPair{
				Left:  p.Left,
				Right: lmid,
			}
			rpair := lib.IpPair{
				Left:  rmid,
				Right: p.Right,
			}
			q.Push(lpair)
			q.Push(rpair)
		}
	}
	return res
}

func getIpNum(left int, right int, ipList []int) (num int, edge bool, leftIndexInList int, rightIndexInList int) {
	leftIndexInList = 0
	rightIndexInList = 0
	edge = false
	for index, value := range ipList {
		leftIndexInList = index
		if value < left {
		} else {
			if value == left {
				edge = true
			}
			break
		}
	}
	for i := len(ipList) - 1; i >= 0; i-- {
		rightIndexInList = i
		if (ipList)[i] > right {
		} else {
			if (ipList)[i] == right {
				edge = true
			}
			break
		}
	}
	return (rightIndexInList - leftIndexInList + 1), edge, leftIndexInList, rightIndexInList
}
