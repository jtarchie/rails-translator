package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

type Token interface{}
type Language map[string]Token

type Payload map[string]Language

type CLI struct {
	Filename          string `help:"filename of the original language" required:""`
	FromLanguage      string `help:"language to translate from" required:"" default:"en"`
	OpenAIAccessToken string `help:"the API token for the OpenAI API" required:"" env:"OPENAI_ACCESS_TOKEN"`
	ToLanguage        string `help:"language to translate to" required:""`
}

const prompt = `Translate the following message from the locale %q to the locale %q.`

func (c *CLI) translate(value string) (string, error) {
	client := openai.NewClient(c.OpenAIAccessToken)

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: fmt.Sprintf(prompt, c.FromLanguage, c.ToLanguage),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: value,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("could not translate: %w", err)
	}

	return response.Choices[0].Message.Content, nil
}

func (c *CLI) Run() error {
	contents, err := os.ReadFile(c.Filename)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	var payload Payload

	err = yaml.Unmarshal(contents, &payload)
	if err != nil {
		return fmt.Errorf("could not unmarshal: %w", err)
	}

	if _, ok := payload[c.FromLanguage]; !ok {
		return fmt.Errorf("could not find %q in %q", c.FromLanguage, c.Filename)
	}

	translation := Language{}

	for name, token := range payload[c.FromLanguage] {
		switch v := token.(type) {
		case string:
			value, err := c.translate(v)
			if err != nil {
				return fmt.Errorf("could not translate %q to %q: %w", v, c.ToLanguage, err)
			}

			translation[name] = value
		default:
			return fmt.Errorf("do not understand %v to translated", v)
		}
	}

	newPayload := Payload{}
	newPayload[c.ToLanguage] = translation

	contents, err = yaml.Marshal(&newPayload)
	if err != nil {
		return fmt.Errorf("could not translate to YAML: %w", err)
	}

	newFilename := filepath.Join(filepath.Dir(c.Filename), fmt.Sprintf("%s.yaml", c.ToLanguage))

	err = os.WriteFile(newFilename, contents, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not write file %q: %w", newFilename, err)
	}

	return nil
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
