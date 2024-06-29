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

const (
	maxTokens = 4000
	model     = anthropic.ModelClaude3Haiku20240307
)

type Client struct {
	*instructor.InstructorAnthropic
}

func NewClient() (*Client, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading environment: %w", err)
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}

	client := instructor.FromAnthropic(
		anthropic.NewClient(apiKey),
		instructor.WithMode(instructor.ModeToolCall),
		instructor.WithMaxRetries(3),
	)

	return &Client{client}, nil
}

func (c *Client) ExtractUserDetail(ctx context.Context, content string) (*output.UserDetail, error) {
	var user output.UserDetail
	prompt := input.Prompt{
		Instruction: "Given a job description, list all of the perceived and real problems and obstacles that the job poster could or is currently facing.",
		Input:       content,
	}

	_, err := c.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:     model,
		Messages:  []anthropic.Message{anthropic.NewUserTextMessage(prompt.Createprompt())},
		MaxTokens: maxTokens,
	}, &user)

	return &user, err
}

func (c *Client) ExtractProblems(ctx context.Context, content string) (*output.Problems, error) {
	var problems output.Problems
	prompt := input.Prompt{
		Instruction: "Given a job description, list all of the perceived and real problems and obstacles that the job poster could or is currently facing.",
		Input:       content,
	}

	_, err := c.CreateMessages(ctx, anthropic.MessagesRequest{
		Model:     model,
		Messages:  []anthropic.Message{anthropic.NewUserTextMessage(prompt.Createprompt())},
		MaxTokens: maxTokens,
	}, &problems)

	return &problems, err
}

func (c *Client) RunUserDetail(ctx context.Context) error {
	userFile, err := os.ReadFile("user_detail.txt")
	if err != nil {
		return fmt.Errorf("reading user_detail.txt: %w", err)
	}

	lines := strings.Split(string(userFile), "\n")
	for _, user := range lines {
		u, err := c.ExtractUserDetail(ctx, user)
		if err != nil {
			return fmt.Errorf("extracting user detail: %w", err)
		}
		fmt.Println("Age:", u.Age)
		fmt.Println("Name:", u.Name)
	}

	return nil
}

func (c *Client) RunProblems(ctx context.Context) error {
	jobDescription, err := os.ReadFile("job_description.txt")
	if err != nil {
		return fmt.Errorf("reading job_description.txt: %w", err)
	}

	problems, err := c.ExtractProblems(ctx, string(jobDescription))
	if err != nil {
		return fmt.Errorf("extracting problems: %w", err)
	}

	for _, problem := range problems.Problems {
		fmt.Println("ID:", problem.Id)
		fmt.Println("Problem:", problem.Problem)
	}

	return nil
}

func main() {
	client, err := NewClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	// Uncomment the function you want to run
	// if err := client.RunUserDetail(ctx); err != nil {
	// 	log.Fatalf("Error running user detail: %v", err)
	// }

	if err := client.RunProblems(ctx); err != nil {
		log.Fatalf("Error running problems: %v", err)
	}
}
