/*
 * @Author: fyfishie
 * @Date: 2023-03-21:08
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-14:10
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package onetrace

import (
	"fmt"
	"sort"
	"vp/lib"
	"vp/status"
)

type Parser struct {
	//send result out
	ResultOutChan chan lib.AntiDescriptor
	//raw resources input
	traceInChan chan lib.TraceRoute
	//status monitor
	Eye status.StatusEye
}

func NewParser(traceInChan chan lib.TraceRoute) *Parser {
	return &Parser{
		traceInChan:   traceInChan,
		ResultOutChan: make(chan lib.AntiDescriptor, 1000),
		Eye:           status.StatusEye{},
	}
}

func (p *Parser) WithRun() *Parser {
	go func() {
		count := 0
		for trace := range p.traceInChan {
			count++
			if count%50000 == 0 {
				fmt.Printf("count: %v\n", count)
			}
			sort.Ints(trace.Trace)
			for i := 0; i < len(trace.Trace)-1; i++ {
				ipArray := lib.AntiDescriptor{
					Ip:           trace.Trace[i],
					AntiAliasSet: trace.Trace[i+1:],
				}
				p.ResultOutChan <- ipArray
			}
		}
		fmt.Println("outChan closed")
		close(p.ResultOutChan)
	}()
	return p
}
