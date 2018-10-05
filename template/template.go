package template

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	texttemplate "text/template"
)

// Vars is a loose type alias to map[string]interface{}
type Vars = map[string]interface{}

// New creates a new template.
func New() *Template {
	temp := &Template{
		vars: Vars{},
		env:  ParseEnvVars(os.Environ()),
	}
	temp.funcs = ViewFuncs.FuncMap()
	return temp
}

// NewFromFile creates a new template from a file.
func NewFromFile(filepath string) (*Template, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return New().WithName(filepath).WithBody(string(contents)), nil
}

// Template is a wrapper for html.Template.
type Template struct {
	name       string
	body       string
	vars       Vars
	env        map[string]string
	includes   []string
	funcs      texttemplate.FuncMap
	leftDelim  string
	rightDelim string
}

// WithName sets the template name.
func (t *Template) WithName(name string) *Template {
	t.name = name
	return t
}

// Name returns the template name if set, or if not set, just "template" as a constant.
func (t *Template) Name() string {
	if len(t.name) > 0 {
		return t.name
	}
	return "template"
}

// WithDelims sets the template action delimiters, treating empty string as default delimiter.
func (t *Template) WithDelims(left, right string) *Template {
	t.leftDelim = left
	t.rightDelim = right
	return t
}

// WithBody sets the template body and returns a reference to the template object.
func (t *Template) WithBody(body string) *Template {
	t.body = body
	return t
}

// WithInclude includes a (sub) template into the rendering assets.
func (t *Template) WithInclude(body string) *Template {
	t.includes = append(t.includes, body)
	return t
}

// Body returns the template body.
func (t *Template) Body() string {
	return t.body
}

// WithVar sets a variable and returns a reference to the template object.
func (t *Template) WithVar(key string, value interface{}) *Template {
	t.SetVar(key, value)
	return t
}

// WithVars reads a map of variables into the template.
func (t *Template) WithVars(vars Vars) *Template {
	for key, value := range vars {
		t.SetVar(key, value)
	}
	return t
}

// SetVar sets a var in the template.
func (t *Template) SetVar(key string, value interface{}) {
	t.vars[key] = value
}

// HasVar returns if a variable is set.
func (t *Template) HasVar(key string) bool {
	_, hasKey := t.vars[key]
	return hasKey
}

// Var returns the value of a variable, or panics if the variable is not set.
func (t *Template) Var(key string, defaults ...interface{}) (interface{}, error) {
	if value, hasVar := t.vars[key]; hasVar {
		return value, nil
	}

	if len(defaults) > 0 {
		return defaults[0], nil
	}

	return nil, fmt.Errorf("template variable `%s` is unset and no default is provided", key)
}

// Env returns an environment variable.
func (t *Template) Env(key string, defaults ...string) (string, error) {
	if value, hasVar := t.env[key]; hasVar {
		return value, nil
	}

	if len(defaults) > 0 {
		return defaults[0], nil
	}

	return "", fmt.Errorf("template env variable `%s` is unset and no default is provided", key)
}

// HasEnv returns if an env var is set.
func (t *Template) HasEnv(key string) bool {
	_, hasKey := t.env[key]
	return hasKey
}

// Process processes the template.
func (t *Template) Process(dst io.Writer) error {
	base := texttemplate.New(t.Name()).Funcs(t.ViewFuncs()).Delims(t.leftDelim, t.rightDelim)

	var err error
	for _, include := range t.includes {
		_, err = base.New(t.Name()).Parse(include)
		if err != nil {
			return err
		}
	}

	final, err := base.New(t.Name()).Parse(t.body)
	if err != nil {
		return err
	}
	return final.Execute(dst, t)
}

// ViewFuncs returns the view funcs.
func (t *Template) ViewFuncs() texttemplate.FuncMap {
	return t.funcs
}
