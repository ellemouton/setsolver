package setsolver

import "fmt"

var (
	Colour = Attribute{
		Name: "Colour",
		Variants: map[VariantName]struct{}{
			"Red":    {},
			"Green":  {},
			"Purple": {},
		},
	}
)

type Variant struct {
	Name      VariantName
	Attribute Attribute
}

func NewVariant(name VariantName, attribute Attribute) (*Variant, error) {
	if _, ok := attribute.Variants[name]; !ok {
		return nil, fmt.Errorf("%s is not a valid variant of "+
			"attribute %s", name, attribute.Name)
	}

	return &Variant{
		Name:      name,
		Attribute: attribute,
	}, nil
}

type VariantName string

type AttributeName string

// Attribute describes one attribute of a card. If the game of SET uses a
// variant count of 3 then the attribute must have 3 different variants.
type Attribute struct {
	// Name uniquely describes this attribute. It must not clash with the
	// name of another attribute.
	Name AttributeName

	// Variants returns the variants of this attribute. This must match
	// the variant count of the game being played and will be checked on
	// attribution registration. It's a map so that uniqueness of variants
	// is enforced.
	Variants map[VariantName]struct{}
}
