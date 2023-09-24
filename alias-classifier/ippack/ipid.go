/*
 * @Author: fyfishie
 * @Date: 2023-04-23:16
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-02:20
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package ippack

import (
	"alias_article/lib"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/fyfishie/esyerr"
	"golang.org/x/net/ipv4"
)

func ListenIPID(wtpath string) {
	wfi, err := os.OpenFile(wtpath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	esyerr.AutoPanic(err)
	defer wfi.Close()
	wtr := bufio.NewWriter(wfi)
	conn, err := net.ListenIP("ip4:icmp", nil)
	esyerr.AutoPanic(err)
	rawConn, err := ipv4.NewRawConn(conn)
	esyerr.AutoPanic(err)
	buf := make([]byte, 1024)
	// tkr := time.NewTicker(time.Second * 5)
	count := 0
	go func() {
		for {
			wtr.Flush()
			time.Sleep(time.Second * 10)
		}
	}()
	for {
		count++
		if count%10000 == 0 {
			fmt.Printf("count: %v\n", count)
		}
		h, p, _, _ := rawConn.ReadFrom(buf)
		ipid := lib.IPAID{
			IP:                 h.Src.String(),
			ID:                 h.ID,
			SequenceNum:        int(p[7]),
			OriginateTimeStamp: int(binary.BigEndian.Uint32(p[8:12])),
			ReceiveTimeStamp:   int(binary.BigEndian.Uint32(p[12:16])),
			TransmitTimeStamp:  int(binary.BigEndian.Uint32(p[16:20])),
		}
		bs, _ := json.Marshal(ipid)
		wtr.Write(bs)
		wtr.WriteString("\n")
	}
}
