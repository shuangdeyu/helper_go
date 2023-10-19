package test

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"testing"
)

var (
	uri          = flag.String("uri", "amqp://test:test@192.168.2.110:5672/test", "AMQP URI")
	exchange     = flag.String("exchange", "test.exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "test.queue", "Ephemeral AMQP queue name")
	bindingKey   = flag.String("key", "normal", "AMQP binding key")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
)

func init() {
	flag.Parse()
}

func TestMqConsumer(t *testing.T) {
	_, err := MqConsumer(*uri, *exchange, *exchangeType, *queue, *bindingKey, *consumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func MqConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
	for {
		c := &Consumer{
			conn:    nil,
			channel: nil,
			tag:     ctag,
			done:    make(chan error),
		}

		var err error
		var notify chan *amqp.Error = make(chan *amqp.Error)

		// 连接mq服务器
		c.conn, err = amqp.Dial(amqpURI)
		if err != nil {
			return nil, fmt.Errorf("Dial: %s", err)
		}
		c.conn.NotifyClose(notify)

		// 获取通道
		c.channel, err = c.conn.Channel()
		if err != nil {
			return nil, fmt.Errorf("Channel: %s", err)
		}

		// 声明交换机
		if err = c.channel.ExchangeDeclare(
			exchange,     // name of the exchange
			exchangeType, // type
			true,         // durable
			false,        // delete when complete
			false,        // internal
			false,        // noWait
			nil,          // arguments
		); err != nil {
			return nil, fmt.Errorf("Exchange Declare: %s", err)
		}

		// 声明队列
		queue, err := c.channel.QueueDeclare(
			queueName, // name of the queue
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // noWait
			nil,       // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("Queue Declare: %s", err)
		}

		//log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		//	queue.Name, queue.Messages, queue.Consumers, key)

		// 绑定
		if err = c.channel.QueueBind(
			queue.Name, // name of the queue
			key,        // bindingKey
			exchange,   // sourceExchange
			false,      // noWait
			nil,        // arguments
		); err != nil {
			return nil, fmt.Errorf("Queue Bind: %s", err)
		}

		// 开始消费消息
		deliveries, err := c.channel.Consume(
			queue.Name, // name
			c.tag,      // consumerTag,
			false,      // 自动发送ACK
			false,      // exclusive
			false,      // noLocal
			false,      // noWait
			nil,        // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("Queue Consume: %s", err)
		}

		// 消息处理，出错的话跳出消息消费，重新连接，记录出错次数和接收次数
		for {
			select {
			case notifyErr := <-notify:
				//common.MetricsInc("mq.total.error", int64(1))
				log.Printf("AMQP close: %s", notifyErr.Error())
				goto reconnect
			case d := <-deliveries:
				//common.MetricsInc("mq.total.receive", int64(1))
				go handle(d, c.done)
			case <-c.done:
				//common.MetricsInc("mq.total.error", int64(1))
				goto reconnect
			}
		}

	reconnect:
		log.Println("reconnection")
		// 重连超出预定次数，则断开连接
		//if {
		//	return c, nil
		//}
	}
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries amqp.Delivery, done chan error) {
	log.Printf(
		"got %dB delivery: [%v] %q",
		len(deliveries.Body),
		deliveries.DeliveryTag,
		deliveries.Body,
	)
	deliveries.Ack(false) // 处理成功，发送ACK信息
	//done <- nil // 关闭通道
}
