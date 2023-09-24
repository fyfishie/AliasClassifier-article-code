/*
 * @Author: fyfishie
 * @Date: 2023-04-02:20
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-12:10
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package rabbitmq

import (
	"fmt"
	"testing"
)

func Test_pullcon(t *testing.T) {
	// go pub()
	exchangeInfo := ExchangeInfo{
		Name:       "ex_test",
		Kind:       "direct",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     true,
		Args:       nil,
	}
	queueInfo := QueueInfo{
		Name:       "t_queue",
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		nowait:     true,
		args:       nil,
	}
	consumer := NewConsumer("amqp://alias_parse:tyWFX2hJ@81.70.76.237:5672/alias_parse_vhost", 5, "tkey", queueInfo, exchangeInfo)
	taskin, err := consumer.RunAsPullMod()
	if err != nil {
		panic(err.Error())
	}
	for delivery := range taskin {
		fmt.Println(string(delivery.Body))
		delivery.Ack(false)
	}

}

func Test_pushCon(t *testing.T) {
	exchangeInfo := ExchangeInfo{
		Name:       "ex_test",
		Kind:       "direct",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     true,
		Args:       nil,
	}
	queueInfo := QueueInfo{
		Name:       "t_queue",
		durable:    true,
		autoDelete: false,
		exclusive:  false,
		nowait:     true,
		args:       nil,
	}
	consumer := NewConsumer("amqp://alias_parse:tyWFX2hJ@81.70.76.237:5672/alias_parse_vhost", 5, "tkey", queueInfo, exchangeInfo)
	taskin, err := consumer.RunAsPushMod(1)
	if err != nil {
		panic(err.Error())
	}
	for d := range taskin {
		fmt.Println(string(d.Body))
		d.Ack(false)
	}
}
