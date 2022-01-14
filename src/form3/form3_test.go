package form3

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Create two new clients and confirm that they are unique.
	c := NewClient()
	if got, want := c.BaseURL.String(), BaseURL(); got != want {
		t.Errorf("form3::NewClient BaseURL returned %v, want %v", got, want)
	}

	c2 := NewClient()
	if c.client == c2.client {
		t.Error("form3::NewClient returned the same http.Client, want different")
	}
}

func TestClientDoWithoutContext(t *testing.T) {
	// Perform an API request with a nil context and confirm that an error is produced.
	c := NewClient()
	var ctx context.Context
	_, err := c.Do(ctx, nil, nil)
	if err == nil {
		t.Error("form3::Do without context returned nil error")
	}
}
