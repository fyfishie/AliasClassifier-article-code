package utils

import (
	"aliasParseMaster/lib"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var ISIPREG = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)

func MakeMessage(status bool, message string) lib.Message {
	return lib.Message{
		Status:  status,
		Message: message,
	}
}

func MakeMessageBytesWithID(workID int, status bool, message string) []byte {
	m := lib.MessageWithID{
		WorkID:  workID,
		Status:  status,
		Message: message,
	}
	bs, _ := json.Marshal(m)
	return bs
}

func MakeMessageBytes(status bool, message string) []byte {
	m := lib.Message{
		Status:  status,
		Message: message,
	}
	bs, _ := json.Marshal(m)
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
func Float642Bytes(f float64) []byte {
	bits := math.Float64bits(f)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
func Bytes2Float64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func MustFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
func VPEquals(a, b lib.VP) bool {
	return a.UID == b.UID && a.URL == b.URL
}

func UniqueMongoCollection(topTaskID, ChildID int) (collectionName string) {
	return strconv.Itoa(topTaskID) + "_" + strconv.Itoa(ChildID)
}

func ParseVPNameAndVPID(uid string) (vpName, vpID string, err error) {
	reg := regexp.MustCompile(`([\W]+?_(\d+))`)
	match := reg.FindStringSubmatch(uid)
	if len(match) != 3 {
		return "", "", errors.New("invalid format of vp UID")
	}
	return match[1], match[2], nil
}
func ParseMessage(bs []byte) (lib.Message, error) {
	m := lib.Message{}
	err := json.Unmarshal(bs, &m)
	if err != nil {
		return lib.Message{}, err
	}
	return m, nil
}
func IsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
func MustInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
