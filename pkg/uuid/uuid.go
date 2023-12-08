package uuid

import (
	"math/rand"
	"slices"
)

type UUID string

const (
	Nil UUID = "|"
)

const (
	chars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	length = 7
)

type Gen struct {
	r *rand.Rand
}

func NewGen(r *rand.Rand) *Gen {
	return &Gen{
		r: r,
	}
}

func (g *Gen) New(typ rune) UUID {
	c := ""
	for i := 0; i < length; i++ {
		c += string(chars[g.r.Intn(len(chars))])
	}
	return UUID(string(typ) + c)
}

func (u UUID) Type() rune {
	return rune(u[0])
}

// Intersect performs set intersection on a ∩ ...b
func Intersect(a []UUID, b ...[]UUID) []UUID {
	intersect := make([]UUID, 0)
	for _, it := range a {
		found := true
		for _, l := range b {
			if !slices.Contains(l, it) {
				found = false
				break
			}
		}
		if found {
			intersect = append(intersect, it)
		}
	}
	return intersect
}

func Union(a []UUID, b ...[]UUID) []UUID {
	lists := [][]UUID{a}
	items := make([]UUID, 0)
	union := make([]UUID, 0)
	for _, l := range b {
		lists = append(lists, l)
	}
	for _, l := range lists {
		items = append(items, l...)
	}
	for _, it := range items {
		if !slices.Contains(union, it) {
			union = append(union, it)
		}
	}
	return union
}

// Diff performs set difference on a \ b
func Diff(a, b []UUID) []UUID {
	diff := make([]UUID, 0)
	for _, it := range a {
		if !slices.Contains(b, it) {
			diff = append(diff, it)
		}
	}
	return diff
}
