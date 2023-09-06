# Rails Language File Translator

Translate Rails internationalization files from one language to another using
the power of OpenAI's GPT-3.5 Turbo. This tool is designed to handle nested YAML
structures and ensures that HTML tags and placeholders are preserved during the
translation process.

## Features

- Translates Rails language files while preserving HTML tags.
- Does not translate placeholders surrounded by `%%{` and `}`.
- Supports nested YAML structures.
- Automatically saves the translated file with the appropriate language code.

## Usage

1. Set the environment variable `OPENAI_ACCESS_TOKEN` to your OpenAI API token.
1. Translate the example from English to Japanese:

   ```bash
   go run main.go \
     --filename examples/en.yaml \
     --to-language "jp"
   ```

1. Check the file `examples/jp.yml`.

English:

```
en:
  welcome_html: "<b>Welcome %{username}!</b>"
  static_pages:
    index:
      welcome: "Welcome!"
      services_html: "We provide some fancy services to <em>good people</em>."
```

Japanese (translation):

```
jp:
  static_pages:
    index:
      services_html: 私たちは<em>善良な人々</em>にいくつかの華やかなサービスを提供します。
      welcome: ようこそ！
  welcome_html: <b>ようこそ %{username} さん！</b>
```

Keeps the HTML and `%{}` placeholders.

## Limitations

- Ensure that your OpenAI API quota is sufficient for the number of
  translations.
- The tool assumes that the provided YAML file is correctly formatted.
