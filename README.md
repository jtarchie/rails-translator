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

## Installation

1. **Download Builder**:

   Using Homebrew:
   
   ```bash
   brew tap jtarchie/translation https://github.com/jtarchie/rails-translator
   brew install rails-translator
   ```

## Usage

1. Translate the example from English to Japanese:

   ```bash
   rails-translator \
    --from-filename examples/en.yml \
    --from-language "British English" \
    --to-filename examples/jp.yml \
    --to-language "Japanese" \
    --openai-access-token $YOUR_OPENAI_API_TOKEN
   ```

1. Check the file `examples/jp.yml`.

   English:

   ```yaml
   en:
     welcome_html: "<b>Welcome %{username}!</b>"
     static_pages:
       index:
         welcome: "Welcome!"
         services_html: "We provide some fancy services to <em>good people</em>."
   ```

   Japanese (translation):

   ```yaml
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
