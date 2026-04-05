package llm

import "testing"

func TestPromptTemplate_Format(t *testing.T) {
	tmpl := NewPromptTemplate("{{name}}의 오늘 일정을 요약해줘")
	result := tmpl.Format(map[string]string{"name": "정지범"})
	expected := "정지범의 오늘 일정을 요약해줘"
	if result != expected {
		t.Errorf("want %q, got %q", expected, result)
	}
}

func TestPromptTemplate_MultipleVars(t *testing.T) {
	tmpl := NewPromptTemplate("{{project}} 프로젝트에서 {{member}}에게 할당된 업무를 알려줘")
	result := tmpl.Format(map[string]string{
		"project": "dooray-sdk",
		"member":  "정지범",
	})
	expected := "dooray-sdk 프로젝트에서 정지범에게 할당된 업무를 알려줘"
	if result != expected {
		t.Errorf("want %q, got %q", expected, result)
	}
}

func TestPromptTemplate_UnmatchedPlaceholder(t *testing.T) {
	tmpl := NewPromptTemplate("{{name}}의 {{task}}를 확인해줘")
	result := tmpl.Format(map[string]string{"name": "정지범"})
	expected := "정지범의 {{task}}를 확인해줘"
	if result != expected {
		t.Errorf("want %q, got %q", expected, result)
	}
}
