package commands

type RuleSet []Rule

type Rule struct {
	ID      string        `yaml:"id"`
	Path    string        `yaml:"path"`
	Method  string        `yaml:"method"`
	Extract *ExtractBlock `yaml:"extract,omitempty"`
	Inject  *InjectBlock  `yaml:"inject,omitempty"`
}

type ExtractBlock struct {
	Headers []HeaderExtract `yaml:"headers,omitempty"`
	Body    []BodyExtract   `yaml:"body,omitempty"`
}

type InjectBlock struct {
	Headers []HeaderInject `yaml:"headers,omitempty"`
	Body    []BodyInject   `yaml:"body,omitempty"`
}

type HeaderExtract struct {
	Name string `yaml:"name"`
	As   string `yaml:"as"`
}

type BodyExtract struct {
	Path string `yaml:"path"`
	As   string `yaml:"as"`
	Type string `yaml:"type"`
}

type HeaderInject struct {
	Name string `yaml:"name"`
	From string `yaml:"from"`
}

type BodyInject struct {
	Path string `yaml:"path"`
	From string `yaml:"from"`
	Type string `yaml:"type"`
}
