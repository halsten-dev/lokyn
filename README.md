# lokyn
Lokyn is a simple lib to allow for easy localized applications developpement.
## main.go
```go
package main

import "embed"
import "github.com/halsten-dev/lokyn"

//go:embed translations
var translations embed.FS

func main() {
	lokyn.Init()
	err = lokyn.AddTranslationFS(translations, "translations")

	if err != nil {
		log.Fatal(err)
	}

	lokyn.SetLanguage("en")


  fmt.Println(lokyn.L("translation"))
  fmt.Println(lokyn.P("apple", 2))
}
```

## projects structure
```
project root /
  translations /
    en.json
    fr.json
```

## en.json
```json
{
  "translation": "translation", 
  "apple": {
        "one": "apple",
        "other": "apples"
  }
}
```

## fr.json
```json
{
  "translation": "traduction",
  "apple": {
        "one": "pomme",
        "other": "pommes"
  }
}
```
