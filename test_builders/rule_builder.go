package test_builder

import "github.com/alexplayer15/parmesan/data"

type RuleBuilder struct {
	rule   data.Rule
	parent *RuleSetBuilder
}

func (b *RuleBuilder) AddHeaderExtract(name, as string) *RuleBuilder {
	if b.rule.Extract == nil {
		b.rule.Extract = &data.ExtractBlock{}
	}
	b.rule.Extract.Headers = append(b.rule.Extract.Headers, data.HeaderExtract{Name: name, As: as})
	return b
}

func (b *RuleBuilder) AddBodyExtract(path, as, typ string) *RuleBuilder {
	if b.rule.Extract == nil {
		b.rule.Extract = &data.ExtractBlock{}
	}
	b.rule.Extract.Body = append(b.rule.Extract.Body, data.BodyExtract{Path: path, As: as, Type: typ})
	return b
}

func (b *RuleBuilder) AddHeaderInject(name, from string) *RuleBuilder {
	if b.rule.Inject == nil {
		b.rule.Inject = &data.InjectBlock{}
	}
	b.rule.Inject.Headers = append(b.rule.Inject.Headers, data.HeaderInject{Name: name, From: from})
	return b
}

func (b *RuleBuilder) AddBodyInject(path, from, typ string) *RuleBuilder {
	if b.rule.Inject == nil {
		b.rule.Inject = &data.InjectBlock{}
	}
	b.rule.Inject.Body = append(b.rule.Inject.Body, data.BodyInject{Path: path, From: from, Type: typ})
	return b
}

func (b *RuleBuilder) EndRule() *RuleSetBuilder {
	b.parent.rules = append(b.parent.rules, b.rule)
	return b.parent
}
