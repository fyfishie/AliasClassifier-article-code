/*
 * @Author: fyfishie
 * @Date: 2023-04-02:15
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-04-07:16
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
package rabbitmq

import (
	"testing"
)

func Test_pub(t *testing.T) {
	// exchangeInfo := ExchangeInfo{
	// 	Name:       "ex_test",
	// 	kind:       "direct",
	// 	durable:    true,
	// 	autoDelete: false,
	// 	internal:   false,
	// 	noWait:     true,
	// 	args:       nil,
	// }
	// p := NewPublisher("amqp://antialias:qAwR3Y@81.70.76.237:5672/alias_vhost", 5, exchangeInfo, 10, context.TODO(), Confirm_Mode)
	// _, err := p.Run()
	// if err != nil {
	// 	panic(err.Error())
	// }
	// for {
	// 	time.Sleep(time.Second)
	// 	// bs, _ := json.Marshal("hello rabbit!")
	// 	data := []any{}
	// 	data = append(data, "bs")
	// 	err := p.Publish(data, "tkey")
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// }
}
