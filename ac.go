package main

import "fmt"

// State is a node of Trie or Pattern Matching Automaton.
type State struct {
	ID       int
	Accept   bool
	Success  map[rune]*State
	Failure  *State
	Matching []int
}

// Trie is a kind of search tree, an ordered tree structure.
type Trie struct {
	Init        *State
	StatesCount int
}

// CreateTrie returns an initial Trie.
func CreateTrie() *Trie {
	init := &State{
		Success:  map[rune]*State{},
		Matching: []int{},
	}
	return &Trie{init, 1}
}

// NewState returns new state having fresh ID.
func (tp *Trie) NewState() *State {
	s := &State{
		ID:      tp.StatesCount,
		Success: map[rune]*State{},
	}
	tp.StatesCount++
	return s
}

// Add adds an word as array of rune to Trie.
func (tp *Trie) Add(word []rune) int {
	s := tp.Init
	for _, c := range word {
		if _, ok := s.Success[c]; !ok {
			ns := tp.NewState()
			s.Success[c] = ns
		}
		s = s.Success[c]
	}

	s.Accept = true

	return s.ID
}

// AddWord adds an word as string to Trie.
func (tp *Trie) AddWord(word string) int {
	return tp.Add([]rune(word))
}

// AddWords adds words as array of string to Trie.
func (tp *Trie) AddWords(words []string) map[int]int {
	ids := make(map[int]int)

	for i, w := range words {
		ids[tp.AddWord(w)] = i
	}

	return ids
}

// Size returns the size of Trie.
func (tp *Trie) Size() int {
	q := newQueue()
	q.enqueue(tp.Init)
	size := 0
	for !q.isEmpty() {
		s := q.dequeue()
		for _, v := range s.Success {
			q.enqueue(v)
		}
		size++
	}
	return size
}

// PMA is Pattern Matching Automaton (?).
type PMA struct {
	*Trie
}

// ToPMA converts Trie to PMA.
func (tp *Trie) ToPMA() *PMA {
	pma := &PMA{Trie: tp}
	makeFailureLink(pma)

	return pma
}

// makeFilureLink adds links used when matching is falied.
func makeFailureLink(pma *PMA) {
	que := newQueue()

	pma.Init.Failure = pma.Init
	for _, s := range pma.Init.Success {
		s.Failure = pma.Init
		que.enqueue(s)
	}

	for !que.isEmpty() {
		s := que.dequeue()

		f := s.Failure
		if f.Accept {
			s.Matching = append(s.Matching, f.ID)
		}
		s.Matching = append(s.Matching, f.Matching...)

		for c, n := range s.Success {
			f := s.Failure

			if nn, ok := f.Success[c]; ok {
				n.Failure = nn
			} else {
				n.Failure = f.Failure
			}

			que.enqueue(n)
		}
	}
}

// Match represent matching.
type Match struct {
	ID  int // State ID
	Pos int // End of the word
}

// Match returns matchings obtained from searching array of rune.
func (pma *PMA) Match(input []rune) []Match {
	match := make([]Match, 0)

	now := pma.Init
	for pos, c := range input {
		next, ok := now.Success[c]
		for !ok && now != pma.Init {
			now = now.Failure
			next, ok = now.Success[c]
		}

		if ok {
			now = next
		} else {
			now = pma.Init
		}

		if now.Accept {
			match = append(match, Match{now.ID, pos})
		}

		for _, id := range now.Matching {
			m := Match{id, pos}
			match = append(match, m)
		}
	}

	return match
}

// MatchWord returns matchings obtained from searching string.
func (pma *PMA) MatchWord(input string) []Match {
	return pma.Match([]rune(input))
}

// ShowMatchingWords shows matching words.
func (pma *PMA) ShowMatchingWords(input string, words []string, corr map[int]int) {
	result := pma.MatchWord(input)

	for _, m := range result {
		w := []rune(words[corr[m.ID]])
		fmt.Printf("%v,%v: %v\n", m.Pos-len(w)+1, m.Pos, string(w))
	}
}

// Queue

type queue struct {
	container []*State
}

func newQueue() *queue {
	return &queue{make([]*State, 0)}
}

func (q *queue) enqueue(s *State) *State {
	q.container = append(q.container, s)
	return s
}

func (q *queue) dequeue() *State {
	s := q.container[0]
	q.container = q.container[1:]
	return s
}

func (q *queue) isEmpty() bool {
	return len(q.container) == 0
}
