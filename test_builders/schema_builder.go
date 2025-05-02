package test_builder

import oas_struct "github.com/alexplayer15/parmesan/data"

type SchemaBuilder struct {
	schema oas_struct.Schema
}

func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		schema: oas_struct.Schema{
			Type:       "object",
			Properties: make(map[string]oas_struct.Property),
		},
	}
}

func (b *SchemaBuilder) WithType(schemaType string) *SchemaBuilder {
	b.schema.Type = schemaType
	return b
}

func (b *SchemaBuilder) AddProperty(name string, property oas_struct.Property) *SchemaBuilder {
	b.schema.Properties[name] = property
	return b
}

func (b *SchemaBuilder) Build() *oas_struct.Schema {
	return &b.schema
}
