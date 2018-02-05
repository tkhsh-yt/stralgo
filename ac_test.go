package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddWord(t *testing.T) {
	assert := assert.New(t)

	trie := CreateTrie()
	words := []string{"a", "ab", "aca", "ba", "c", "cab"}
	accIDs := trie.AddWords(words)

	assert.Equal(10, trie.StatesCount)

	ids := []int{1, 2, 4, 6, 7, 9}
	for i, id := range ids {
		assert.Equal(i, accIDs[id])
	}
}

func TestMachWord(t *testing.T) {
	assert := assert.New(t)

	trie := CreateTrie()
	words := []string{"a", "ab", "aca", "ba", "c", "cab"}
	trie.AddWords(words)

	pma := trie.ToPMA()
	matches := pma.MatchWord("caba")

	assert.Equal(6, len(matches))
	assert.Contains(matches, Match{7, 0})
	assert.Contains(matches, Match{1, 1})
	assert.Contains(matches, Match{9, 2})
	assert.Contains(matches, Match{2, 2})
	assert.Contains(matches, Match{1, 3})
	assert.Contains(matches, Match{6, 3})

	matches = pma.MatchWord("acab")

	assert.Equal(6, len(matches))
	assert.Contains(matches, Match{1, 0})
	assert.Contains(matches, Match{7, 1})
	assert.Contains(matches, Match{4, 2})
	assert.Contains(matches, Match{1, 2})
	assert.Contains(matches, Match{9, 3})
	assert.Contains(matches, Match{2, 3})

	matches = pma.MatchWord("acac")

	assert.Equal(5, len(matches))
	assert.Contains(matches, Match{1, 0})
	assert.Contains(matches, Match{7, 1})
	assert.Contains(matches, Match{1, 2})
	assert.Contains(matches, Match{4, 2})
	assert.Contains(matches, Match{7, 3})
}
