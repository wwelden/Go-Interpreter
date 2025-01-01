package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEN_OBJ  = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

type Null struct{}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}
type Boolean struct {
	Value bool
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEN_OBJ }

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }
