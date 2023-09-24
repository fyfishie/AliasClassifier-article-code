package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	uri          string
	queueInfo    QueueInfo
	exchangeInfo ExchangeInfo
	rabbit       rabbit
	connected    bool
	bindKey      string
	retry        int
}

func NewConsumer(uri string, retry int, bindKey string, queueInfo QueueInfo, exchangeInfo ExchangeInfo) *Consumer {
	c := Consumer{
		uri:          uri,
		retry:        retry,
		bindKey:      bindKey,
		rabbit:       rabbit{uri: uri, retry: retry},
		connected:    false,
		queueInfo:    queueInfo,
		exchangeInfo: exchangeInfo,
	}
	return &c
}
func (c *Consumer) connect() (err error) {
	connection, err := amqp.Dial(c.rabbit.uri)
	if err != nil {
		logrus.Error("err in dial rabbitmq, err = " + err.Error())
		return err
	}
	c.rabbit.connection = connection
	channel, err := connection.Channel()
	if err != nil {
		logrus.Error("err in get channel, err = " + err.Error())
		return err
	}
	_, err = channel.QueueDeclare(c.queueInfo.Name, c.queueInfo.durable, c.queueInfo.autoDelete, c.queueInfo.exclusive, c.queueInfo.nowait, nil)
	if err != nil {
		logrus.Error("error in declare queue, err = " + err.Error())
		return err
	}
	err = channel.QueueBind(c.queueInfo.Name, c.bindKey, c.exchangeInfo.Name, c.queueInfo.nowait, nil)
	if err != nil {
		logrus.Error("error in queue bind, err = " + err.Error())
		return
	}
	closeSigChan := channel.NotifyClose(make(chan *amqp.Error))
	c.rabbit.errSigChan = closeSigChan
	c.rabbit.channel = channel
	c.connected = true
	return nil
}

func (c *Consumer) RunAsPushMod(prefetchCount int) (taskin chan amqp.Delivery, err error) {
	taskin = make(chan amqp.Delivery)
	for i := 0; i < c.rabbit.retry; i++ {
		c.Clean()
		err = c.connect()
		if err == nil {
			c.connected = true
			break
		}
	}
	if !c.connected {
		logrus.Panic("consumer cant not connect to mq, err =" + err.Error())
		return nil, err
	}
	c.rabbit.channel.Qos(prefetchCount, -1, false)
	deliveries, err := c.rabbit.channel.Consume(c.queueInfo.Name, c.bindKey, false, false, false, true, nil)
	if err != nil {
		logrus.Panic("error in get deliveries, err = " + err.Error())
		return nil, err
	}
	go func() {
		for {
			select {
			case v := <-deliveries:
				taskin <- v
			case er := <-c.rabbit.errSigChan:
				logrus.Error("an err from notifyClose, err = " + er.Error())
				var err error
				for i := 0; i < c.retry; i++ {
					err = c.reconnect()
					if err != nil {
						logrus.Panic("error in reconnect to rabbitmq, err = " + err.Error())
						continue
					}
					deliveries, err = c.rabbit.channel.Consume(c.queueInfo.Name, c.bindKey, false, false, false, true, nil)
					if err != nil {
						logrus.Error("error in get deliveries, err = " + err.Error())
					}
				}
				if err != nil {
					logrus.Panic("error in reconnect to rabbitmq, err = " + err.Error())
				}
			}
		}
	}()
	return taskin, err
}
func (c *Consumer) Clean() {
	if c.rabbit.channel != nil {
		err := c.rabbit.channel.Close()
		if err != nil {
			logrus.Error("error in close channel of consumer, err = " + err.Error())
		}
	}
	if c.rabbit.connection != nil {
		if !c.rabbit.connection.IsClosed() {
			err := c.rabbit.connection.Close()
			if err != nil {
				logrus.Error("error in close connection of consumer, err = " + err.Error())
			}
		}
	}
	if c.rabbit.confirms != nil {
		close(c.rabbit.confirms)
	}
	if c.rabbit.errSigChan != nil {
		close(c.rabbit.errSigChan)
	}
}

func (c *Consumer) RunAsPullMod() (taskin chan amqp.Delivery, err error) {
	taskin = make(chan amqp.Delivery)
	for i := 0; i < c.rabbit.retry; i++ {
		c.Clean()
		err = c.connect()
		if err == nil {
			c.connected = true
			break
		}
	}
	if !c.connected {
		logrus.Panic("error after connect many times, err = " + err.Error())
		return nil, err
	}
	go func() {
		for {
			delivery, _, err := c.rabbit.channel.Get(c.queueInfo.Name, true)
			if err == nil {
				taskin <- delivery
			} else {
				logrus.Error("error in get delivery,trying to reconnect, err = " + err.Error())
				err := c.reconnect()
				if err != nil {
					logrus.Panic("error in reconnect,err = " + err.Error())
				}
			}
		}
	}()
	return taskin, err
}

func (c *Consumer) reconnect() error {
	c.connected = false
	var err error
	for i := 0; i < c.rabbit.retry; i++ {
		c.Clean()
		err = c.connect()
		if err == nil {
			c.connected = true
			return nil
		}
	}
	return err
}
