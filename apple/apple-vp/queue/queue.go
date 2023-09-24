/*
 * @Author: fyfishie
 * @Date: 2023-05-08:10
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-08:10
 * @@email: fyfishie@outlook.com
 * @Description: :)
 */
package queue

import "vp/lib"

type Queue struct {
	Head   *lib.IpPair
	Tail   *lib.IpPair
	Length int
}

func (q *Queue) Push(pair lib.IpPair) {
	if q.Length == 0 {
		q.Head = &pair
		q.Tail = &pair
		q.Length++
	} else {
		q.Tail.Next = &pair
		q.Tail = &pair
		q.Length++
	}
}
func (q *Queue) Pop() *lib.IpPair {
	if q.Length == 0 {
		return nil
	} else {
		ret := q.Head
		q.Head = ret.Next
		q.Length--
		return ret
	}
}
func NewQueue() Queue {
	q := Queue{
		Head:   nil,
		Tail:   nil,
		Length: 0,
	}
	return q
}
