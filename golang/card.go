package setsolver

import (
	"fmt"
	"strings"
)

type Card struct {
	// ID is an identifier that the caller can pass in order to maintain a
	// mapping to the card.
	ID int

	Attributes map[AttributeName]VariantName
}

func NewDefaultCard(id int, colour, count, shape, fill VariantName) *Card {
	return &Card{
		ID: id,
		Attributes: map[AttributeName]VariantName{
			Colour.Name: colour,
			Count.Name:  count,
			Shape.Name:  shape,
			Fill.Name:   fill,
		},
	}
}

func (c *Card) String() string {
	var desc string
	for attribute, variant := range c.Attributes {
		desc += fmt.Sprintf("%s:%s|", attribute, variant)
	}

	return strings.TrimRight(desc, "|")
}
