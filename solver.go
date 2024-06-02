package setsolver

import "fmt"

type SetSolver interface {
	Solve(cards []Card) (*Set, error)
}

type Config struct {
	// VariantCount is the number of different variants that each attribute
	// in the game should have.
	VariantCount uint

	// SetSize is the number of cards in a set. SetSize can never be greater
	// than VariantCount. It can only be less than or equal to the
	// VariantCount.
	SetSize uint
}

type setSolver struct {
	cfg *Config

	attributes map[AttributeName]Attribute
}

func NewSetSolver(cfg *Config, attributes ...Attribute) (SetSolver, error) {
	if cfg.SetSize > cfg.VariantCount {
		return nil, fmt.Errorf("the set size of %d is greater than "+
			"the variant count %d which is not a valid game",
			cfg.SetSize, cfg.VariantCount)
	}

	attributeList := make(map[AttributeName]Attribute)
	for _, attribute := range attributes {
		name := attribute.Name()

		if _, ok := attributeList[name]; ok {
			return nil, fmt.Errorf("an attribute named %s has "+
				"already been registered", name)
		}

		count := len(attribute.Variants())
		if count != int(cfg.VariantCount) {
			return nil, fmt.Errorf("expected a variant count of "+
				"%d but attribute %s has a variant count of "+
				"%d", cfg.VariantCount, name, count)
		}
	}

	return &setSolver{
		cfg:        cfg,
		attributes: attributeList,
	}, nil
}

func (s *setSolver) Solve(cards []Card) (*Set, error) {
	if err := s.validateCards(cards); err != nil {
		return nil, err
	}

}

func (s *setSolver) validateCards(cards []Card) error {
	// Check that all the cards given have all the attributes that have
	// been registered to this solver and that the variant of that attribute
	// is known.
	for _, card := range cards {
		if len(card.Attributes) != len(s.attributes) {
			return fmt.Errorf("card(%s) does not have the "+
				"correct number of attributes of %d", card,
				len(s.attributes))
		}

		for _, attribute := range s.attributes {
			variant, ok := card.Attributes[attribute.Name()]
			if !ok {
				return fmt.Errorf("card(%s) is missing "+
					"attribute %s", cards, attribute.Name())
			}

			_, ok = attribute.Variants()[variant]
			if !ok {
				return fmt.Errorf("card(%s) has an unknown "+
					"varient (%s) for attribute %s. "+
					"Expected one of: %v", cards, variant,
					attribute.Name(), attribute.Variants())
			}
		}
	}

	return nil
}
