package test_builder

import oas_struct "github.com/alexplayer15/parmesan/data"

type ParameterBuilder struct {
	parameter oas_struct.Parameter
}

func NewParameterBuilder() *ParameterBuilder {
	return &ParameterBuilder{
		parameter: oas_struct.Parameter{},
	}
}

func (b *ParameterBuilder) WithName(name string) *ParameterBuilder {
	b.parameter.Name = name
	return b
}

func (b *ParameterBuilder) WithIn(in string) *ParameterBuilder {
	b.parameter.In = in
	return b
}

func (b *ParameterBuilder) WithExample(example string) *ParameterBuilder {
	b.parameter.Example = example
	return b
}

func (b *ParameterBuilder) Build() *oas_struct.Parameter {
	return &b.parameter
}
