package chat

import (
	"testing"
	"time"
)

func TestHubRegisterclient(t *testing.T) {
	testHub := newHub()

	testClient := &client{
		hub:  testHub,
		send: make(chan []byte, 1),
	}

	go testHub.run()

	testHub.register <- testClient

	testMessage := []byte("hello cabin chat!")
	testHub.broadcast <- testMessage

	select {
	case got := <-testClient.send:
		if string(got) != string(testMessage) {
			t.Fatalf("expected to receive %q, got %q", string(testMessage), string(got))
		}
	case <-time.After(time.Second):
		t.Fatal("expected registered client to receive broadcast message")
	}
}

func TestHubUnregisterClient(t *testing.T) {
	testHub := newHub()

	client := &client{
		hub:  testHub,
		send: make(chan []byte, 1),
	}

	go testHub.run()

	testHub.register <- client

	testHub.broadcast <- []byte("first message")

	select {
	case <-client.send:
	case <-time.After(time.Second):
		t.Fatal("expected registered client to receive first broadcast")
	}

	testHub.unregister <- client

	testHub.broadcast <- []byte("second message")

	select {
	case _, ok := <-client.send:
		if ok {
			t.Fatal("expected send channel to be closed after unregister")
		}
	case <-time.After(time.Second):
		t.Fatal("expected client send channel to be closed after unregister")
	}
}

func TestHubDropsSlowClientDuringBroadcast(t *testing.T) {
	testHub := newHub()

	testClient := &client{
		hub:  testHub,
		send: make(chan []byte, 1),
	}

	go testHub.run()

	testHub.register <- testClient

	testHub.broadcast <- []byte("first broadcast")
	testHub.broadcast <- []byte("second broadcast")

	select {
	case gotMessage := <-testClient.send:
		if string(gotMessage) != "first broadcast" {
			t.Fatalf("expected first broadcast, got %q", string(gotMessage))
		}
	case <-time.After(time.Second):
		t.Fatal("expected to receive first broadcast")
	}

	select {
	case _, ok := <-testClient.send:
		if ok {
			t.Fatal("expected the send channel to be closed after dropping a slow client")
		}
	case <-time.After(time.Second):
		t.Fatal("expected the client send channel to be closed after slow client was dropped")
	}
}

func TestHubDropsClosedClientDuringBroadcast(t *testing.T) {
	testHub := newHub()

	liveClient := &client{
		hub:  testHub,
		send: make(chan []byte, 1),
	}

	closedClient := &client{
		hub:  testHub,
		send: make(chan []byte, 1),
	}

	go testHub.run()

	testHub.register <- liveClient
	testHub.register <- closedClient

	close(closedClient.send)

	testHub.broadcast <- []byte("first broadcast")
	testHub.broadcast <- []byte("second broadcast")

	for _, expected := range []string{"first broadcast", "second broadcast"} {
		select {
		case gotMessage := <-liveClient.send:
			if string(gotMessage) != expected {
				t.Fatalf("expected %q, got %q", expected, string(gotMessage))
			}
		case <-time.After(time.Second):
			t.Fatal("expected to receive broadcast message")
		}
	}
}

func TestHubBroadcastsToConnectedClients(t *testing.T) {
	testHub := newHub()

	firstClient := &client{
		hub:  testHub,
		send: make(chan []byte, 1),
	}

	secondClient := &client{
		hub:  testHub,
		send: make(chan []byte, 1),
	}

	go testHub.run()

	testHub.register <- firstClient
	testHub.register <- secondClient

	testMessage := []byte("hello cabin chat!")
	testHub.broadcast <- testMessage

	for _, testClient := range []*client{firstClient, secondClient} {
		select {
		case got := <-testClient.send:
			if string(got) != string(testMessage) {
				t.Fatalf("expected to receive broadcast %q, got %q", string(testMessage), string(got))
			}
		case <-time.After(time.Second):
			t.Fatal("expected to receive broadcast message")
		}
	}
}
