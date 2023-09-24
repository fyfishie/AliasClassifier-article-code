/*
 * @Author: fyfishie
 * @Date: 2023-05-14:11
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-14:19
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package end2

import (
	"fmt"
	"sort"
	"vp/lib"
)

type set2DescriptorConvor struct {
	inChan  chan []int
	outChan chan lib.AntiDescriptor
}

func newCovor(resultoutChan chan lib.AntiDescriptor) *set2DescriptorConvor {
	return &set2DescriptorConvor{
		inChan:  make(chan []int),
		outChan: resultoutChan,
	}
}
func (c *set2DescriptorConvor) withRun() *set2DescriptorConvor {
	go func() {
		for set := range c.inChan {
			if len(set) > 50000 {
				fmt.Println("?")
			}
			sort.Ints(set)
			for i := 0; i < len(set)-1; i++ {
				descriptor := lib.AntiDescriptor{}
				descriptor.Ip = set[i]
				descriptor.AntiAliasSet = set[i+1:]
				c.outChan <- descriptor
			}
		}
		close(c.outChan)
	}()
	return c
}
