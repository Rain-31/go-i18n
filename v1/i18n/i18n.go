package i18n

import (
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var p printers

// PluralRule is Plural rule
type PluralRule struct {
	Pos   int
	Expr  string
	Value int
	Text  string
}

type printers struct {
	sessionMap sync.Map
}

type PrinterSession struct {
	printer *message.Printer
}

func RegistPrinter(id string, lang language.Tag) *PrinterSession {
	session := &PrinterSession{printer: message.NewPrinter(lang)}
	p.sessionMap.Store(id, session)

	return session
}

func DeletePrinter(id string) {
	p.sessionMap.Delete(id)
}

//Session load session with target id
func Session(id string) *PrinterSession {
	p.sessionMap.Load(id)
	if session, exist := p.sessionMap.Load(id); exist == true {
		return session.(*PrinterSession)
	} else {
		return &PrinterSession{printer: message.NewPrinter(language.AmericanEnglish)}
	}
}

// Printf is like fmt.Printf, but using language-specific formatting.
func (s *PrinterSession) Printf(format string, args ...interface{}) (n int, err error) {
	format, args = preArgs(format, args...)
	return s.printer.Printf(format, args...)
}

// Sprintf is like fmt.Sprintf, but using language-specific formatting.
func (s *PrinterSession) Sprintf(format string, args ...interface{}) string {
	format, args = preArgs(format, args...)
	return s.printer.Sprintf(format, args...)
}

// Fprintf is like fmt.Fprintf, but using language-specific formatting.
func (s *PrinterSession) Fprintf(w io.Writer, key message.Reference, a ...interface{}) (n int, err error) {
	format, args := preArgs(key.(string), a...)
	key = message.Reference(format)
	return s.printer.Fprintf(w, key, args...)
}

func Printf(id string, format string, args ...interface{}) (n int, err error) {
	return Session(id).Printf(format, args...)
}

func Sprintf(id string, format string, args ...interface{}) string {
	return Session(id).Sprintf(format, args...)
}

func Fprintf(id string, w io.Writer, key message.Reference, a ...interface{}) (n int, err error) {
	return Session(id).Fprintf(w, key, a...)
}

// Preprocessing parameters in plural form
func preArgs(format string, args ...interface{}) (string, []interface{}) {
	length := len(args)
	if length > 0 {
		lastArg := args[length-1]
		switch lastArg.(type) {
		case []PluralRule:
			rules := lastArg.([]PluralRule)
			// parse rule
			for _, rule := range rules {
				curPosVal := args[rule.Pos-1].(int)
				// Support comparison expression
				if (rule.Expr == "=" && curPosVal == rule.Value) || (rule.Expr == ">" && curPosVal > rule.Value) {
					format = rule.Text
					break
				}
			}
			args = args[0:strings.Count(format, "%")]
		}
	}
	return format, args
}

// Plural is Plural function
func Plural(cases ...interface{}) []PluralRule {
	rules := []PluralRule{}
	// %[1]d=1, %[1]d>1
	re := regexp.MustCompile(`\[(\d+)\][^=>]\s*(\=|\>)\s*(\d+)$`)
	for i := 0; i < len(cases); {
		expr := cases[i].(string)
		if i++; i >= len(cases) {
			return rules
		}
		text := cases[i].(string)
		// cannot match continue
		if !re.MatchString(expr) {
			continue
		}
		matches := re.FindStringSubmatch(expr)
		pos, _ := strconv.Atoi(matches[1])
		value, _ := strconv.Atoi(matches[3])
		rules = append(rules, PluralRule{
			Pos:   pos,
			Expr:  matches[2],
			Value: value,
			Text:  text,
		})
		i++
	}
	return rules
}
