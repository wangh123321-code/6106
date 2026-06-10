package algorithm

import (
	"math/rand"
	"sort"
)

type SeedEntry struct {
	PlayerID string
	Ranking  int
}

func SnakeSeed(players []SeedEntry, bracketSize int) [][]string {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Ranking < players[j].Ranking
	})

	slots := make([]string, bracketSize)
	for i := 0; i < bracketSize; i++ {
		slots[i] = ""
	}

	seedCount := len(players)
	if seedCount > bracketSize {
		seedCount = bracketSize
	}

	for i := 0; i < seedCount; i++ {
		pos := snakePosition(i, bracketSize)
		slots[pos] = players[i].PlayerID
	}

	rounds := buildRounds(slots, bracketSize)
	return rounds
}

func snakePosition(seedIndex, bracketSize int) int {
	half := bracketSize / 2

	seedSlots := []int{}

	seedSlots = append(seedSlots, 0, half-1)

	round := 1
	for len(seedSlots) < bracketSize {
		chunk := []int{}
		for i := 0; i < (1 << round); i++ {
			pos := seedSlots[i] / 2
			chunk = append(chunk, pos, half-1-pos)
		}
		seedSlots = chunk
		round++
		if round > 10 {
			break
		}
	}

	if seedIndex < len(seedSlots) {
		return seedSlots[seedIndex]
	}
	return seedIndex
}

func RandomDraw(playerIDs []string, bracketSize int) [][]string {
	shuffled := make([]string, len(playerIDs))
	copy(shuffled, playerIDs)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	slots := make([]string, bracketSize)
	for i := 0; i < bracketSize; i++ {
		slots[i] = ""
	}
	for i, pid := range shuffled {
		if i < bracketSize {
			slots[i] = pid
		}
	}

	rounds := buildRounds(slots, bracketSize)
	return rounds
}

func buildRounds(firstRoundSlots []string, bracketSize int) [][]string {
	numRounds := 0
	for size := bracketSize; size > 1; size /= 2 {
		numRounds++
	}

	rounds := make([][]string, numRounds)
	rounds[0] = firstRoundSlots

	for r := 1; r < numRounds; r++ {
		matchCount := bracketSize / (1 << (r + 1))
		rounds[r] = make([]string, matchCount*2)
		for i := range rounds[r] {
			rounds[r][i] = ""
		}
	}

	return rounds
}

func NextPowerOf2(n int) int {
	p := 1
	for p < n {
		p *= 2
	}
	return p
}
