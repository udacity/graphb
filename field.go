package graphb

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Field is a recursive data struct which represents a GraphQL query field.
type Field struct {
	Name      string
	Alias     string
	Arguments []Argument
	Fields    []*Field
}

// implements fieldContainer
func (f *Field) getFields() []*Field {
	return f.Fields
}

func (f *Field) setFields(fs []*Field) {
	f.Fields = fs
}

// StringChan returns read only string token channel or an error.
// It checks if there is a circle.
func (f *Field) StringChan() (<-chan string, error) {
	if err := f.check(); err != nil {
		// return a closed channel instead of nil for receiving from nil blocks forever, hard to debug and confusing to users.
		ch := make(chan string)
		close(ch)
		return ch, errors.WithStack(err)
	}
	return f.stringChan(), nil
}

// One may have noticed that there is a public StringChan and a private stringChan.
// The different being the public method checks the validity of the Field structure
// while the private counterpart assumes the validity.
func (f *Field) stringChan() <-chan string {

	tokenChan := make(chan string)

	go func() {
		// emit alias and names
		if f.Alias != "" {
			tokenChan <- f.Alias
			tokenChan <- ":" // todo: change it to const token
		}
		tokenChan <- f.Name

		// emit argument tokens
		if len(f.Arguments) > 0 {
			str := buildArgumentString(f.Arguments)

			tokenChan <- "("
			tokenChan <- str
			tokenChan <- ")"
		}

		// emit field tokens
		if len(f.Fields) > 0 {
			tokenChan <- "{"
			for _, field := range f.Fields {
				if field != nil {
					strs := field.stringChan()
					for str := range strs {
						tokenChan <- str
					}
				}
				tokenChan <- ","
			}
			tokenChan <- "}"
		}
		close(tokenChan)
	}()
	return tokenChan
}

func (f *Field) check() error {
	if err := f.checkCycle(); err != nil {
		return errors.WithStack(err)
	}
	if err := f.checkOther(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// checkOther checks the validity of this Field and returns nil on valid Field.
func (f *Field) checkOther() error {
	// Check validity of names
	if !validName.MatchString(f.Name) {
		return errors.WithStack(InvalidNameErr{fieldName, f.Name})
	}
	if f.Alias != "" && !validName.MatchString(f.Alias) {
		return errors.WithStack(InvalidNameErr{aliasName, f.Alias})
	}
	for _, arg := range f.Arguments {
		if !validName.MatchString(arg.Name) {
			return errors.WithStack(InvalidNameErr{argumentName, arg.Name})
		}
	}

	// Check sub fields
	for _, subF := range f.Fields {
		if err := subF.checkOther(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// todo: reports the cycle path
func (f *Field) checkCycle() error {
	if err := reach(f, f); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

////////////////
// Public API //
////////////////

// MakeField constructs a Field of given name and return the pointer to it.
func MakeField(name string) *Field {
	return &Field{Name: name}
}

func (f *Field) SetArguments(arguments ...Argument) *Field {
	f.Arguments = arguments
	return f
}

func (f *Field) SetFields(fs ...*Field) *Field {
	f.Fields = fs
	return f
}

func (f *Field) SetAlias(alias string) *Field {
	f.Alias = alias
	return f
}

/////////////
// Helpers //
/////////////
// reach checks if f1 can be reached by f2 either directly (itself) or indirectly (children)
func reach(f1, f2 *Field) error {
	if f1 == nil || f2 == nil {
		return errors.WithStack(NilFieldErr{})
	}
	for _, field := range f2.Fields {
		if f1 == field {
			return errors.WithStack(CyclicFieldErr{*f1})
		}
		if err := reach(f1, field); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func strRepresentation(str string) string {
	return fmt.Sprintf(`"%s"`, str)
}

func mapStringSliceToStrRepSlice(strs []string) []string {
	newStrings := make([]string, len(strs))
	for i, s := range strs {
		newStrings[i] = strRepresentation(s)
	}
	return newStrings
}

func buildArgumentString(arguments []Argument) string {
	argumentStrings := make([]string, len(arguments))
	for i, v := range arguments {
		argumentStrings[i] = v.Name + ":" + v.Value
	}
	return strings.Join(argumentStrings, ",")
}
