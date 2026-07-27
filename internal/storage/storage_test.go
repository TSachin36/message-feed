package storage

import (
	"context"
	"message-feed/internal/models"
	"os"
	"path/filepath"
	"testing"
)

func TestGetLast10Messages(t *testing.T) {

	filename := "test_messages.txt"

	content := `Alice: Hello
Bob: Good Morning
Charlie: How are you?`

	err := os.WriteFile(
		filename,
		[]byte(content),
		0644,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(filename)

	ctx := context.WithValue(
		context.Background(),
		"traceID",
		"TEST-123",
	)

	messages, err := GetLast10Messages(ctx, filename)
	if err != nil {
		t.Fatal(err)
	}

	if len(messages) != 3 {
		t.Fatalf(
			"Expected 3 messages, got %d",
			len(messages),
		)
	}

	if messages[0].UserID != "Alice" || messages[0].Text != "Hello" {
		t.Errorf(
			"Expected Alice: Hello, got %+v",
			messages[0],
		)
	}

	if messages[1].UserID != "Bob" || messages[1].Text != "Good Morning" {
		t.Errorf(
			"Expected Bob: Good Morning, got %+v",
			messages[1],
		)
	}

	if messages[2].UserID != "Charlie" || messages[2].Text != "How are you?" {
		t.Errorf(
			"Expected Charlie: How are you?, got %+v",
			messages[2],
		)
	}
}

func TestSaveMessage(t *testing.T) {

	filename := "test_save.txt"

	defer os.Remove(filename)

	ctx := context.WithValue(
		context.Background(),
		"traceID",
		"TEST-456",
	)

	msg := models.Message{
		UserID: "Alice",
		Text:   "Hello World",
	}

	err := SaveMessage(ctx, filename, msg)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)

	expected := "Alice: Hello World\n"

	if result != expected {
		t.Errorf(
			"Expected %q, got %q",
			expected,
			result,
		)
	}
}

func TestGetLast10MessagesForUser(t *testing.T) {

	filename := filepath.Join(
		t.TempDir(),
		"messages.txt",
	)

	data := `alice: message 1
bob: message 1
alice: message 2
bob: message 2
alice: message 3
`

	err := os.WriteFile(
		filename,
		[]byte(data),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	messages, err := GetLast10MessagesForUser(
		context.Background(),
		filename,
		"alice",
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(messages) != 3 {
		t.Fatalf(
			"expected 3 messages, got %d",
			len(messages),
		)
	}

	for _, msg := range messages {
		if msg.UserID != "alice" {
			t.Fatalf(
				"expected only alice messages, got user %q",
				msg.UserID,
			)
		}
	}
}
