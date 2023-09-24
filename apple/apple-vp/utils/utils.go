/*
 * @Author: fyfishie
 * @Date: 2023-03-01:20
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-08:17
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package utils

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strings"
	"vp/lib"

	"github.com/fyfishie/ipop"
)

var ISIPREG = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)

func MakeMessage(status bool, message string) lib.Message {
	return lib.Message{
		Status:  status,
		Message: message,
	}
}

func MakeMessageBytes(status bool, message string) []byte {
	m := lib.Message{
		Status:  status,
		Message: message,
	}
	bs, _ := json.Marshal(m)
	return bs
}
func MakeMessageBytesWithID(status bool, id lib.ID, msg string) []byte {
	bs, _ := json.Marshal(lib.MessageWithID{Status: status, ID: id, Message: msg})
	return bs
}

func Gcd(a int, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// 计算一个数组的Gcd，输入必须为非负自然数,
func Gcds(list []int) int {
	currentGcd := Gcd(list[0], list[1])
	if len(list) > 2 {
		currentGcd = Gcd(list[0], list[1])
		for i := 1; i < len(list); i++ {
			currentGcd = Gcd(currentGcd, list[i])
		}
	}
	return currentGcd
}

// 16843009 -> 1.1.1.1
func IPConv_atoi(i int) string {
	a := i & 255
	b := i & 65280 >> 8
	c := i & 16711680 >> 16
	d := i & 4278190080 >> 24
	return IntToStringMap[d] + "." + IntToStringMap[c] + "." + IntToStringMap[b] + "." + IntToStringMap[a]
}

// 1.1.1.1 -> 16843009
func IPConv_itoa(s string) int {
	sArray := strings.Split(s, ".")
	var res int = 0
	res = String2IntMap[sArray[3]] + String2IntMap[sArray[2]]*256 + String2IntMap[sArray[1]]*65536 + String2IntMap[sArray[0]]*16777216
	return res
}

// check whether a string is in valid ipv4 format or not by regexp
func IsIP(s string) bool {
	match := ISIPREG.FindStringSubmatch(s)
	return len(match) == 1
}

// return a list of ip which is in the correct format ipv4 from input ip list
func FiltIPList(ips []string) []string {
	validIPs := []string{}
	for _, ip := range ips {
		if IsIP(ip) {
			validIPs = append(validIPs, ip)
		}
	}
	return validIPs
}

func IPsConv_atoi(ss []string) []int {
	is := make([]int, len(ss))
	for index, s := range ss {
		is[index] = IPConv_itoa(s)
	}
	return is
}

func IPsConv_itoa(is []int) []string {
	ss := make([]string, len(is))
	for index, i := range is {
		ss[index] = IPConv_atoi(i)
	}
	return ss
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

/*
 * @description: convert lib.RawTrace to lib.TraceRoute
 * @param {[]lib.RawTrace} Raws
 * @return {*}
 */
func RawTrace2Trace(input lib.RawTrace) lib.TraceRoute {
	length := len(input.Results)
	trace := lib.TraceRoute{
		Target:   IPConv_itoa(input.Ip),
		End:      IPConv_itoa(input.Results[length-1].Ip),
		TTLStart: input.Results[0].TTL,
	}
	if length > 1 {
		trace.End2 = IPConv_itoa(input.Results[length-2].Ip)
	}
	for i := 0; i < length; i++ {
		trace.Trace = append(trace.Trace, IPConv_itoa(input.Results[i].Ip))
		trace.Times = append(trace.Times, float64(input.Results[i].Rtt))
		trace.TTLS = append(trace.TTLS, input.Results[i].TTL)
		trace.Status = append(trace.Status, input.Results[i].Status)
	}
	return trace
}

/*
 * @description: convert more RawTrace at once
 * @param {[]lib.RawTrace} input
 * @return {*}
 */
func RawTraces2Traces(inputs []lib.RawTrace) []lib.TraceRoute {
	res := []lib.TraceRoute{}
	for _, input := range inputs {
		length := len(input.Results)
		trace := lib.TraceRoute{
			Target:   IPConv_itoa(input.Ip),
			End:      IPConv_itoa(input.Results[length-1].Ip),
			TTLStart: input.Results[0].TTL,
		}
		if length > 1 {
			trace.End2 = IPConv_itoa(input.Results[length-2].Ip)
		}
		for i := 0; i < length; i++ {
			trace.Trace = append(trace.Trace, IPConv_itoa(input.Results[i].Ip))
			trace.Times = append(trace.Times, float64(input.Results[i].Rtt))
			trace.TTLS = append(trace.TTLS, input.Results[i].TTL)
			trace.Status = append(trace.Status, input.Results[i].Status)
		}
		res = append(res, trace)
	}
	return res
}

// // get public ip addr of localhost
//
//	func GetLocalIPList() (publicIPList []string) {
//		interfaceAddrs, err := net.InterfaceAddrs()
//		if err != nil {
//			logrus.Panic("failed to get net interface addr, err = " + err.Error())
//		}
//		for _, addr := range interfaceAddrs {
//			ipNet, isValidIpNet := addr.(*net.IPNet)
//			if isValidIpNet && !ipNet.IP.IsLoopback() {
//				if ipNet.IP.To4() != nil {
//					if ipNet.IP.To4() != nil && !ipNet.IP.IsPrivate() {
//						publicIPList = append(publicIPList, ipNet.IP.String())
//					}
//				}
//			}
//		}
//		return publicIPList
//	}
// func GetLocalIPList() []string {
// 	return []string{"127.0.0.1"}
// }

func ParseMessage(bs []byte) (lib.Message, error) {
	m := lib.Message{}
	err := json.Unmarshal(bs, &m)
	if err != nil {
		return lib.Message{}, err
	}
	return m, nil
}
func Load[T any](rdpath string) ([]*T, error) {
	rfi, err := os.OpenFile(rdpath, os.O_RDONLY, 0000)
	if err != nil {
		return nil, err
	}
	defer rfi.Close()
	rdr := bufio.NewReader(rfi)
	res := []*T{}
	for {
		line, _, err := rdr.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		var item T
		err = json.Unmarshal(line, &item)
		if err != nil {
			continue
		}
		res = append(res, &item)
	}
	return res, nil
}

func StringIPs2Ints(ipList []string) []int {
	res := []int{}
	for _, ip := range ipList {
		res = append(res, ipop.String2Int(ip))
	}
	return res
}
