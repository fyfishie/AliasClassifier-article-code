/*
 * @Author: fyfishie
 * @Date: 2023-05-30:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-06-21:15
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package lib

import (
	"strconv"

	"github.com/fyfishie/ipop"
)

type Pair struct {
	IPA   string `json:"ipa"`
	IPB   string `json:"ipb"`
	ANode string `json:"a_node"`
	BNode string `json:"b_node"`
}

func (i *Pair) ID() string {
	ipa := ipop.String2Int(i.IPA)
	ipb := ipop.String2Int(i.IPB)
	if ipa < ipb {
		return strconv.Itoa(ipa) + strconv.Itoa(ipb)
	} else {
		return strconv.Itoa(ipb) + strconv.Itoa(ipa)
	}
}
