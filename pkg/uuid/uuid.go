package uuid

import "github.com/google/uuid"

type UUID string

func New(typ rune) UUID {
	return UUID(string(typ) + uuid.New().String()[:8])
}

func Type(uuid UUID) rune {
	return rune(uuid[0])
}
