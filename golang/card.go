package setsolver

import (
	"fmt"
	"strings"
)

type Card struct {
	Attributes map[AttributeName]VariantName
}

func (c *Card) String() string {
	var desc string
	for attribute, variant := range c.Attributes {
		desc += fmt.Sprintf("%s:%s|", attribute, variant)
	}

	return strings.TrimRight(desc, "|")
}
