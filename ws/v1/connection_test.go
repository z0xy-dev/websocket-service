package v1_test

import (
	"testing"
	v1 "websocketservice/ws/v1"
)

// TestAddEvent ensures AddEvent adds an event only once.
func TestAddEvent(t *testing.T) {
	conn := &v1.Connection{}
	conn.AddEvent("test")
	if len(conn.SubscribedEvents) != 1 {
		t.Errorf("expected 1 event, got %d", len(conn.SubscribedEvents))
	}
	conn.AddEvent("test")
	if len(conn.SubscribedEvents) != 1 {
		t.Errorf("expected event not to be duplicated")
	}
}

// TestRemoveEvent ensures RemoveEvent actually removes the subscribed event.
func TestRemoveEvent(t *testing.T) {
	conn := &v1.Connection{
		SubscribedEvents: []string{"event1", "event2"},
	}
	conn.RemoveEvent("event1")
	if len(conn.SubscribedEvents) != 1 || conn.SubscribedEvents[0] != "event2" {
		t.Errorf("expected [event2], got %v", conn.SubscribedEvents)
	}
	conn.RemoveEvent("event2")
	if len(conn.SubscribedEvents) != 0 {
		t.Errorf("expected no events, got %v", conn.SubscribedEvents)
	}
}

// TestHasEvent checks HasEvent returns correct boolean.
func TestHasEvent(t *testing.T) {
	conn := &v1.Connection{
		SubscribedEvents: []string{"foo", "bar"},
	}
	if !conn.HasEvent("foo") {
		t.Errorf("expected true for existing event")
	}
	if conn.HasEvent("nope") {
		t.Errorf("expected false for unexisting event")
	}
}
