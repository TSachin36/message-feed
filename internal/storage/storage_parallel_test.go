package storage

import (
	"context"
	"fmt"
	"message-feed/internal/models"
	"os"
	"sync"
	"testing"
)

func TestSaveMessageParallel(t *testing.T) {

	t.Parallel()

	filename := "parallel_test.txt"

	defer os.Remove(filename)

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {

		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			msg := models.Message{
				UserID: fmt.Sprintf("User%d", i),
				Text:   fmt.Sprintf("Message%d", i),
			}

			err := SaveMessage(
				context.Background(),
				filename,
				msg,
			)

			if err != nil {
				t.Error(err)
			}

		}(i)
	}

	wg.Wait()

	messages, err := GetLast10Messages(
		context.Background(),
		filename,
	)

	if err != nil {
		t.Fatal(err)
	}

	if len(messages) != 10 {
		t.Errorf(
			"expected 10 messages, got %d",
			len(messages),
		)
	}
}
