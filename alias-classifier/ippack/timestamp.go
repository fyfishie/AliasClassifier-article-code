/*
 * @Author: fyfishie
 * @Date: 2023-04-25:21
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-02:09
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package ippack

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// construct ip datagrams which are expected as request of ipid data
func MakeICMPEchoBody() []byte {
	msg := icmp.Message{
		Type:     ipv4.ICMPTypeEcho,
		Code:     0,
		Checksum: 0,
		Body: &icmp.RawBody{
			Data: []byte(strconv.FormatInt(time.Now().UnixMilli(), 10)),
		},
	}
	bs, _ := msg.Marshal(nil)
	fmt.Println(bs)
	return bs
}
