package card

import "math/rand"

type Suite string

const (
	Spades   Suite = ":spades:"
	Hearts   Suite = ":heart:"
	Diamonds Suite = ":diamonds:"
	Clubs    Suite = ":clubs:"
)

const (
	Joker = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

var rank = map[int]string{
	Joker: "Joker",
	Ace:   "Ace",
	Two:   "Two",
	Three: "Three",
	Four:  "Four",
	Five:  "Five",
	Six:   "Six",
	Seven: "Seven",
	Eight: "Eight",
	Nine:  "Nine",
	Ten:   "Ten",
	Jack:  "Jack",
	Queen: "Queen",
	King:  "King",
}

type card struct {
	suit Suite
	rank int
}

type Card interface {
	Suit() Suite
	RankString() string
	Rank() int
	IsJoker() bool
	BJscore() int
}

func (c card) Suit() Suite {
	return c.suit
}

func (c card) RankString() string {
	return rank[c.rank]
}

func (c card) Rank() int {
	return c.rank
}

func (c card) IsJoker() bool {
	return c.rank == Joker
}

func (c card) BJscore() int {
	if c.rank == Ace {
		return 11
	}

	if c.rank > 10 {
		return 10
	}
	return c.rank
}

func New(suit Suite, rank int) Card {
	return card{suit: suit, rank: rank}
}

func NewDeck() []Card {
	cards := []Card{}

	// Add 4 suite of cards
	for _, suit := range []Suite{Spades, Hearts, Diamonds, Clubs} {
		for _, rank := range []int{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King} {
			cards = append(cards, New(suit, rank))
		}
	}
	// Add 2 jokers
	cards = append(cards, New("", Joker))
	cards = append(cards, New("", Joker))

	return cards
}

func Shuffle(cards []Card) []Card {
	shuffled := make([]Card, len(cards))
	perm := rand.Perm(len(cards))

	for i, v := range perm {
		shuffled[v] = cards[i]
	}

	return shuffled
}
