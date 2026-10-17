package setup

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/testcontainers/testcontainers-go"
)

type Exchange struct {
	testcontainers.ExecOptions
	Name       string
	VHost      string
	Type       string
	AutoDelete bool
	Internal   bool
	Durable    bool
	Args       map[string]interface{}
}

func (e Exchange) AsCommand() []string {
	cmd := []string{"rabbitmqadmin"}

	if e.VHost != "" {
		cmd = append(cmd, "--vhost="+e.VHost)
	}

	cmd = append(cmd, "declare", "exchange", fmt.Sprintf("name=%s", e.Name), fmt.Sprintf("type=%s", e.Type))

	if e.AutoDelete {
		cmd = append(cmd, "auto_delete=true")
	}
	if e.Internal {
		cmd = append(cmd, "internal=true")
	}
	if e.Durable {
		cmd = append(cmd, fmt.Sprintf("durable=%t", e.Durable))
	}

	if len(e.Args) > 0 {
		bytes, err := json.Marshal(e.Args)
		if err != nil {
			return cmd
		}

		cmd = append(cmd, "arguments="+string(bytes))
	}

	return cmd
}

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

func readExecOutput(reader io.Reader) string {
	out, err := io.ReadAll(reader)
	if err != nil {
		return ""
	}
	return string(out)
}
