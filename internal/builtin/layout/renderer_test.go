package layout

import (
	"regexp"
	"testing"

	"github.com/Masterminds/sprig"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestRenderer(t *testing.T) {
	suite.Run(t, new(rendererTestSuite))
}

// --- Suites ---

type rendererTestSuite struct {
	suite.Suite
}

func (s *rendererTestSuite) TestRenderBoxError() {
	source := `/some/file`

	box, _ := NewBox()
	renderer := NewRenderer(box, nil)

	content, err := renderer.Render(source, generator.RunValues{})
	s.Error(err)
	s.Regexp(regexp.MustCompile(`^failed to get source file "/some/file" content: open /some/file: file does not exist`), err)
	s.Nil(content)
}

func (s *rendererTestSuite) TestRenderTemplateParseError() {
	source := `/some/file`
	content := `Some text`

	mockedBox := new(mockBox)
	mockedBox.On(`ReadFile`, source).Return([]byte(content), nil)

	mockedTemplate := new(mockTemplateInterface)
	mockedTemplate.On(`Parse`, content).Return(nil, errors.New(`expected error`))

	renderer := NewRenderer(mockedBox, mockedTemplate)

	contents, err := renderer.Render(source, generator.RunValues{})
	s.Error(err)
	s.EqualError(err, `failed to parse source file "/some/file": expected error`)
	s.Nil(contents)
}

func (s *rendererTestSuite) TestRenderTemplateExecuteError() {
	source := `/some/file`
	content := `Some text`

	mockedBox := new(mockBox)
	mockedBox.On(`ReadFile`, source).Return([]byte(content), nil)

	template := NewTemplate(`embed`, nil)

	mockedTemplate := new(mockTemplateInterface)
	mockedTemplate.On(`Parse`, content).Return(template, nil)

	renderer := NewRenderer(mockedBox, mockedTemplate)

	contents, err := renderer.Render(source, generator.RunValues{})
	s.Error(err)
	s.EqualError(err, `failed to render source file "/some/file": template: embed: "embed" is an incomplete or empty template`)
	s.Nil(contents)
}

func (s *rendererTestSuite) TestRender() {
	source := `/some/file`

	mockedBox := new(mockBox)
	mockedBox.On(`ReadFile`, source).Return([]byte(`This is {{ .Application.Name }} application`), nil)

	template := NewTemplate(`embed`, sprig.TxtFuncMap())

	renderer := NewRenderer(mockedBox, template)

	contents, err := renderer.Render(source, generator.RunValues{Application: generator.RunValuesApplication{Name: `testing`}})
	s.NoError(err)
	s.Equal(`This is testing application`, string(contents))
}
