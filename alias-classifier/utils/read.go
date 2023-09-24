/*
 * @Author: fyfishie
 * @Date: 2023-04-22:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-04-26:19
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package utils

import (
	"bufio"
	"github.com/fyfishie/esyerr"
	"os"
)

func ReadDir(dirpath string) []string {
	etrys, err := os.ReadDir(dirpath)
	esyerr.AutoPanic(err)
	files := []string{}
	for _, e := range etrys {
		files = append(files, e.Name())
	}
	return files
}

func ReadALLLine(path string) []string {
	rfi, err := os.OpenFile(path, os.O_RDONLY, 0000)
	if err != nil {
		panic(err)
	}
	rdr := bufio.NewScanner(rfi)
	res := []string{}
	for rdr.Scan() {
		res = append(res, rdr.Text())
	}
	return res
}
