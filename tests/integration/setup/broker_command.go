package setup

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

// PublishMessage publishes a raw payload to an exchange through the broker's
// admin CLI, so tests can exercise consumers over real RabbitMQ.
func PublishMessage(exchange, routingKey, payload string) error {
	cmd := []string{"rabbitmqadmin", "publish",
		"exchange=" + exchange,
		"routing_key=" + routingKey,
		"payload=" + payload,
	}

	exitCode, output, err := brokerContainer.Exec(brokerCtx, cmd)
	if err != nil {
		return fmt.Errorf("failed publishing message: %w", err)
	}
	if exitCode != 0 {
		return fmt.Errorf("failed publishing message, exit code %d: %s", exitCode, readExecOutput(output))
	}
	return nil
}

// DeadLetteredPayloads returns the payloads currently parked in a queue, so
// tests can assert a failed event was dead-lettered. The broker CLI's get
// defaults to requeueing, so inspection does not drain the queue.
func DeadLetteredPayloads(queueName string) ([]string, error) {
	cmd := []string{"rabbitmqadmin", "-f", "pretty_json", "get", "queue=" + queueName}

	exitCode, output, err := brokerContainer.Exec(brokerCtx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed inspecting queue %s: %w", queueName, err)
	}
	if exitCode != 0 {
		return nil, fmt.Errorf("failed inspecting queue %s, exit code %d: %s", queueName, exitCode, readExecOutput(output))
	}

	var messages []struct {
		Payload string `json:"payload"`
	}
	if err := json.Unmarshal([]byte(readExecOutput(output)), &messages); err != nil {
		return nil, fmt.Errorf("failed parsing inspect output for queue %s: %w", queueName, err)
	}

	payloads := make([]string, 0, len(messages))
	for _, m := range messages {
		payloads = append(payloads, m.Payload)
	}
	return payloads, nil
}

// readExecOutput drains a container exec output stream. The docker exec
// stream frames stdout and stderr as repeated [type, 0, 0, 0, size] headers
// followed by the payload; this concatenates the payloads in order.
func readExecOutput(reader io.Reader) string {
	if reader == nil {
		return ""
	}

	var out bytes.Buffer
	header := make([]byte, 8)
	for {
		if _, err := io.ReadFull(reader, header); err != nil {
			break
		}
		size := binary.BigEndian.Uint32(header[4:8])
		if _, err := io.CopyN(&out, reader, int64(size)); err != nil {
			break
		}
	}
	return out.String()
}
