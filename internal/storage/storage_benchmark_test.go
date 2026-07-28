package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"message-feed/internal/models"
)

func BenchmarkSaveMessage(b *testing.B) {

	filename := filepath.Join(
		b.TempDir(),
		"messages.txt",
	)

	ctx := context.Background()

	msg := models.Message{
		UserID: "alice",
		Text:   "Benchmark test message",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		err := SaveMessage(
			ctx,
			filename,
			msg,
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetLast10MessagesForUser(b *testing.B) {

	filename := filepath.Join(
		b.TempDir(),
		"messages.txt",
	)

	var data string

	for i := 0; i < 1000; i++ {

		user := "alice"

		if i%2 == 0 {
			user = "bob"
		}

		data += fmt.Sprintf(
			"%s: message %d\n",
			user,
			i,
		)
	}

	err := os.WriteFile(
		filename,
		[]byte(data),
		0644,
	)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		_, err := GetLast10MessagesForUser(
			ctx,
			filename,
			"alice",
		)
		if err != nil {
			b.Fatal(err)
		}
	}
}
