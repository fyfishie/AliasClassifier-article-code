/*
 * @Author: fyfishie
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-08-02:09
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package data4sk

import (
	"bufio"
	"os"
	"strings"
)

// merge 2 charactor vectors file to get origin data for machine learn of multi-vp
func Merge2(first, second string, wtpath string, appendID string) {
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	rfi, err := os.Open(first)
	if err != nil {
		panic(err.Error())
	}
	defer rfi.Close()
	scaner1 := bufio.NewScanner(rfi)
	scaner1.Scan()
	head := scaner1.Text()
	rfi1, err := os.Open(second)
	if err != nil {
		panic(err.Error())
	}
	defer rfi1.Close()
	scaner2 := bufio.NewScanner(rfi1)
	ss := strings.Split(head, ",")
	for i := 0; i < len(ss)-1; i++ {
		wtr.WriteString(ss[i] + ",")
	}
	scaner2.Scan()
	head = scaner2.Text()
	ss = strings.Split(head, ",")
	for i := 0; i < len(ss)-1; i++ {
		if i == 5 || i == 6 || i == 7 || i == 8 || i == 9 {
			continue
		}
		wtr.WriteString(ss[i] + appendID + ",")
	}
	wtr.WriteString(ss[len(ss)-1] + "\n")
	for scaner1.Scan() && scaner2.Scan() {
		ss := strings.Split(scaner1.Text(), ",")
		for i := 0; i < len(ss)-1; i++ {
			if i == 5 || i == 6 || i == 7 || i == 8 || i == 9 {
				continue
			}
			wtr.WriteString(ss[i] + ",")
		}
		wtr.WriteString(scaner2.Text() + "\n")
	}
	wtr.Flush()
}
