// Конекторы для получения соединения с очередями
//
// Пример реализации конектера к ampq
//
// type AMQPConnector struct {
// 	appCtx context.Context
// 	conn *amqp.Connection
// 	ch *amqp.Channel
// }
// func NewAMQPConnector(appCtx context.Context, connUrl string) (c *AMQPConnector, err error) {
// 	c = &AMQPConnector{appCtx: appCtx}
// 	c.conn, err = amqp.Dial(connUrl)
// 	if err != nil {
// 		return
// 	}
// 	c.ch, err = c.conn.Channel()
// 	if err != nil {
// 		c.conn.Close()
// 		return
// 	}
// 	context.AfterFunc(appCtx, func() {
// 		c.ch.Close()
// 		c.conn.Close()
// 	})
// 	return
// }
//
// func (c *AMQPConnector) Connect(ctx context.Context, queueID model.QueueID) (<-chan model.Message, func(context.Context, model.QueueID, model.Message) error, error) {
// 		q, err := c.ch.QueueDeclare(string(queueID), false, false, false, false, nil)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		msgs, err := c.ch.Consume(q.Name,"",true,false,false,false,nil)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		queue := make(chan model.Message)
// 		go func() {
// 			for msg := range msgs {
// 				queue <- msg.Body
// 			}
// 			close(queue)
// 		}()
// 		return queue, func(ctx context.Context, qi model.QueueID, msg model.Message) error {
// 			return c.ch.Publish("",q.Name,false,false,amqp.Publishing{Body:msg})
// 		}, nil
// }
//
package connector
