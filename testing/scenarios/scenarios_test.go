package scenarios

import "fmt"

type lang string

const (
	en lang = "EN"
	es lang = "ES"
)

type greeter struct {
	lang lang
}

func NewGreeter(l string) greeter {
	return greeter{
		lang: lang(l),
	}
}

func (g greeter) Hello(name string) (string, error) {
	switch g.lang {
	case es:
		return fmt.Sprintf("hola, %s", name), nil
	case en:
		return fmt.Sprintf("hi, %s", name), nil
	default:
		return "", fmt.Errorf("I don't know how to say hello in %s!", g.lang)
	}
}

func (g greeter) Goodbye() (string, error) {
	switch g.lang {
	case es:
		return "adios!", nil
	case en:
		return "goodbye!", nil
	default:
		return "", fmt.Errorf("I don't know how to say goodbye in %s!", g.lang)
	}
}
