package test_builder

import oas_struct "github.com/alexplayer15/parmesan/data"

type PropertyBuilder struct {
	name     string
	property oas_struct.Property
}

func NewPropertyBuilder() *PropertyBuilder {
	return &PropertyBuilder{
		property: oas_struct.Property{},
	}
}

func (b *PropertyBuilder) WithName(name string) *PropertyBuilder {
	b.name = name
	return b
}

func (b *PropertyBuilder) WithType(t string) *PropertyBuilder {
	b.property.Type = t
	return b
}

func (b *PropertyBuilder) WithDescription(desc string) *PropertyBuilder {
	b.property.Description = desc
	return b
}

func (b *PropertyBuilder) WithExample(example any) *PropertyBuilder {
	b.property.Example = example
	return b
}

func (b *PropertyBuilder) Build() (string, oas_struct.Property) {
	return b.name, b.property
}
