package queue

import (
	"encoding/json"
	"testing"

	"go.uber.org/zap"
)

func TestQueue(t *testing.T) {

	// setup queue connections
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger, error: %v", err)
	}
	pubClient, err := NewQueueClient(logger, &QConfig{Host: "127.0.0.1", Port: 5672, Password: "guest", User: "guest", QueueName: "test-queue"})
	if err != nil {
		t.Fatalf("failed to connect to db, error: %v", err)
	}

	conClient, err := NewQueueClient(logger, &QConfig{Host: "127.0.0.1", Port: 5672, Password: "guest", User: "guest", QueueName: "test-queue"})
	if err != nil {
		t.Fatalf("failed to connect to db, error: %v", err)
	}
	defer pubClient.Close()
	defer conClient.Close()

	messages := []map[string]string{
		{
			"contract": "0x0contract1",
			"counter":  "1",
		},
		{
			"contract": "0x0contract2",
			"counter":  "1",
		},
		{
			"contract": "0x0contract3",
			"counter":  "2",
		},
	}

	t.Run("create channels", func(t1 *testing.T) {
		if _, err = pubClient.Channel(); err != nil {
			t1.Fatalf("failed to create publisher channel, err:%v", err)
		}

		if _, err = conClient.Channel(); err != nil {
			t1.Fatalf("failed to create consumer channel, err:%v", err)
		}
	})

	t.Run("setup queues", func(t1 *testing.T) {
		if _, err = pubClient.DeclareExchangeQueue(); err != nil {
			t1.Fatalf("failed to declare publisher queue, err:%v", err)
		}

		if _, err = conClient.DeclareExchangeQueue(); err != nil {
			t1.Fatalf("failed to declare consuner queue, err:%v", err)
		}
	})

	t.Run("publish messages", func(t1 *testing.T) {
		_, err := pubClient.SetupPublisherConfirms()
		if err != nil {
			t1.Fatalf("failed to create pub confirms channel, err:%v", err)
		}
		for i, msg := range messages {
			data, err := json.Marshal(msg)
			if err != nil {
				t1.Fatalf("failed to marshal msg, err:%v", err)
			}
			if err = pubClient.Publish(string(data), 10000); err != nil {
				t1.Fatalf("failed to publish to channel, err:%v", err)
			}
			if err = pubClient.ConfirmPublish(); err != nil {
				t1.Fatalf("%v for message %v", err.Error(), i+1)
			}
		}
	})

	t.Run("consume messages", func(t1 *testing.T) {
		if err = conClient.Qos(1); err != nil {
			t1.Fatalf("failed to setup qos for consumer, err:%v", err)
		}
		msgs, err := conClient.Consume()

		if err != nil {
			t1.Fatalf("failed get msgs for consumer, err:%v", err)
		}

		for i := 0; i < len(messages); i++ {
			msg, ok := <-msgs
			if !ok {
				t1.Fatalf("consumer channel closed!")
			}
			data := map[string]string{}
			if err = json.Unmarshal(msg.Body, &data); err != nil {
				t1.Fatalf("failed to unmarshal msg, err:%v", err)
			}
			for item := range messages[i] {
				if data[item] != messages[i][item] {
					t1.Fatalf("data mismatch expected=%v got=%v", messages[i][item], data[item])
				}
			}
			err := conClient.AckDelivery(&msg, false)
			if err != nil {
				t1.Fatalf("failed to ack message, err:%v", err)
			}
		}
	})
}
