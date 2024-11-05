package keygen

import (
	"math/rand"
	"strings"
	"sync"
)

type Generator struct {
	maxLen  int64
	letters []rune
}

func New(maxLen int64) *Generator {
	return &Generator{
		maxLen:  maxLen,
		letters: []rune("abcdefghijklmnopqrstuvwxyz0123456789"),
	}
}

func (g *Generator) Generate() string {
	partLen := rand.Intn(6) + 3
	builderParts := make([]strings.Builder, 4)

	line := ""

	wg := sync.WaitGroup{}
	for _, builder := range builderParts {
		wg.Add(1)
		go func(builder *strings.Builder) {
			defer wg.Done()
			for i := 0; i < partLen; i++ {
				builder.WriteRune(g.letters[rand.Intn(len(g.letters))])
			}
		}(&builder)
	}
	wg.Wait()
	line = builderParts[0].String() + "-" + builderParts[1].String() + "-" + builderParts[2].String() + "-" + builderParts[3].String()
	return line
}
