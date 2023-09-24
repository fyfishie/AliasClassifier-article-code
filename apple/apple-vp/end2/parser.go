/*
 * @Author: fyfishie
 * @Date: 2023-03-21:08
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-14:21
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package end2

import (
	"fmt"
	"sort"
	"vp/lib"
	"vp/status"
)

type Parser struct {
	traceInChan chan lib.TraceRoute
	//send result out
	ResultOutChan chan lib.AntiDescriptor

	// //how many traces later write data out
	// bufferLen int

	//status monitor port
	Eye status.StatusEye
}

func NewParser(traceInChan chan lib.TraceRoute) *Parser {
	p := Parser{
		traceInChan:   traceInChan,
		Eye:           status.StatusEye{},
		ResultOutChan: make(chan lib.AntiDescriptor, 100),
	}
	return &p
}

/*
returns some sets, each ip in one set is not alias of others of the same set
output from resultOutChan is not ordered
*/
func (p *Parser) WithRun() *Parser {
	go func() {
		end2Map := map[int]lib.AntiAliasSet{}
		for trace := range p.traceInChan {
			if end2Map[trace.End2] == nil {
				end2Map[trace.End2] = lib.AntiAliasSet{trace.End}
			} else {
				end2Map[trace.End2] = append(end2Map[trace.End2], trace.End)
			}
		}
		for _, antiSet := range end2Map {
			if len(antiSet) < 2 {
				continue
			}
			if len(antiSet) > 1000 {
				fmt.Println("?")
			}
			for _, d := range set2AntiDescriptor(antiSet) {
				p.ResultOutChan <- *d
			}
		}
		close(p.ResultOutChan)
	}()
	return p
}

func set2AntiDescriptor(set []int) []*lib.AntiDescriptor {
	res := []*lib.AntiDescriptor{}
	sort.Ints(set)
	for i := 0; i < len(set)-1; i++ {
		descriptor := lib.AntiDescriptor{}
		descriptor.Ip = set[i]
		descriptor.AntiAliasSet = set[i+1:]
		res = append(res, &descriptor)
	}
	return res
}
