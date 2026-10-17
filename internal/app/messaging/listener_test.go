package messaging

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeHandler struct {
	outcome error
	calls   int
}

func (f *fakeHandler) HandleMessage(ctx context.Context, msg Delivery) error {
	f.calls++
	return f.outcome
}

func TestProcessDeliveryAcksSuccess(t *testing.T) {
	handler := &fakeHandler{}
	msg := testDelivery(`{}`)

	processDelivery(context.Background(), handler, msg)

	require.Equal(t, 1, handler.calls)
	require.True(t, msg.acked, "success must be acknowledged")
	require.False(t, msg.nacked)
}

func TestProcessDeliveryDeadLettersPermanentOnFirstAttempt(t *testing.T) {
	handler := &fakeHandler{outcome: Permanent(errors.New("invalid fields"))}
	msg := testDelivery(`{}`)

	processDelivery(context.Background(), handler, msg)

	require.Equal(t, 1, handler.calls)
	require.False(t, msg.acked)
	require.True(t, msg.nacked, "a permanent failure must be nacked")
	require.False(t, msg.nackRequeue, "a permanent failure must not be requeued")
}

func TestProcessDeliveryDeadLettersTransient(t *testing.T) {
	handler := &fakeHandler{outcome: errors.New("database unavailable")}
	msg := testDelivery(`{}`)

	processDelivery(context.Background(), handler, msg)

	require.Equal(t, 1, handler.calls)
	require.False(t, msg.acked)
	require.True(t, msg.nacked, "an unconsumed event must be nacked")
	require.False(t, msg.nackRequeue, "nacking must dead-letter, not requeue")
}
