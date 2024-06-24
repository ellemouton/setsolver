package setsolver

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSolver(t *testing.T) {
	solver, err := NewSetSolver(&Config{
		VariantCount: 3,
		SetSize:      3,
	}, Colour, Shape, Fill, Count)
	require.NoError(t, err)

	tests := []struct {
		name         string
		cards        []*Card
		expectedSets [][]*Card
	}{
		{
			name: "basic one set with diff shape",
			cards: []*Card{
				card("one", "red", "diamond", "hollow"),
				card("one", "red", "squiggle", "hollow"),
				card("one", "red", "oval", "hollow"),
			},
			expectedSets: [][]*Card{
				{
					card("one", "red", "diamond", "hollow"),
					card("one", "red", "squiggle", "hollow"),
					card("one", "red", "oval", "hollow"),
				},
			},
		},
		{
			name: "no sets",
			cards: []*Card{
				card("one", "red", "diamond", "hollow"),
				card("one", "red", "squiggle", "hollow"),
				card("one", "purple", "oval", "hollow"),
			},
			expectedSets: nil,
		},
		{
			name: "one set with different count",
			cards: []*Card{
				card("one", "red", "oval", "hollow"),
				card("two", "red", "oval", "hollow"),
				card("two", "red", "oval", "solid"),
				card("three", "red", "oval", "hollow"),
			},
			expectedSets: [][]*Card{
				{
					card("one", "red", "oval", "hollow"),
					card("two", "red", "oval", "hollow"),
					card("three", "red", "oval", "hollow"),
				},
			},
		},
		{
			name: "one set with all different attributes",
			cards: []*Card{
				card("one", "red", "squiggle", "hollow"),
				card("two", "purple", "oval", "solid"),
				card("three", "green", "diamond", "shaded"),
			},
			expectedSets: [][]*Card{
				{
					card("one", "red", "squiggle", "hollow"),
					card("two", "purple", "oval", "solid"),
					card("three", "green", "diamond", "shaded"),
				},
			},
		},
		{
			name: "two possible sets",
			cards: []*Card{
				card("one", "red", "squiggle", "hollow"),
				card("two", "purple", "oval", "solid"),
				card("two", "green", "oval", "shaded"),
				card("three", "green", "diamond", "shaded"),
				card("three", "purple", "diamond", "solid"),
			},
			expectedSets: [][]*Card{
				{
					card("one", "red", "squiggle", "hollow"),
					card("two", "purple", "oval", "solid"),
					card("three", "green", "diamond", "shaded"),
				},
				{
					card("one", "red", "squiggle", "hollow"),
					card("two", "green", "oval", "shaded"),
					card("three", "purple", "diamond", "solid"),
				},
			},
		},
		{
			name: "today's online puzzle",
			cards: []*Card{
				card("two", "purple", "squiggle", "shaded"),
				card("three", "red", "diamond", "solid"),
				card("one", "red", "squiggle", "hollow"),

				card("two", "red", "diamond", "hollow"),
				card("one", "green", "diamond", "solid"),
				card("two", "green", "oval", "solid"),

				card("one", "red", "oval", "shaded"),
				card("three", "green", "oval", "solid"),
				card("one", "red", "diamond", "solid"),

				card("three", "red", "diamond", "shaded"),
				card("three", "green", "squiggle", "solid"),
				card("three", "purple", "squiggle", "hollow"),
			},
			expectedSets: [][]*Card{
				{
					card("two", "purple", "squiggle", "shaded"),
					card("one", "red", "squiggle", "hollow"),
					card("three", "green", "squiggle", "solid"),
				},
				{
					card("two", "purple", "squiggle", "shaded"),
					card("two", "red", "diamond", "hollow"),
					card("two", "green", "oval", "solid"),
				},
				{
					card("one", "red", "squiggle", "hollow"),
					card("one", "red", "oval", "shaded"),
					card("one", "red", "diamond", "solid"),
				},
				{
					card("two", "red", "diamond", "hollow"),
					card("one", "red", "diamond", "solid"),
					card("three", "red", "diamond", "shaded"),
				},
				{
					card("one", "green", "diamond", "solid"),
					card("two", "green", "oval", "solid"),
					card("three", "green", "squiggle", "solid"),
				},
				{
					card("three", "green", "oval", "solid"),
					card("three", "red", "diamond", "shaded"),
					card("three", "purple", "squiggle", "hollow"),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sets, err := solver.Solve(test.cards)
			require.NoError(t, err)

			var resultSets [][]*Card
			for _, s := range sets {
				resultSets = append(resultSets, s.Cards)
			}

			require.EqualValues(t, test.expectedSets, resultSets)
		})
	}
}

func card(count, colour, shape, fill VariantName) *Card {
	return &Card{Attributes: map[AttributeName]VariantName{
		Colour.Name: colour,
		Count.Name:  count,
		Shape.Name:  shape,
		Fill.Name:   fill,
	}}
}
