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
