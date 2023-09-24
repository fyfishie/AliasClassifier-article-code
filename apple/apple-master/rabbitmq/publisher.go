/*
 * @Author: fyfishie
 * @Date: 2023-04-07:15
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-05-13:16
 * @Description: :)
 * @email: fyfishie@outlook.com
 */
/*
 * @Author: fyfishie
 * @Date: 2023-02-20:09
 * @LastEditors: fyfishie
 * @LastEditTime: 2023-04-07:15
 * @Description: provide publisher publishes data with publish-key only
 * @email: fyfishie@outlook.com
 */
package rabbitmq

import (
	"context"
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const Confirm_Mode = 0
const Tx_Mode = 1

type Publisher struct {
	rabbit       rabbit
	exchangeInfo ExchangeInfo

	connected bool
	//publish消息的超时时间/second
	timeout     int
	publishChan chan publishTask
	ctx         context.Context
}

func NewPublisher(uri string, retry int, exchangeInfo ExchangeInfo, publishTimeout int, ctx context.Context) *Publisher {
	p := Publisher{
		rabbit:       rabbit{uri: uri, retry: retry},
		timeout:      publishTimeout,
		ctx:          ctx,
		exchangeInfo: exchangeInfo,
		connected:    false,
	}
	return &p
}
func (p *Publisher) Run() (err error) {
	p.publishChan = make(chan publishTask, 100)
	for i := 0; i < p.rabbit.retry; i++ {
		err = p.connect()
		if err != nil {
			p.clean()
			continue
		} else {
			break
		}
	}
	if err != nil {
		log.Fatalf("error in connect rabbitmq, err = %v\n", err.Error())
		return err
	}
	go func() {
		for {
			select {
			case task := <-p.publishChan:
				p.sendMany(task)
			case err := <-p.rabbit.errSigChan:
				p.closeProcess(err)
			}
		}
	}()
	return
}
func (p *Publisher) connect() (err error) {
	logrus.Infoln("publisher tring to connect... ")
	connection, err := amqp.Dial(p.rabbit.uri)
	if err != nil {
		log.Printf("error in dial rabbitMQ, err =  %v\n", err.Error())
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		log.Printf("error in get channel, err =  %v\n", err.Error())
		return err
	}
	p.rabbit.channel = channel

	p.rabbit.errSigChan = channel.NotifyClose(make(chan *amqp.Error))
	err = p.rabbit.channel.Tx()
	if err != nil {
		return err
	}
	// exchangeInfo := p.exchangeInfo
	// err = channel.ExchangeDeclare(exchangeInfo.Name, exchangeInfo.Kind, exchangeInfo.Durable, exchangeInfo.AutoDelete, exchangeInfo.Internal, exchangeInfo.NoWait, exchangeInfo.Args)
	// if err != nil {
	// 	log.Printf("error in declare exchange, err =  %v\n", err.Error())
	// 	return err
	// }
	p.connected = true
	return nil
}

func (p *Publisher) clean() {
	if p.rabbit.channel != nil {
		err := p.rabbit.channel.Close()
		if err != nil {
			log.Printf("error in close rebbit channel of publisher, err = %v\n", err.Error())
		}
	}
	if p.rabbit.connection != nil {
		if !p.rabbit.connection.IsClosed() {
			err := p.rabbit.connection.Close()
			if err != nil {
				log.Printf("error in close rabbit connect of publisher, err = %v\n", err.Error())
			}
		}
	}
	// if p.rabbit.errSigChan != nil {
	// 	close(p.rabbit.errSigChan)
	// }
}
func (p *Publisher) closeProcess(err error) {
	logrus.Errorf("lose connection to rabbitmq, tring to reconnect...")
	p.connected = false
	for i := 0; i < p.rabbit.retry; i++ {
		err := p.connect()
		if err != nil {
			p.clean()
			continue
		}
		p.connected = true
		break
	}
	if !p.connected {
		log.Fatalf("error in reconnect rabbitmq, err = %v\n", err.Error())
	}
}
func (p *Publisher) sendMany(task publishTask) {
	if !p.connected {
		task.result <- WRITE_RESULT_RECONNECTING
		return
	}
	err := p.rabbit.channel.Tx()
	if err != nil {
		task.result <- WRITE_RESULT_PUBLISH_ERR
	}
	for key, bs := range task.datas {
		err = p.rabbit.channel.PublishWithContext(p.ctx, p.exchangeInfo.Name, key, true, false, amqp.Publishing{
			Timestamp:    time.Now(),
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         bs,
		})
	}
	if err != nil {
		task.result <- WRITE_RESULT_PUBLISH_ERR
		err = p.rabbit.channel.TxRollback()
	}
	if err != nil {
		logrus.Error("error in rollback rabbitmq!!!, err = " + err.Error())
		return
	}
	err = p.rabbit.channel.TxCommit()
	if err != nil {
		task.result <- WRITE_RESULT_PUBLISH_ERR
		err = p.rabbit.channel.TxRollback()
	}
	if err != nil {
		logrus.Error("error in rollback rabbitmq!!!, err = " + err.Error())
		return
	}
	task.result <- WRITE_RESULT_OK
}

// datas: map[key]data
func (p *Publisher) Publish(datas map[string][]byte) error {
	resultChan := make(chan int)
	task := publishTask{
		result: resultChan,
		datas:  datas,
	}
	p.publishChan <- task
	res := <-resultChan
	if res != WRITE_RESULT_OK {
		return errors.New("failed to publish tx")
	}
	return nil
}
