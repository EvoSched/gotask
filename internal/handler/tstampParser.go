package handler

import "time"
import "errors"
import "strings"
import "strconv"

type Parser struct {
	tokens []Token
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens}
}

func (parser *Parser) Parse() (start *time.Time, end *time.Time, err error) {
	current := 0

	start, err = parseStartTimestamp(parser.tokens, &current)
	if err != nil {
		return start, &time.Time{}, err
	}

	err = parseDash(parser.tokens, &current)
	if err != nil {
		return &time.Time{}, &time.Time{}, err
	}

	end, err = parseEndTimestamp(parser.tokens, &current)
	if err != nil {
		return &time.Time{}, end, err
	}

	if end.Before(*start) {
		return &time.Time{}, &time.Time{}, errors.New("end time cannot occur before start time")
	}

	return start, end, nil
}

func parseStartTimestamp(tokens []Token, index *int) (*time.Time, error) {
	startHour, err := parseTimeUnit(tokens, index)
	if err != nil {
		return &time.Time{}, err
	}

	next, err := parseSymbolAfterHour(tokens, index)
	if err != nil {
		return &time.Time{}, err
	}

	startFormat := ""
	startMinute := 0

	if next == "format" {
		startFormat, err = parseFormat(tokens, index)
		if err != nil {
			return &time.Time{}, err
		}
	} else if next == "minute" {
		if tokens[*index].tType == TokenTimeUnit {
			if len(tokens[*index].value) != 2 {
				return &time.Time{}, errors.New("minutes unit must have two digits")
			}

			startMinute, err = parseTimeUnit(tokens, index)
			if err != nil {
				return &time.Time{}, err
			}

			if tokens[*index].tType == TokenTimeFormat {
				startFormat, err = parseFormat(tokens, index)
				if err != nil {
					return &time.Time{}, err
				}
			}
		} else if tokens[*index].tType == TokenTimeFormat {
			startFormat, err = parseFormat(tokens, index)
			if err != nil {
				return &time.Time{}, err
			}
		}
	}

	if strings.ToLower(startFormat) == "am" {
		if startHour > 12 {
			return &time.Time{}, errors.New("hours in 12-hour format cannot be higher than 12")
		}
	} else if strings.ToLower(startFormat) == "pm" {
		if startHour > 12 {
			return &time.Time{}, errors.New("hours in 12-hour format cannot be higher than 12")
		}

		startHour += 12
	} else {
		if startHour > 23 {
			return &time.Time{}, errors.New("hours in 24-hour format cannot be higher than 23")
		}
	}

	if startMinute > 59 {
		return &time.Time{}, errors.New("minutes cannot be higher than 59")
	}

	startTimestamp := time.Date(0, 0, 0, startHour, startMinute, 0, 0, time.UTC)
	return &startTimestamp, nil
}

func parseDash(tokens []Token, index *int) error {
	if tokens[*index].tType != TokenDash {
		return errors.New("expected dash after start timestamp")
	}

	*index++
	return nil
}

func parseEndTimestamp(tokens []Token, index *int) (*time.Time, error) {
	endHour, err := parseTimeUnit(tokens, index)
	if err != nil {
		return &time.Time{}, err
	}

	next, err := parseSymbolAfterHour(tokens, index)
	if err != nil {
		return &time.Time{}, err
	}

	endFormat := ""
	endMinute := 0

	if next == "format" {
		endFormat, err = parseFormat(tokens, index)
		if err != nil {
			return &time.Time{}, err
		}
	} else if next == "minute" {
		if tokens[*index].tType == TokenTimeUnit {
			if len(tokens[*index].value) != 2 {
				return &time.Time{}, errors.New("minutes unit must have two digits")
			}

			endMinute, err = parseTimeUnit(tokens, index)
			if err != nil {
				return &time.Time{}, err
			}

			if tokens[*index].tType == TokenTimeFormat {
				endFormat, err = parseFormat(tokens, index)
				if err != nil {
					return &time.Time{}, err
				}
			}
		} else if tokens[*index].tType == TokenTimeFormat {
			endFormat, err = parseFormat(tokens, index)
			if err != nil {
				return &time.Time{}, err
			}
		}
	}

	if strings.ToLower(endFormat) == "am" {
		if endHour > 12 {
			return &time.Time{}, errors.New("hours in 12-hour format cannot be higher than 12")
		}
	} else if strings.ToLower(endFormat) == "pm" {
		if endHour > 12 {
			return &time.Time{}, errors.New("hours in 12-hour format cannot be higher than 12")
		}

		endHour += 12
	} else {
		if endHour > 23 {
			return &time.Time{}, errors.New("hours in 24-hour format cannot be higher than 23")
		}
	}

	if endMinute > 59 {
		return &time.Time{}, errors.New("minutes cannot be higher than 59")
	}

	endTimestamp := time.Date(0, 0, 0, endHour, endMinute, 0, 0, time.UTC)
	return &endTimestamp, nil
}

func parseTimeUnit(tokens []Token, index *int) (int, error) {
	if tokens[*index].tType != TokenTimeUnit {
		return 0, errors.New("invalid timestamp")
	}

	value, err := strconv.Atoi(tokens[*index].value)
	if err != nil {
		return 0, errors.New("invalid time unit")
	}

	*index++
	return value, nil
}

func parseSymbolAfterHour(tokens []Token, index *int) (string, error) {
	output, err := "", errors.New("invalid symbol after time unit")

	if tokens[*index].tType == TokenColon {
		output, err = "minute", nil
		*index++

		return output, err
	} else if tokens[*index].tType == TokenTimeFormat {
		output, err = "format", nil
	} else if tokens[*index].tType == TokenDash {
		output, err = "", nil
	} else if tokens[*index].tType == TokenEnd {
		output, err = "", nil
	}

	return output, err
}

func parseFormat(tokens []Token, index *int) (string, error) {
	output, err := "", errors.New("invalid time format")

	if strings.ToLower(tokens[*index].value) == "am" {
		output, err = "AM", nil
	} else if strings.ToLower(tokens[*index].value) == "pm" {
		output, err = "PM", nil
	}

	*index++
	return output, err
}