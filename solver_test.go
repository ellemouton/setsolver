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
				card("One", "Red", "Diamond", "Hollow"),
				card("One", "Red", "Squiggle", "Hollow"),
				card("One", "Red", "Oval", "Hollow"),
			},
			expectedSets: [][]*Card{
				{
					card("One", "Red", "Diamond", "Hollow"),
					card("One", "Red", "Squiggle", "Hollow"),
					card("One", "Red", "Oval", "Hollow"),
				},
			},
		},
		{
			name: "no sets",
			cards: []*Card{
				card("One", "Red", "Diamond", "Hollow"),
				card("One", "Red", "Squiggle", "Hollow"),
				card("One", "Purple", "Oval", "Hollow"),
			},
			expectedSets: nil,
		},
		{
			name: "one set with different count",
			cards: []*Card{
				card("One", "Red", "Oval", "Hollow"),
				card("Two", "Red", "Oval", "Hollow"),
				card("Two", "Red", "Oval", "Solid"),
				card("Three", "Red", "Oval", "Hollow"),
			},
			expectedSets: [][]*Card{
				{
					card("One", "Red", "Oval", "Hollow"),
					card("Two", "Red", "Oval", "Hollow"),
					card("Three", "Red", "Oval", "Hollow"),
				},
			},
		},
		{
			name: "one set with all different attributes",
			cards: []*Card{
				card("One", "Red", "Squiggle", "Hollow"),
				card("Two", "Purple", "Oval", "Solid"),
				card("Three", "Green", "Diamond", "Shaded"),
			},
			expectedSets: [][]*Card{
				{
					card("One", "Red", "Squiggle", "Hollow"),
					card("Two", "Purple", "Oval", "Solid"),
					card("Three", "Green", "Diamond", "Shaded"),
				},
			},
		},
		{
			name: "two possible sets",
			cards: []*Card{
				card("One", "Red", "Squiggle", "Hollow"),
				card("Two", "Purple", "Oval", "Solid"),
				card("Two", "Green", "Oval", "Shaded"),
				card("Three", "Green", "Diamond", "Shaded"),
				card("Three", "Purple", "Diamond", "Solid"),
			},
			expectedSets: [][]*Card{
				{
					card("One", "Red", "Squiggle", "Hollow"),
					card("Two", "Purple", "Oval", "Solid"),
					card("Three", "Green", "Diamond", "Shaded"),
				},
				{
					card("One", "Red", "Squiggle", "Hollow"),
					card("Two", "Green", "Oval", "Shaded"),
					card("Three", "Purple", "Diamond", "Solid"),
				},
			},
		},
		{
			name: "today's online puzzle",
			cards: []*Card{
				card("Two", "Purple", "Squiggle", "Shaded"),
				card("Three", "Red", "Diamond", "Solid"),
				card("One", "Red", "Squiggle", "Hollow"),

				card("Two", "Red", "Diamond", "Hollow"),
				card("One", "Green", "Diamond", "Solid"),
				card("Two", "Green", "Oval", "Solid"),

				card("One", "Red", "Oval", "Shaded"),
				card("Three", "Green", "Oval", "Solid"),
				card("One", "Red", "Diamond", "Solid"),

				card("Three", "Red", "Diamond", "Shaded"),
				card("Three", "Green", "Squiggle", "Solid"),
				card("Three", "Purple", "Squiggle", "Hollow"),
			},
			expectedSets: [][]*Card{
				{
					card("Two", "Purple", "Squiggle", "Shaded"),
					card("One", "Red", "Squiggle", "Hollow"),
					card("Three", "Green", "Squiggle", "Solid"),
				},
				{
					card("Two", "Purple", "Squiggle", "Shaded"),
					card("Two", "Red", "Diamond", "Hollow"),
					card("Two", "Green", "Oval", "Solid"),
				},
				{
					card("One", "Red", "Squiggle", "Hollow"),
					card("One", "Red", "Oval", "Shaded"),
					card("One", "Red", "Diamond", "Solid"),
				},
				{
					card("Two", "Red", "Diamond", "Hollow"),
					card("One", "Red", "Diamond", "Solid"),
					card("Three", "Red", "Diamond", "Shaded"),
				},
				{
					card("One", "Green", "Diamond", "Solid"),
					card("Two", "Green", "Oval", "Solid"),
					card("Three", "Green", "Squiggle", "Solid"),
				},
				{
					card("Three", "Green", "Oval", "Solid"),
					card("Three", "Red", "Diamond", "Shaded"),
					card("Three", "Purple", "Squiggle", "Hollow"),
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
