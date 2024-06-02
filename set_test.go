package setsolver

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestSetBuilding tests that a valid set is always built when the NewSet and
// (Set).Add methods are used.
func TestSetBuilding(t *testing.T) {
	t.Run("basic single attribute", func(t *testing.T) {
		set, err := NewSet(nil)
		require.NoError(t, err)

		// Adding the first card to the set should work.
		card1 := &Card{Attributes: map[AttributeName]VariantName{
			"Colour": "Red",
		}}
		require.True(t, set.Add(card1))

		// Adding a second card to the set should also work. This card a
		// different colour to the first card.
		card2 := &Card{Attributes: map[AttributeName]VariantName{
			"Colour": "Blue",
		}}
		require.True(t, set.Add(card2))

		// Adding a third card that has the same colour as one of the
		// existing cards in the set should fail.
		card3 := &Card{Attributes: map[AttributeName]VariantName{
			"Colour": "Blue",
		}}
		require.False(t, set.Add(card3))

		card3 = &Card{Attributes: map[AttributeName]VariantName{
			"Colour": "Red",
		}}
		require.False(t, set.Add(card3))

		// However, adding a third card with a different colour should
		// work.
		card3 = &Card{Attributes: map[AttributeName]VariantName{
			"Colour": "Purple",
		}}
		require.True(t, set.Add(card3))

		// Also check that NewSet does this full check properly if
		// initialised with a non-nil set of cards.
		cards := []*Card{
			{Attributes: map[AttributeName]VariantName{
				"Colour": "Red",
			}},
			{Attributes: map[AttributeName]VariantName{
				"Colour": "Blue",
			}},
			{Attributes: map[AttributeName]VariantName{
				"Colour": "Green",
			}},
		}
		_, err = NewSet(cards)
		require.NoError(t, err)

		cards = []*Card{
			{Attributes: map[AttributeName]VariantName{
				"Colour": "Red",
			}},
			{Attributes: map[AttributeName]VariantName{
				"Colour": "Blue",
			}},
			{Attributes: map[AttributeName]VariantName{
				"Colour": "Blue",
			}},
		}
		_, err = NewSet(cards)
		require.ErrorContains(t, err, "invalid set")
	})

	t.Run("two attributes", func(t *testing.T) {
		set, err := NewSet(nil)
		require.NoError(t, err)

		// The first card always works.
		ok := set.Add(&Card{Attributes: map[AttributeName]VariantName{
			"Colour": "Red",
			"Shape":  "Diamond",
		}})
		require.True(t, ok)

		// For the second card, let one attribute match and one differ.
		ok = set.Add(&Card{Attributes: map[AttributeName]VariantName{
			"Colour": "Red",
			"Shape":  "Squiggle",
		}})
		require.True(t, ok)

		// Try adding a card with a valid shape but not a valid colour.
		ok = set.Add(&Card{Attributes: map[AttributeName]VariantName{
			"Colour": "Green",
			"Shape":  "Oval",
		}})
		require.False(t, ok)

		// Try adding a card with a valid colour but not a valid shape.
		ok = set.Add(&Card{Attributes: map[AttributeName]VariantName{
			"Colour": "Red",
			"Shape":  "Diamond",
		}})
		require.False(t, ok)

		// Finally, add a card with valid colour and shape.
		ok = set.Add(&Card{Attributes: map[AttributeName]VariantName{
			"Colour": "Red",
			"Shape":  "Oval",
		}})
		require.True(t, ok)
	})
}
