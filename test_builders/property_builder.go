package test_builder

import oas_struct "github.com/alexplayer15/parmesan/data"

type PropertyBuilder struct {
	name       string
	property   oas_struct.Property
	properties map[string]oas_struct.Property
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

func (b *PropertyBuilder) WithDefault(defaultValue any) *PropertyBuilder {
	b.property.Default = defaultValue
	return b
}

func (b *PropertyBuilder) WithItems(items *oas_struct.Schema) *PropertyBuilder {
	b.property.Items = items
	return b
}

func (b *PropertyBuilder) WithItemsRef(ref string) *PropertyBuilder {
	b.property.Items = &oas_struct.Schema{
		Ref: ref,
	}
	return b
}

func (b *PropertyBuilder) WithOneOfRefs(refs []string) *PropertyBuilder {
	for _, ref := range refs {
		b.property.OneOf = append(b.property.OneOf, oas_struct.Schema{
			Ref: ref,
		})
	}
	return b
}

func (b *PropertyBuilder) WithAllOfRefs(refs []string) *PropertyBuilder {
	for _, ref := range refs {
		b.property.AllOf = append(b.property.AllOf, oas_struct.Schema{
			Ref: ref,
		})
	}
	return b
}

func (b *PropertyBuilder) WithFormat(format string) *PropertyBuilder {
	b.property.Format = format
	return b
}

func (b *PropertyBuilder) WithProperty(name string, prop oas_struct.Property) *PropertyBuilder {
	if b.properties == nil {
		b.properties = make(map[string]oas_struct.Property)
	}
	b.properties[name] = prop
	return b
}

func (b *PropertyBuilder) WithRef(ref string) *PropertyBuilder {
	b.property.Ref = ref
	return b
}

func (b *PropertyBuilder) Build() (string, oas_struct.Property) {

	if len(b.properties) > 0 {
		b.property.Properties = b.properties
	}
	return b.name, b.property
}
