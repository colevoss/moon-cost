package moon

import "github.com/gofrs/uuid/v5"

var DefaultUUIDGenerator = NewUUIDGenerator()
var NilID = uuid.Nil

type UUID = uuid.UUID

type IDGenerator interface {
	ID() (uuid.UUID, error)
}

type UUIDGenerator struct {
	gen *uuid.Gen
}

func (u *UUIDGenerator) ID() (UUID, error) {
	return u.gen.NewV7()
}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{
		gen: uuid.NewGen(),
	}
}
