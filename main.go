package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

type (
	Language map[string]interface{}
	Payload  map[string]Language
)

type CLI struct {
	BaseURL           string `default:"https://api.openai.com/v1"         help:"url of the OpenAI HTTP domain"`
	FromFilename      string `help:"filename of the original language"    required:""`
	FromLanguage      string `default:"American English"                  help:"language to translate from"       required:""`
	OpenAIAccessToken string `env:"OPENAI_ACCESS_TOKEN"                   help:"the API token for the OpenAI API" required:""`
	ToFilename        string `help:"filename to save the translations to" required:""`
	ToLanguage        string `help:"language to translate to"             required:""`
	Model             string `default:"gpt-3.5-turbo"                     enum:"gpt-3.5-turbo,gpt-4"              help:"model to use from OpenAI translation" required:""`
}

const prompt = `
Translate the following message from the locale %q to the locale %q.
Please use the following criteria:
- Ensure HTML tags are maintained.
- Do not translate placeholders that are surrounded by token '%%{' and '}'.
- Please only return the translation, no extraneous explanation.
`

func (c *CLI) translate(value string) (string, error) {
	config := openai.DefaultConfig(c.OpenAIAccessToken)
	config.BaseURL = c.BaseURL
	client := openai.NewClientWithConfig(config)

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: c.Model,
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

func (c *CLI) iterate(node Language) (Language, error) {
	translation := Language{}

	for name, token := range node {
		switch original := token.(type) {
		case string:
			translatedLine, err := c.translate(original)
			if err != nil {
				return nil, fmt.Errorf("could not translate %q to %q: %w", original, c.ToLanguage, err)
			}

			translation[name] = translatedLine
		case Language:
			t, err := c.iterate(original)
			if err != nil {
				return nil, fmt.Errorf("could not translate embedded %q: %w", name, err)
			}

			translation[name] = t
		default:
			return nil, fmt.Errorf("do not understand %#v to be translated", original)
		}
	}

	return translation, nil
}

func (c *CLI) Run() error {
	contents, err := os.ReadFile(c.FromFilename)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	var payload Payload

	err = yaml.Unmarshal(contents, &payload)
	if err != nil {
		return fmt.Errorf("could not unmarshal: %w", err)
	}

	keys := []string{}
	for key := range payload {
		keys = append(keys, key)
	}

	if 1 < len(keys) {
		return fmt.Errorf("input language file (%q) has more than one language", c.FromFilename)
	}

	translation, err := c.iterate(payload[keys[0]])
	if err != nil {
		return fmt.Errorf("could not iterate through language file %q: %w", c.FromFilename, err)
	}

	filename := filepath.Base(c.ToFilename)
	parts := strings.Split(filename, ".")

	newPayload := Payload{}
	newPayload[parts[0]] = translation

	contents, err = yaml.Marshal(&newPayload)
	if err != nil {
		return fmt.Errorf("could not translate to YAML: %w", err)
	}

	err = os.WriteFile(c.ToFilename, contents, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not write file %q: %w", c.ToFilename, err)
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
