# Translation

This is demo of translation of Rails internationalization file to another
locale.

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
