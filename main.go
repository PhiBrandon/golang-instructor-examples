package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/instructor-ai/instructor-go/pkg/instructor"
	"github.com/joho/godotenv"
	"github.com/liushuangls/go-anthropic/v2"
)

type UserDetail struct {
	Name string `json:"name" jsonschema:"title=the name,description="he name of the user"`
	Age  int    `json:"age" jsonschema:"title=the age, description="the age of the user"`
}

func ExtractUserDetail(client *instructor.InstructorAnthropic, ctx context.Context, content string) (anthropic.MessagesResponse, *UserDetail, error) {
	var user UserDetail
	resp, err := client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:     anthropic.ModelClaude3Haiku20240307,
		Messages:  []anthropic.Message{anthropic.NewUserTextMessage(fmt.Sprintf("Extract the user information: %s", content))},
		MaxTokens: 4000,
	},
		&user,
	)
	return resp, &user, err
}

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load environment, check your files.")
	}
	userFile, err := os.ReadFile("user_detail.txt")
	if err != nil {
		log.Fatal()
	}
	lines := strings.Split(string(userFile), "\n")

	client := instructor.FromAnthropic(anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY")), instructor.WithMode(instructor.ModeToolCall), instructor.WithMaxRetries(3))
	for _, user := range lines {
		r, u, e := ExtractUserDetail(client, ctx, user)
		if e != nil {
			log.Fatal(e)
		}
		fmt.Println(*r.Content[0].Text)
		fmt.Println(u.Age)
		fmt.Println(u.Name)
	}

}
