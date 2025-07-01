package engine

// func Translate(strokes []string, dict dictionary.Dictionary) string {
// 	var output []string
// 	for _, stroke := range strokes {
// 		if word, ok := dict[stroke]; ok {
// 			output = append(output, word)
// 		} else {
// 			output = append(output, "["+stroke+"]") // show unmapped strokes
// 		}
// 	}
// 	return strings.Join(output, " ")
// }

type Engine struct {
	Dict Dictionary
}

func NewEngine(dict Dictionary) *Engine {
	return &Engine{Dict: dict}
}

func (e *Engine) TranslateSteno(strokeText string) string {
	stroke, err := NewStrokeFromSteno(strokeText)
	if err != nil {
		return "[error]"
	}
	word, ok := e.Dict.Lookup(stroke.Steno())
	if !ok {
		return "[" + stroke.Steno() + "]"
	}
	return word
}
