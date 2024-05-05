package messaging

import (
	"github.com/dherik/ddd-golang-project/internal/domain"
)

type TaskMessagingHandler struct {
	taskRepository domain.TaskRepository
}

func NewTaskMessagingHandler(taskRepository domain.TaskRepository) {

}

func ReadNewTasks() {
	// consumer, err := rabbitmq.NewConsumer(
	// 	conn,
	// 	"my_queue",
	// 	rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
	// 	rabbitmq.WithConsumerOptionsExchangeName("events"),
	// 	rabbitmq.WithConsumerOptionsExchangeDeclare,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer consumer.Close()
}
