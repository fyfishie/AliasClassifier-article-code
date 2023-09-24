/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-07-02:16
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package statistic

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fyfishie/esyerr"
)

// write data out in CDF form
func WriteAccuDataFrom(data map[int]int, wtpath string, fileHead string) {
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	wtr.WriteString(fileHead)
	total := 0
	//累计
	keyMax := 0
	for k, cnt := range data {
		if k > keyMax {
			keyMax = k
		}
		total += cnt
	}
	vTotal := 0
	for i := 0; i < keyMax+1; i++ {
		if cnt, ok := data[i]; ok {
			vTotal += cnt
		}
		wtr.WriteString(fmt.Sprintf("%d,%.3f\n", i, float64(vTotal)/float64(total)))
	}
	esyerr.AutoPanic(wtr.Flush())
}

// write data out in CDF form
func WriteAccuDataFromFloat(data map[float64]int, wtpath string, fileHead string, step float64) {
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	wtr.WriteString(fileHead)
	total := 0
	//累计
	var keyMax float64
	for k, cnt := range data {
		if k > keyMax {
			keyMax = k
		}
		total += cnt
	}
	var i float64
	for i = 0; i <= keyMax; i += step {
		vTotal := 0
		for k, v := range data {
			if k <= i {
				vTotal += v
			}
		}
		wtr.WriteString(fmt.Sprintf("%.3f,%.3f\n", i, float64(vTotal)/float64(total)))
	}
	esyerr.AutoPanic(wtr.Flush())
}

// write ipid data out
func WriteIPID(data map[float64]int, wtpath string, fileHead string) {
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	wtr.WriteString(fileHead)
	total := 0
	//累计
	var keyMax float64
	for k, cnt := range data {
		if k > keyMax {
			keyMax = k
		}
		total += cnt
	}
	var i float64
	var step float64 = 10
	fmt.Printf("max k is:%v\n", keyMax)
	for i = 10; i <= keyMax; i *= step {
		vTotal := 0
		for k, v := range data {
			if k <= i {
				vTotal += v
			}
		}
		wtr.WriteString(fmt.Sprintf("%.3f,%.3f\n", i, float64(vTotal)/float64(total)))
	}
	esyerr.AutoPanic(wtr.Flush())
}
