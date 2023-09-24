/*
 * @Author: fyfishie
 * @Date: 2023-03-08:16
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-10:21
 * @Description: :)
 * @email: muren.zhuang@outlook.com
 */
package graph

type Graph struct {
	G map[int][]int
}

func NewGraph() Graph {
	g := Graph{
		G: map[int][]int{},
	}
	return g
}

func (g *Graph) AddEdge(nodeA int, nodeB int) {
	if g.G[nodeA] == nil {
		g.G[nodeA] = []int{nodeB}
	} else {
		g.G[nodeA] = append(g.G[nodeA], nodeB)
	}

	if g.G[nodeB] == nil {
		g.G[nodeB] = []int{nodeA}
	} else {
		g.G[nodeB] = append(g.G[nodeB], nodeA)
	}
}

func (g *Graph) MaxConnectedSubGraph() [][]int {
	res := [][]int{}
	distinctMap := map[int]bool{}
	q := NewQueue4Int()
	for keyIp, _ := range g.G {
		_, ok := distinctMap[keyIp]
		if ok {
			continue
		}
		q.Push(keyIp)
		distinctMap[keyIp] = true
		antiAliasSet := []int{}
		for {
			if q.Length == 0 {
				break
			}
			ip := q.Pop()
			antiAliasSet = append(antiAliasSet, ip)
			for _, v := range g.G[ip] {
				_, ok := distinctMap[v]
				if !ok {
					q.Push(v)
					distinctMap[v] = true
				}
			}
		}
		res = append(res, antiAliasSet)
	}
	return res
}
