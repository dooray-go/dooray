package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dooray-go/dooray-sdk/llm"
	"github.com/dooray-go/dooray-sdk/llm/anthropic"
	"github.com/dooray-go/dooray-sdk/openapi/messenger"
)

func main() {
	// 1. LLM Provider 생성
	provider, err := anthropic.New(
		llm.WithModel("claude-sonnet-4-20250514"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. 프롬프트 템플릿 사용
	tmpl := llm.NewPromptTemplate("{{name}}의 오늘 일정을 요약해줘")
	prompt := tmpl.Format(map[string]string{"name": "정지범"})

	// 3. LLM 질의
	answer, err := provider.Query(context.Background(), prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("LLM 응답:", answer)

	// 4. Dooray 메신저로 결과 전송
	m := messenger.NewDefaultMessenger()
	_, err = m.DirectSend(os.Getenv("DOORAY_API_KEY"), &messenger.DirectSendRequest{
		Text:                 answer,
		OrganizationMemberId: "member-id",
	})
	if err != nil {
		log.Fatal(err)
	}
}
