package handler

import "strconv"
import "strings"

type TokenType uint8

const (
	TokenTimeUnit TokenType = iota
	TokenTimeFormat
	TokenColon
	TokenDash
	TokenEnd
)

type Token struct {
	value string
	tType TokenType
}

type Lexer struct {
	source string
}

func NewLexer(source string) *Lexer {
	return &Lexer{source}
}

func (lexer *Lexer) Scan() []Token {
	var tokens []Token

	for index := 0; index < len(lexer.source); {
		currentChar := string(lexer.source[index])

		if _, err := strconv.Atoi(currentChar); err == nil {
			tokens = addTimeUnitToken(&tokens, lexer.source, &index)
		} else if strings.ToLower(currentChar) == "a" || strings.ToLower(currentChar) == "p" {
			tokens = addTimeFormatToken(&tokens, lexer.source, &index)
		} else if currentChar == ":" {
			tokens = addSingleToken(&tokens, currentChar, TokenColon, &index)
		} else if currentChar == "-" {
			tokens = addSingleToken(&tokens, currentChar, TokenDash, &index)
		}
	}

	tokens = append(tokens, Token{"", TokenEnd})

	return tokens
}

func addTimeUnitToken(tokens *[]Token, timestampStr string, index *int) []Token {
	value := ""

	// Exits loop if:
	// - there are more than two characters per unit.
	// - a non-digit character is encountered
	// - the end of the string is reached
	for count := 0; count < 2 && *index < len(timestampStr); count++ {
		currentChar := string(timestampStr[*index])

		if _, err := strconv.Atoi(currentChar); err != nil {
			break
		}

		value += currentChar
		*index++
	}

	token := Token{value, TokenTimeUnit}

	return append(*tokens, token)
}

func addTimeFormatToken(tokens *[]Token, timestampStr string, index *int) []Token {
	value := ""

	// Breaks out of loop if there are more than two characters per format.
	for count := 0; count < 2; count++ {
		currentChar := string(timestampStr[*index])

		value += currentChar
		*index++
	}

	token := Token{value, TokenTimeFormat}

	return append(*tokens, token)
}

func addSingleToken(tokens *[]Token, char string, tType TokenType, index *int) []Token {
	token := Token{char, tType}

	*index++

	return append(*tokens, token)
}