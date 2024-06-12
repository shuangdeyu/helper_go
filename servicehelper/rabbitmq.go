package servicehelper

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"

	"github.com/streadway/amqp"

	"github.com/shuangdeyu/helper_go/comhelper"
)

// 主体参数
type MqArgs struct {
	Id       string
	User_id  string
	Event    string
	App_id   int
	Ip       string
	Dateline string
	Endpoint string
	Agent    string
	Content  string
}

// 发送event参数
func SendEventParam(user_id string, event string, event_params interface{}) string {
	content, _ := json.Marshal(event_params)
	content_str := string(content)
	mq_params := &MqArgs{
		Id:      comhelper.Get_uuid(4),
		User_id: user_id,
		Event:   event,
		Content: content_str,
	}
	mq, _ := json.Marshal(mq_params)
	return string(mq)
}

/**
 * 发布消息
 * @param uri 			string 		mq地址
 * @param exchange		string 		交换机名称
 * @param queue		 	string 		队列名称
 * @param routing_key	string 		路由键名
 * @param queue_args 	amqp.Table	队列参数
 * @param expiration 	string 		消息过期时间,单位毫秒
 * @param message 		string 		消息内容
 */
func MqPublish(uri, exchange, queue, routing_key string, queue_args amqp.Table, expiration string, message string) error {

	// 建立连接
	connection, err := amqp.Dial(uri)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err.Error())
		return err
	}
	defer connection.Close()

	// 创建一个Channel
	channel, err := connection.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err.Error())
		return err
	}
	defer channel.Close()

	// 声明exchange(可以在控制台创建好)
	if err := channel.ExchangeDeclare(
		exchange, //name
		"direct", //exchangeType
		true,     //durable
		false,    //auto-deleted
		false,    //internal
		false,    //noWait
		nil,      //arguments
	); err != nil {
		log.Println("Failed to declare a exchange:", err.Error())
		return err
	}

	// 声明一个queue
	if _, err := channel.QueueDeclare(
		queue,      // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		queue_args, // arguments
	); err != nil {
		log.Println("Failed to declare a queue:", err.Error())
		return err
	}

	// exchange 绑定 queue(同样可以在控制台绑定好)
	channel.QueueBind(queue, routing_key, exchange, false, nil)

	// 发送
	if err = channel.Publish(
		exchange,    // exchange
		routing_key, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(message),
			Expiration:      expiration, // 消息过期时间
		},
	); err != nil {
		log.Println("Failed to publish a message:", err.Error())
		return err
	}

	return nil
}

/**
 * mq消息相关类型获取转换
 */
func ParseEventContent(content string) map[string]interface{} {
	ret := make(map[string]interface{})
	if err := json.Unmarshal([]byte(content), &ret); err != nil {

	}
	return ret
}

func GetStringFromObj(obj map[string]interface{}, key string, defaultValue string) string {
	value := defaultValue
	ret, ok := obj[key]
	if ok {
		value, ok = ret.(string)
		if !ok {
			switch ret.(type) {
			case int:
				value = strconv.Itoa(ret.(int))
			case int64:
				value = strconv.FormatInt(ret.(int64), 10)
			case float32:
				value = strconv.FormatFloat(float64(ret.(float32)), 'f', -1, 64)
			case float64:
				value = strconv.FormatFloat(ret.(float64), 'f', -1, 64)
			case bool:
				value = strconv.FormatBool(ret.(bool))
			default:
				return defaultValue
			}
		}
	}
	return value
}

// 从结构中获取类型值，如果类型不匹配，则返回默认值,
// 请注意默认值的类型，以方便后续代码处理
// 不支持的类型：匿名的struct,map[xxx]xxx类型,指针类型
func GetTypeDataFromObj(obj map[string]interface{}, key string, vType string, defaultValue interface{}) interface{} {
	ret, ok := obj[key]
	if ok {
		t := reflect.TypeOf(ret)
		name := t.Name()
		if name == vType {
			return ret
		}
	}
	return defaultValue
}
