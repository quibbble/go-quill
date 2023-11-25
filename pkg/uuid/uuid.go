package uuid

import "github.com/google/uuid"

type UUID string

const (
	Nil UUID = "|"
)

func New(typ rune) UUID {
	return UUID(string(typ) + uuid.New().String()[:8])
}

func (u UUID) Type() rune {
	return rune(u[0])
}
