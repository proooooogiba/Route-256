package sql_parser

import (
	"regexp"
	"strings"
)

const ( // iota is reset to 0
	SELECT = iota
	INSERT
	INTO
	Identificator
	Template
	OpenBrackets
	Assign
	STAR
	CloseBrackets
	Comma
	Semicolon
	FROM
	WHERE
	AND
	VALUES
	ERROR
)

type Lex int
type slice []string

func containsOnlyChars(word string) bool {
	pattern := "a-zA-Z"
	matched, _ := regexp.MatchString("^["+pattern+"]+$", word)
	return matched
}

func isTemplate(expression string) bool {
	regex := regexp.MustCompile(`^\$\d+$`)
	return regex.MatchString(expression)
}

func getLex(iterator *SliceIterator) Lex {
	iterator.Next()
	word := iterator.Value()
	switch word {
	case "INSERT":
		return INSERT
	case "SELECT":
		return SELECT
	case "INTO":
		return INTO
	case "FROM":
		return FROM
	case "WHERE":
		return WHERE
	case "AND":
		return AND
	case "VALUES":
		return VALUES
	case "*":
		return STAR
	case "=":
		return Assign
	case ",":
		return Comma
	case "(":
		return OpenBrackets
	case ")":
		return CloseBrackets
	case ";":
		return Semicolon
	}

	if containsOnlyChars(word) {
		return Identificator
	}

	if isTemplate(word) {
		return Template
	}

	return ERROR
}

func assign(lex Lex, iterator *SliceIterator) (Lex, []string, []string) {
	var identSlice, templateSlice []string

	for lex == Identificator {
		identSlice = append(identSlice, iterator.Value())
		lex = getLex(iterator)
		if lex == Assign {
			lex = getLex(iterator)
		} else {
			return ERROR, nil, nil
		}

		if lex == Template {
			templateSlice = append(templateSlice, iterator.Value())
			lex = getLex(iterator)
		} else {
			return ERROR, nil, nil
		}

		if lex == AND {
			lex = getLex(iterator)
			continue
		} else if lex == Semicolon {
			break
		} else {
			return ERROR, nil, nil
		}
	}
	return lex, identSlice, templateSlice
}

func varDeclaration(lex Lex, iterator *SliceIterator) (Lex, []string) {
	var identSlice []string
	if lex == OpenBrackets {
		lex = getLex(iterator)
	} else {
		return ERROR, nil
	}
	for lex == Identificator {
		identSlice = append(identSlice, iterator.Value())
		lex = getLex(iterator)
		if lex == Comma {
			lex = getLex(iterator)
		} else if lex == CloseBrackets {
			lex = getLex(iterator)
			break
		} else {
			return ERROR, nil
		}
	}
	return lex, identSlice
}

func templateDeclaration(lex Lex, iterator *SliceIterator) (Lex, []string) {
	var templateSlice []string
	if lex == OpenBrackets {
		lex = getLex(iterator)
	} else {
		return ERROR, nil
	}
	for lex == Template {
		templateSlice = append(templateSlice, iterator.Value())
		lex = getLex(iterator)
		if lex == Comma {
			lex = getLex(iterator)
		} else if lex == CloseBrackets {
			lex = getLex(iterator)
			break
		} else {
			return ERROR, nil
		}
	}
	return lex, templateSlice
}

func insertQuery(iterator *SliceIterator) (QueryInfo, error) {
	var lex Lex
	var queryInfo QueryInfo
	queryInfo.Operation = "INSERT"
	lex = getLex(iterator)

	if lex == INTO {
		lex = getLex(iterator)
	} else {
		return QueryInfo{}, ParseSQLError
	}

	if lex == Identificator {
		queryInfo.Table = iterator.Value()
		lex = getLex(iterator)
	} else {
		return QueryInfo{}, ParseSQLError
	}

	lex, identSlice := varDeclaration(lex, iterator)
	if lex == ERROR {
		return QueryInfo{}, ParseSQLError
	}
	queryInfo.IdentSlice = identSlice

	if lex == VALUES {
		lex = getLex(iterator)
	} else {
		return QueryInfo{}, ParseSQLError
	}

	lex, templateSlice := templateDeclaration(lex, iterator)
	if lex == ERROR {
		return QueryInfo{}, ParseSQLError
	}
	queryInfo.TemplateSlice = templateSlice

	if len(identSlice) != len(templateSlice) {
		return QueryInfo{}, ParseSQLError
	}

	if lex == Semicolon {
		return queryInfo, nil
	} else {
		return QueryInfo{}, ParseSQLError
	}
}

func selectQuery(iterator *SliceIterator) (QueryInfo, error) {
	var lex Lex
	var queryInfo QueryInfo
	queryInfo.Operation = "SELECT"
	lex = getLex(iterator)

	if lex == STAR {
		lex = getLex(iterator)
	} else {
		return QueryInfo{}, ParseSQLError
	}

	if lex == FROM {
		lex = getLex(iterator)
	} else {
		return QueryInfo{}, ParseSQLError
	}

	if lex == Identificator {
		queryInfo.Table = iterator.Value()
		lex = getLex(iterator)
	} else {
		return QueryInfo{}, ParseSQLError
	}

	if lex == WHERE {
		lex = getLex(iterator)
	} else {
		return QueryInfo{}, ParseSQLError
	}

	lex, identSlice, templateSlice := assign(lex, iterator)
	if lex == ERROR {
		return QueryInfo{}, ParseSQLError
	}

	queryInfo.IdentSlice = identSlice
	queryInfo.TemplateSlice = templateSlice

	if lex == Semicolon {
		return queryInfo, nil
	} else {
		return QueryInfo{}, ParseSQLError
	}
}

type SliceIterator struct {
	slice []string
	index int
}

func NewSliceIterator(slice []string) *SliceIterator {
	return &SliceIterator{
		slice: slice,
		index: -1,
	}
}

func (s *SliceIterator) Next() bool {
	s.index++
	return s.index < len(s.slice)
}

func (s *SliceIterator) Value() string {
	return s.slice[s.index]
}

type QueryInfo struct {
	Operation     string
	Table         string
	IdentSlice    slice
	TemplateSlice slice
}

func (s slice) Equals(other slice) bool {
	if len(s) != len(other) {
		return false
	}
	for i, _ := range s {
		if s[i] != other[i] {
			return false
		}
	}
	return true
}

func (q QueryInfo) Equals(other QueryInfo) bool {
	return q.Operation == other.Operation &&
		q.Table == other.Table &&
		q.IdentSlice.Equals(other.IdentSlice) &&
		q.TemplateSlice.Equals(other.TemplateSlice)
}

func Parse(sql string) (QueryInfo, error) {
	query := strings.Split(sql, " ")
	iterator := NewSliceIterator(query)

	switch getLex(iterator) {
	case INSERT:
		queryInfo, err := insertQuery(iterator)
		if err != nil {
			return QueryInfo{}, err
		}
		return queryInfo, nil
	case SELECT:
		queryInfo, err := selectQuery(iterator)
		if err != nil {
			return QueryInfo{}, err
		}
		return queryInfo, nil
	default:
		return QueryInfo{}, ParseSQLError
	}
	return QueryInfo{}, nil
}
