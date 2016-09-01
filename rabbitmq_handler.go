package common

import (
	"sync"
	"runtime/debug"
	"third/amqp"
	"time"
)

type AMQPReceipt struct {
	uris     []string
	alive    bool
	i        int
	channel  *amqp.Channel
	openMu  sync.Mutex
	delivery *amqp.Delivery
}

func (r *AMQPReceipt) Ack() {
	r.delivery.Ack(false)
}

func (r *AMQPReceipt) Requeue() {
	r.delivery.Reject(true)
}

func (r *AMQPReceipt) Reject() {
	r.delivery.Reject(false)
}

//type AMQPDriver struct {
//	uris    []string
//	alive   bool
//	i       int
//	channel *amqp.Channel
//}

func (r *AMQPReceipt) Connect(uri string, queue string, durable bool) (err error) {
	r.openMu.Lock()
	defer r.openMu.Unlock()
	if r.alive {
		return
	}

	conn, err := amqp.Dial(uri)
	if err != nil {
		Errorf("amqp.dial %s,%s error :%v", uri, queue, err)
		return
	}

	r.channel, err = conn.Channel()
	if err != nil {
		Errorf("conn.Channel %s,%s error :%v", uri, queue, err)
		return
	}
	r.alive = true

	var arguments amqp.Table
	_, err = r.channel.QueueDeclare(
		queue,
		durable,
		false,
		false, // exclusive
		false, // noWait
		arguments,
	)
	amqpMessagePool.New = func() interface{} {
		var message AmqpMessage
		message.Receipt = r
		//Infof("amqpMessagePool get new")
		return &message
	}
	Noticef("amqp connect %s,%s success ", uri, queue)
	return
}

func (r *AMQPReceipt) IsConnected() bool {
	return r.alive
}

func (r *AMQPReceipt) Publish(queue string, payload []byte) (err error) {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         payload,
	}
	err = r.channel.Publish(
		"",
		queue,
		false,
		false,
		msg,
	)
	if err != nil {
		Errorf("amqp.Publish %s,%s error :%v", amqp_client_uri, queue, err)
	}

	return err
}

var AmqpReceipt AMQPReceipt
var amqp_client_uri string
var amqp_client_queue string

func InitRabbitmqClient(uri, queue string) error {
	amqp_client_uri = uri
	amqp_client_queue = queue
	return AmqpReceipt.Connect(uri, queue, false)
}

func ReconnAmqpClient() error {
	return AmqpReceipt.Connect(amqp_client_uri, amqp_client_queue, false)
}

type AmqpMessage struct {
	MessageId     int64
	RetryTime     int
	Acked         bool
	ConnMessageId uint
	LastSendTime  int64
	ContentType   string
	Body          []byte
	Receipt       *AMQPReceipt
}

var amqpMessagePool sync.Pool

func NewAmqpMessage(r *AMQPReceipt) *AmqpMessage {
	var message AmqpMessage
	message.Receipt = r

	return &message
}

func PutMessage(message *AmqpMessage) {
	//message.Receipt.Ack()

		if nil != message {
			amqpMessagePool.Put(message)
		} else {
		Errorf("put nil message")
		Infof(string(debug.Stack()))
		}
}

func (r *AMQPReceipt) GetMessages(queue_name string, rate int) (<-chan *AmqpMessage, error) {

	r.channel.Qos(rate, 0, false)
	deliveries, err := r.channel.Consume(
		queue_name,
		"",    // consumerTag
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		r.alive = false
		return nil, err
	}

	amqpMessages := make(chan *AmqpMessage, rate)
	go func() {
		for d := range deliveries {
			Infof("amqp recv message")

				message := amqpMessagePool.Get()
				switch message.(type) {
				case *AmqpMessage:
					message.(*AmqpMessage).Body = d.Body
					message.(*AmqpMessage).ContentType = d.ContentType
					message.(*AmqpMessage).Receipt.delivery = &d
					message.(*AmqpMessage).Acked = false
					message.(*AmqpMessage).ConnMessageId = 0
					message.(*AmqpMessage).LastSendTime = 0
					message.(*AmqpMessage).MessageId = 0
					message.(*AmqpMessage).RetryTime = 0
					amqpMessages <- message.(*AmqpMessage)
				message.(*AmqpMessage).Receipt.Ack()
				default:
				Errorf("pool type error :%T", message)
				}
			/*

				message := NewAmqpMessage(c)
			message.Body = d.Body
			message.ContentType = d.ContentType
			message.Receipt.delivery = &d
			message.Acked = false
			message.ConnMessageId = 0
			message.LastSendTime = 0
			message.MessageId = 0
			message.RetryTime = 0
			amqpMessages <- message
			message.Receipt.Ack()
			*/
		}
		// connection was lost
		Errorf("!!!!!!!!amqp quit, close channel!!!!!!!!!!!!!!!")
		r.alive = false
		close(amqpMessages)
	}()
	return amqpMessages, nil
}
