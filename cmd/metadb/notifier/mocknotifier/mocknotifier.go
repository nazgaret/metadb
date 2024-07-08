package mocknotifier

import (
	"context"
)

func NewMock() *mockNotifier {
	return &mockNotifier{}
}

type mockNotifier struct {
}

func (n *mockNotifier) Notify(ctx context.Context, message string) error {
	return nil
}
