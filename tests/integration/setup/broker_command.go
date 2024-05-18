package setup

import (
	"encoding/json"
	"fmt"

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
