package setsolver

import "fmt"

type Variant struct {
	Name      VariantName
	Attribute Attribute
}

func NewVariant(name VariantName, attribute Attribute) (*Variant, error) {
	if _, ok := attribute.Variants()[name]; !ok {
		return nil, fmt.Errorf("%s is not a valid variant of "+
			"attribute %s", name, attribute.Name())
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
type Attribute interface {
	// Name uniquely describes this attribute. It must not clash with the
	// name of another attribute.
	Name() AttributeName

	// Variants returns the variants of this attribute. This must match
	// the variant count of the game being played and will be checked on
	// attribution registration. It's a map so that uniqueness of variants
	// is enforced.
	Variants() map[VariantName]struct{}
}

type Colour struct{}

func (c *Colour) Name() AttributeName {
	return "Colour"
}

func (c *Colour) Variants() map[VariantName]struct{} {
	return map[VariantName]struct{}{
		"Red":    {},
		"Green":  {},
		"Purple": {},
	}
}

var _ Attribute = (*Colour)(nil)

type Shape struct{}

func (s *Shape) Name() AttributeName {
	return "Shape"
}

func (s *Shape) Variants() map[VariantName]struct{} {
	return map[VariantName]struct{}{
		"Oval":     {},
		"Diamond":  {},
		"Squiggle": {},
	}
}

var _ Attribute = (*Shape)(nil)

type Count struct{}

func (c *Count) Name() AttributeName {
	return "Count"
}

func (c *Count) Variants() map[VariantName]struct{} {
	return map[VariantName]struct{}{
		"1": {},
		"2": {},
		"3": {},
	}
}

var _ Attribute = (*Count)(nil)

type Shading struct{}

func (s *Shading) Name() AttributeName {
	return "Shading"
}

func (s *Shading) Variants() map[VariantName]struct{} {
	return map[VariantName]struct{}{
		"Solid":  {},
		"Hollow": {},
		"Shaded": {},
	}
}

var _ Attribute = (*Shading)(nil)
