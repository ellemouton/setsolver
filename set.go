package setsolver

import "fmt"

// Set is a set of Cards that make up a valid set. A set is valid if for all the
// available attributes, either all cards in the set share the variant of that
// attribute OR none of them do. In other words it should never the case that
// two cards share something that is not shared by other cards in the set.
type Set struct {
	Cards []*Card

	// properties defines if a given attribute is shared across all cards
	// in the set (ie all cards have the same variant for that attribute)
	// or if the attribute variant differs for all cards. This will be nil
	// until the Cards list has at least 2 cards in it.
	properties map[AttributeName]propertyState
}

// NewSet uses the given cards to construct a Set. An error will be returned if
// the given set is not valid.
func NewSet(cards []*Card) (*Set, error) {
	set := &Set{
		properties: make(map[AttributeName]propertyState),
	}

	for _, card := range cards {
		if !set.Add(card) {
			return nil, fmt.Errorf("invalid set")
		}
	}

	return set, nil
}

func (s *Set) Add(card *Card) bool {
	switch len(s.Cards) {
	// If the Set doesn't have any card in it yet, then adding this card is
	// a valid set.
	case 0:
		s.Cards = []*Card{card}

		return true

	// If the Set currently has only 1 card in it then this is also a
	// special case since adding any new card will still be a valid set
	// _and_ this is where we can now define the properties of the set.
	case 1:
		for attribute, variant := range s.Cards[0].Attributes {
			if variant == card.Attributes[attribute] {
				s.properties[attribute] = &propertyStateSame{
					variant: variant,
				}

				continue
			}

			s.properties[attribute] = &propertyStateDifferent{
				usedVariants: map[VariantName]bool{
					variant:                    true,
					card.Attributes[attribute]: true,
				},
			}
		}

		s.Cards = append(s.Cards, card)

		return true

	default:
	}

	// Otherwise, we have already established the properties of the set and
	// so now need to decide if the new card belongs or not.
	for attribute, propertyState := range s.properties {
		if !propertyState.canAdd(card.Attributes[attribute]) {
			return false
		}
	}

	s.Cards = append(s.Cards, card)

	return true
}

type propertyState interface {
	canAdd(variant VariantName) bool
}

type propertyStateSame struct {
	// variant is the variant of an attribute that all the cards in the set
	// should have.
	variant VariantName
}

func (ps *propertyStateSame) canAdd(variant VariantName) bool {
	return variant == ps.variant
}

type propertyStateDifferent struct {
	// usedVariants is a map of all the variants of an attribute that are
	// already taken. Any new card to the set should have a variant that is
	// not already in this list for the attribute in question.
	usedVariants map[VariantName]bool
}

func (ps *propertyStateDifferent) canAdd(variant VariantName) bool {
	return !ps.usedVariants[variant]
}
