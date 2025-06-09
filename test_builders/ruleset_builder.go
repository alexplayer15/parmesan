package test_builder

import "github.com/alexplayer15/parmesan/data"

type RuleSetBuilder struct {
	rules []data.Rule
}

func NewRuleSetBuilder() *RuleSetBuilder {
	return &RuleSetBuilder{}
}

func (b *RuleSetBuilder) AddRule(rule data.Rule) *RuleSetBuilder {
	b.rules = append(b.rules, rule)
	return b
}

func (b *RuleSetBuilder) WithRule(id, path, method string) *RuleBuilder {
	rule := data.Rule{
		ID:     id,
		Path:   path,
		Method: method,
	}
	return &RuleBuilder{rule: rule, parent: b}
}

func (b *RuleSetBuilder) Build() data.RuleSet {
	return b.rules
}
