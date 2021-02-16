package balanced_string

const (
	openBrace       = "{"
	openBracket     = "["
	openParenthesis = "("

	closeBrace       = "}"
	closeBracket     = "]"
	closeParenthesis = ")"
)

func isBalanced(str string) bool {
	var symbolsStack []string
	for _, symbolByte := range str {
		symbol := string(symbolByte)

		if symbol == openBrace || symbol == openBracket || symbol == openParenthesis {
			symbolsStack = append(symbolsStack, symbol)
		}

		if len(symbolsStack) == 0 {
			return false
		}

		if symbol == closeBrace || symbol == closeParenthesis || symbol == closeBracket {
			prevSymbol := symbolsStack[len(symbolsStack)-1]
			if (prevSymbol == openBrace && symbol == closeBrace) ||
				(prevSymbol == openBracket && symbol == closeBracket) ||
				(prevSymbol == openParenthesis && symbol == closeParenthesis) {
				symbolsStack = symbolsStack[:len(symbolsStack)-1]
			} else {
				symbolsStack = append(symbolsStack, symbol)
			}
		}
	}
	return len(symbolsStack) == 0
}
