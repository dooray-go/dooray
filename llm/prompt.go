package llm

import "strings"

// PromptTemplate provides simple {{key}} variable substitution in prompt strings.
type PromptTemplate struct {
	Template string
}

// NewPromptTemplate creates a new PromptTemplate.
func NewPromptTemplate(tmpl string) *PromptTemplate {
	return &PromptTemplate{Template: tmpl}
}

// Format replaces all {{key}} placeholders with corresponding values from the map.
// Unmatched placeholders are left as-is.
func (p *PromptTemplate) Format(values map[string]string) string {
	result := p.Template
	for k, v := range values {
		result = strings.ReplaceAll(result, "{{"+k+"}}", v)
	}
	return result
}
