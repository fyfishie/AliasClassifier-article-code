/*
 * @Author: fyfishie
 * @Date: 2023-02-20:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-04-07:16
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	WRITE_RESULT_TIMEOUT = iota
	WRITE_RESULT_OK
	WRITE_RESULT_RECONNECTING
	WRITE_RESULT_JSON_MARSHAL_ERROR
	WRITE_RESULT_PUBLISH_ERR
)

type rabbit struct {
	uri          string
	retry        int
	connection   *amqp.Connection
	channel      *amqp.Channel
	reconnecting bool
	errSigChan   chan *amqp.Error
	confirms     chan amqp.Confirmation
}

type QueueInfo struct {
	Name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	nowait     bool
	args       map[string]interface{}
}
type ExchangeInfo struct {
	Name       string
	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	args       map[string]interface{}
}
type publishTask struct {
	result chan int
	datas  map[string][]byte
}
