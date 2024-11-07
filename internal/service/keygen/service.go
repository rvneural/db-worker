package keygen

import (
	"math/rand"
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
	parts := make([]string, 4)

	line := ""

	for i := range parts {
		parts[i] = ""
	}

	for i := range parts {
		for j := 0; j < partLen; j++ {
			parts[i] += string(g.letters[rand.Intn(len(g.letters))])
		}
	}

	for i := range parts {
		line += parts[i]
		if i < len(parts)-1 {
			line += "-"
		}
	}
	return line
}
