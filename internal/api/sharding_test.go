package api

import (
	"fmt"
	"testing"
)

func TestHashRingHasVirtualNodes(t *testing.T) {

	buildHashRing()

	expected := len(shardAddresses) * virtualNodesPerShard

	if len(hashRing) != expected {
		t.Fatalf(
			"expected %d hash ring positions, got %d",
			expected,
			len(hashRing),
		)
	}
}

func TestShardForUserIsDeterministic(t *testing.T) {

	buildHashRing()

	userID := "alice"

	first := shardForUser(userID)

	for i := 0; i < 100; i++ {

		got := shardForUser(userID)

		if got != first {
			t.Fatalf(
				"expected %q, got %q",
				first,
				got,
			)
		}
	}
}

func TestVirtualNodesDistributeUsersAcrossShards(t *testing.T) {

	buildHashRing()

	counts := make(map[string]int)

	const users = 10000

	for i := 0; i < users; i++ {

		userID := fmt.Sprintf(
			"user-%d",
			i,
		)

		shard := shardForUser(userID)

		counts[shard]++
	}

	for _, address := range shardAddresses {

		if counts[address] == 0 {
			t.Fatalf(
				"shard %q received no users",
				address,
			)
		}

		t.Logf(
			"%s -> %d users",
			address,
			counts[address],
		)
	}
}

func TestAddingShardMovesSubsetOfUsers(t *testing.T) {

	originalAddresses := append(
		[]string(nil),
		shardAddresses...,
	)

	defer func() {
		shardAddresses = originalAddresses
		buildHashRing()
	}()

	const users = 10000

	buildHashRing()

	before := make(map[string]string, users)

	for i := 0; i < users; i++ {

		userID := fmt.Sprintf(
			"user-%d",
			i,
		)

		before[userID] = shardForUser(userID)
	}

	shardAddresses = append(
		shardAddresses,
		"localhost:50054",
	)

	buildHashRing()

	moved := 0
	movedToNewShard := 0

	for i := 0; i < users; i++ {

		userID := fmt.Sprintf(
			"user-%d",
			i,
		)

		after := shardForUser(userID)

		if after != before[userID] {
			moved++

			if after == "localhost:50054" {
				movedToNewShard++
			}
		}
	}

	t.Logf(
		"moved %d of %d users (%.2f%%)",
		moved,
		users,
		float64(moved)/float64(users)*100,
	)

	if moved == 0 {
		t.Fatal("expected some users to move after adding a shard")
	}

	if moved >= users/2 {
		t.Fatalf(
			"too many users moved: %d of %d",
			moved,
			users,
		)
	}

	if moved != movedToNewShard {
		t.Fatalf(
			"expected all moved users to move to new shard: moved=%d, movedToNewShard=%d",
			moved,
			movedToNewShard,
		)
	}
}

func TestRemovingShardMovesOnlyAffectedUsers(t *testing.T) {

	originalAddresses := append(
		[]string(nil),
		shardAddresses...,
	)

	defer func() {
		shardAddresses = originalAddresses
		buildHashRing()
	}()

	const users = 10000

	buildHashRing()

	before := make(map[string]string, users)

	for i := 0; i < users; i++ {

		userID := fmt.Sprintf(
			"user-%d",
			i,
		)

		before[userID] = shardForUser(userID)
	}

	removedShard := "localhost:50052"

	shardAddresses = []string{
		"localhost:50051",
		"localhost:50053",
	}

	buildHashRing()

	moved := 0
	originallyOnRemovedShard := 0

	for i := 0; i < users; i++ {

		userID := fmt.Sprintf(
			"user-%d",
			i,
		)

		oldShard := before[userID]
		newShard := shardForUser(userID)

		if oldShard == removedShard {
			originallyOnRemovedShard++
		}

		if oldShard != newShard {

			moved++

			if oldShard != removedShard {
				t.Fatalf(
					"user %q moved from %q to %q even though its shard was not removed",
					userID,
					oldShard,
					newShard,
				)
			}
		}
	}

	t.Logf(
		"removed shard held %d users; %d users moved (%.2f%%)",
		originallyOnRemovedShard,
		moved,
		float64(moved)/float64(users)*100,
	)

	if moved != originallyOnRemovedShard {
		t.Fatalf(
			"expected exactly %d affected users to move, got %d",
			originallyOnRemovedShard,
			moved,
		)
	}
}
