/*
 * @Author: fyfishie
 * @Date: 2023-04-25:20
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-07-02:16
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package statistic

import (
	"alias_article/lib"
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/fyfishie/dorapock/store"
	"github.com/fyfishie/esyerr"
)

// process ipid data according to article
// "Fixing Ally’s Growing Pains with Velocity Modeling"
func MakeIPIDArray(ipidPath string, wtPath string) {
	ipids, err := store.LoadAny[lib.IPAID](ipidPath)
	IPPointsMap := map[string]lib.IDTimeArray{}
	wfi, err := os.OpenFile(wtPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	for _, ipid := range ipids {
		if _, ok := IPPointsMap[ipid.IP]; !ok {
			IPPointsMap[ipid.IP] = lib.IDTimeArray{}
		}
		IPPointsMap[ipid.IP] = append(IPPointsMap[ipid.IP], lib.IDPoint{IPID: ipid.ID, Time: ipid.TransmitTimeStamp})
	}
	for ip, array := range IPPointsMap {
		sort.Sort(&array)
		t := lib.IPIDTimes{
			IP:          ip,
			IDTimeArray: array,
		}
		bs, _ := json.Marshal(zeroit(dewarp(t)))
		wtr.Write(bs)
		wtr.WriteString("\n")
	}
	wtr.Flush()
}
func zeroit(d lib.IPIDTimes) lib.IPIDTimes {
	zero := d.IDTimeArray[0]
	for index, dp := range d.IDTimeArray {
		d.IDTimeArray[index].IPID = dp.IPID - zero.IPID
		d.IDTimeArray[index].Time = dp.Time - zero.Time
	}
	return d
}
func dewarp(ipidtimes lib.IPIDTimes) lib.IPIDTimes {
	for {
		c := false
		for i := 1; i < len(ipidtimes.IDTimeArray); i++ {
			if ipidtimes.IDTimeArray[i].IPID < ipidtimes.IDTimeArray[i-1].IPID {
				c = true
				ipidtimes.IDTimeArray[i].IPID += 65535
			}
		}
		if !c {
			break
		}
	}
	return ipidtimes
}

func IPIDSta(ipidPath string, pairPath string, wtpath string) {
	ipids, err := store.LoadAny[lib.IPIDTimes](ipidPath)
	esyerr.AutoPanic(err)
	ipidMap := map[string]lib.IPIDTimes{}
	for _, ipid := range ipids {
		ipidMap[ipid.IP] = ipid
	}
	pairs, err := store.LoadAny[lib.Pair](pairPath)
	esyerr.AutoPanic(err)
	disMap := map[float64]int{}
	for _, pair := range pairs {
		if _, ok := ipidMap[pair.IPA]; !ok {
			continue
		}
		if _, ok := ipidMap[pair.IPB]; !ok {
			continue
		}
		disMap[distance(ipidMap[pair.IPA], ipidMap[pair.IPB])]++
	}
	WriteIPID(disMap, wtpath, "")
}

func distance(ipidA, ipidB lib.IPIDTimes) float64 {
	times := allTimes(ipidA, ipidB)
	leftAB, leftPs := leftEnd(ipidA, ipidB, times)
	rightAB, rightPs := rightEnd(ipidA, ipidB, times)
	mAPs, mBPs := middle(ipidA, ipidB, leftPs, rightPs)
	var sum float64
	sum += leftSum(leftAB, leftPs, ipidA, ipidB)
	sum += rightSum(rightAB, rightPs, ipidA, ipidB)
	sum += middleSum(mAPs, mBPs, ipidA, ipidB)
	if sum > 1000000000 {
		fmt.Println("?")
		leftAB, leftPs := leftEnd(ipidA, ipidB, times)
		rightAB, rightPs := rightEnd(ipidA, ipidB, times)
		mAPs, mBPs := middle(ipidA, ipidB, leftPs, rightPs)
		sum1 := leftSum(leftAB, leftPs, ipidA, ipidB)
		sum2 := rightSum(rightAB, rightPs, ipidA, ipidB)
		sum3 := middleSum(mAPs, mBPs, ipidA, ipidB)
		fmt.Println(sum1 + sum2 + sum3)
	}
	return sum
}
func middleSum(mAPs, mBPs []lib.IDPoint, ipidA, ipidB lib.IPIDTimes) float64 {
	var sum float64
	for _, ap := range mAPs {
		leftClose := lib.IDPoint{Time: math.MinInt}
		rightClose := lib.IDPoint{Time: math.MaxInt}
		for _, bID := range ipidB.IDTimeArray {
			if ap.Time == bID.Time {
				sum += math.Abs(float64(ap.IPID) - float64(bID.IPID))
				break
			}
			if bID.Time < ap.Time {
				if bID.Time > leftClose.Time {
					leftClose = bID
				}
			}
			if bID.Time > ap.Time {
				if bID.Time < rightClose.Time {
					rightClose = bID
				}
			}
		}
		if leftClose.Time == math.MinInt || rightClose.Time == math.MaxInt {
			continue
		}
		sum += est(leftClose, rightClose, ap)
	}

	for _, bp := range mBPs {
		leftClose := lib.IDPoint{Time: math.MinInt}
		rightClose := lib.IDPoint{Time: math.MaxInt}
		for _, aID := range ipidA.IDTimeArray {
			if bp.Time == aID.Time {
				sum += math.Abs(float64(bp.IPID - aID.IPID))
				break
			}
			if aID.Time < bp.Time {
				if aID.Time > leftClose.Time {
					leftClose = aID
				}
			}
			if aID.Time > bp.Time {
				if aID.Time < rightClose.Time {
					rightClose = aID
				}
			}
		}
		if leftClose.Time == math.MinInt || rightClose.Time == math.MaxInt {
			continue
		}
		sum += est(leftClose, rightClose, bp)
	}
	return sum
}
func est(p1, p2 lib.IDPoint, dp lib.IDPoint) float64 {
	var k float64
	var b float64
	k = (float64(p1.Time - p2.Time)) / (float64(p1.IPID - p2.IPID))
	b = float64(p1.IPID) - k*(float64(p1.Time))
	gap := float64(dp.IPID) - k*(float64(dp.Time)) - b
	return math.Abs(gap)
}
func leftSum(leftAB string, leftPs []lib.IDPoint, ipidA, ipidB lib.IPIDTimes) float64 {
	if len(leftPs) == 0 {
		return 0
	}
	leftRightTime := leftPs[len(leftPs)-1].Time
	idsk := ipidB
	if leftAB == "B" {
		idsk = ipidA
	}
	firstP, secondP := lib.IDPoint{}, lib.IDPoint{}
	for _, id := range idsk.IDTimeArray {
		if id.Time > leftRightTime {
			firstP = id
			break
		}
	}
	in := 0
	for _, id := range idsk.IDTimeArray {
		if id.Time > leftRightTime {
			in++
		}
		if in == 2 {
			secondP = id
			break
		}
	}
	var k float64
	var b float64
	k = (float64(secondP.IPID - firstP.IPID)) / (float64(secondP.Time - secondP.IPID))
	b = float64(firstP.IPID) - k*(float64(firstP.Time))
	var sum float64 = 0
	for _, p := range leftPs {
		gap := k*float64(p.Time) + b - float64(p.IPID)
		sum += math.Abs(gap)
	}
	return sum
}
func rightSum(rightAB string, rightPs []lib.IDPoint, ipidA, ipidB lib.IPIDTimes) float64 {
	if len(rightPs) == 0 {
		return 0
	}
	rightLeftTime := rightPs[0].Time
	idsk := ipidB
	if rightAB == "B" {
		idsk = ipidA
	}
	rightFirstP, rightSecondP := lib.IDPoint{}, lib.IDPoint{}
	for i := len(idsk.IDTimeArray) - 1; i > -1; i-- {
		id := idsk.IDTimeArray[i]
		if id.Time < rightLeftTime {
			rightFirstP = id
		}
	}
	in := 0
	for i := len(idsk.IDTimeArray) - 1; i > -1; i-- {
		id := idsk.IDTimeArray[i]
		if id.Time < rightLeftTime {
			in++
		}
		if in == 2 {
			rightSecondP = id
			break
		}
	}
	var k float64
	var b float64
	k = (float64(rightSecondP.IPID - rightFirstP.IPID)) / (float64(rightSecondP.Time - rightSecondP.IPID))
	b = float64(rightFirstP.IPID) - k*(float64(rightFirstP.Time))
	var sum float64 = 0
	for _, p := range rightPs {
		gap := k*float64(p.Time) + b - float64(p.IPID)
		sum += math.Abs(gap)
	}
	return sum
}
func allTimes(ipidA, ipidB lib.IPIDTimes) []int {
	res := []int{}
	mid := map[int]struct{}{}
	for _, a := range ipidA.IDTimeArray {
		mid[a.Time] = struct{}{}
	}
	for _, b := range ipidB.IDTimeArray {
		mid[b.Time] = struct{}{}
	}
	for k, _ := range mid {
		res = append(res, k)
	}
	sort.Ints(res)
	return res
}
func middle(ipidA, ipidB lib.IPIDTimes, lefttimes []lib.IDPoint, rightTimes []lib.IDPoint) (APoints, BPoints []lib.IDPoint) {
	aps := []lib.IDPoint{}
	bps := []lib.IDPoint{}
	for _, a := range ipidA.IDTimeArray {
		hs := false
		for _, t := range lefttimes {
			hs = hs && (a.Time == t.Time)
		}
		for _, t := range rightTimes {
			hs = hs && (a.Time == t.Time)
		}
		if !hs {
			aps = append(aps, a)
		}
	}

	for _, b := range ipidB.IDTimeArray {
		hs := false
		for _, t := range lefttimes {
			hs = hs && (b.Time == t.Time)
		}
		for _, t := range rightTimes {
			hs = hs && (b.Time == t.Time)
		}
		if !hs {
			bps = append(bps, b)
		}
	}
	return aps, bps
}
func leftEnd(ipidA, ipidB lib.IPIDTimes, times []int) (AB string, res []lib.IDPoint) {
	left := ""
	for _, time := range times {
		if ipidA.IDTimeArray[0].Time == time {
			left = "A"
		}
	}
	if left == "" {
		left = "B"
	}
	res = []lib.IDPoint{}
	for _, time := range times {
		for _, a := range ipidA.IDTimeArray {
			if a.Time == time {
				if left == "A" {
					res = append(res, lib.IDPoint{IPID: a.IPID, Time: a.Time})
				} else {
					return left, res
				}
			}
		}

		for _, b := range ipidB.IDTimeArray {
			if b.Time == time {
				if left == "B" {
					res = append(res, lib.IDPoint{IPID: b.IPID, Time: b.Time})
				} else {
					return left, res
				}
			}
		}
	}
	return left, res
}
func rightEnd(ipidA, ipidB lib.IPIDTimes, times []int) (AB string, res []lib.IDPoint) {
	right := ""
	for i := len(times) - 1; i > -1; i-- {
		if ipidA.IDTimeArray[0].Time == times[i] {
			right = "A"
		}
	}
	if right == "" {
		right = "B"
	}
	res = []lib.IDPoint{}
	for i := len(times) - 1; i > -1; i-- {
		time := times[i]
		for _, a := range ipidA.IDTimeArray {
			if a.Time == time {
				if right == "A" {
					res = append(res, lib.IDPoint{IPID: a.IPID, Time: a.Time})
				} else {
					return right, res
				}
			}
		}

		for _, b := range ipidB.IDTimeArray {
			if b.Time == time {
				if right == "B" {
					res = append(res, lib.IDPoint{IPID: b.IPID, Time: b.Time})
				} else {
					return right, res
				}
			}
		}
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return right, res
}
