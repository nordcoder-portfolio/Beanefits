package publiccode

import (
	"context"

	"github.com/google/uuid"
)

type Generator struct{}

func NewGenerator() *Generator { return &Generator{} }

func (g *Generator) New(ctx context.Context) (string, error) {
	return uuid.NewString(), nil
}
