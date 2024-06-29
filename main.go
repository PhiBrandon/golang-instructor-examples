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
	"github.com/phibrandon/instructor_demo/pkg/input"
	"github.com/phibrandon/instructor_demo/pkg/output"
)

func ExtractUserDetail(client *instructor.InstructorAnthropic, ctx context.Context, content string) (anthropic.MessagesResponse, *output.UserDetail, error) {
	var user output.UserDetail
	prompt := input.Prompt{Instruction: "Given a job description, list all of the perceived and real problems and obstacles that the job poster could or is currently facing.", Input: content}
	resp, err := client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:     anthropic.ModelClaude3Haiku20240307,
		Messages:  []anthropic.Message{anthropic.NewUserTextMessage(prompt.Createprompt())},
		MaxTokens: 4000,
	},
		&user,
	)
	return resp, &user, err
}

func ExtractProblems(client *instructor.InstructorAnthropic, ctx context.Context, content string) (anthropic.MessagesResponse, *output.Problems, error) {
	var problems output.Problems
	prompt := input.Prompt{Instruction: "Given a job description, list all of the perceived and real problems and obstacles that the job poster could or is currently facing.", Input: content}
	resp, err := client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:     anthropic.ModelClaude3Haiku20240307,
		Messages:  []anthropic.Message{anthropic.NewUserTextMessage(prompt.Createprompt())},
		MaxTokens: 4000,
	},
		&problems,
	)
	return resp, &problems, err
}

func RunUserDetail(client *instructor.InstructorAnthropic, ctx context.Context) {
	userFile, err := os.ReadFile("user_detail.txt")
	if err != nil {
		log.Fatal()
	}
	lines := strings.Split(string(userFile), "\n")
	for _, user := range lines {
		_, u, e := ExtractUserDetail(client, ctx, user)
		if e != nil {
			log.Fatal(e)
		}
		fmt.Println(u.Age)
		fmt.Println(u.Name)
	}
}

func RunProblems(client *instructor.InstructorAnthropic, ctx context.Context) {
	jobDescription, err := os.ReadFile("job_description.txt")
	if err != nil {
		log.Fatal(err)
	}
	_, u, e := ExtractProblems(client, ctx, string(jobDescription))
	if e != nil {
		log.Fatal(e)
	}
	for _, problem := range u.Problems {
		fmt.Println(problem.Id)
		fmt.Println(problem.Problem)
	}
}

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load environment, check your files.")
	}
	client := instructor.FromAnthropic(anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY")), instructor.WithMode(instructor.ModeToolCall), instructor.WithMaxRetries(3))
	//RunUserDetail(client, ctx)
	RunProblems(client, ctx)

}
