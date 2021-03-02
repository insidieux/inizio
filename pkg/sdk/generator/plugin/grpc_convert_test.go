package plugin

import (
	"testing"

	"github.com/insidieux/inizio/pkg/api/protobuf"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/stretchr/testify/suite"
)

// --- Tests ---

func TestConvert(t *testing.T) {
	suite.Run(t, new(convertTestSuite))
}

// --- Suites ---

type convertTestSuite struct {
	suite.Suite
}

func (s *convertTestSuite) TestGeneratorRunArgumentsToProto() {
	var (
		options = generator.RunOptions{
			Override:         false,
			WorkingDirectory: `/dir`,
		}
		values = generator.RunValues{
			Application: generator.RunValuesApplication{},
			Golang:      generator.RunValuesGolang{},
		}
		expected = &protobuf.Run_Request{
			Options: &protobuf.Run_Request_Options{
				Override:         false,
				WorkingDirectory: `/dir`,
			},
			Values: &protobuf.Run_Request_Values{
				Application: &protobuf.Run_Request_Values_Application{},
				Golang:      &protobuf.Run_Request_Values_Golang{},
			},
		}
	)
	got := generatorRunArgumentsToProto(options, values)
	s.Equal(expected, got)
}

func (s *convertTestSuite) TestGeneratorRunResultToProto() {
	var (
		list     = []string{`/some/path`}
		options  = generator.RunResult(list)
		expected = &protobuf.Run_Response{
			Generated: list,
		}
	)
	got := generatorRunResultToProto(options)
	s.Equal(expected, got)
}

func (s *convertTestSuite) TestProtoRunRequestToOptions() {
	var (
		request = &protobuf.Run_Request{
			Options: &protobuf.Run_Request_Options{
				Override:         false,
				WorkingDirectory: `/dir`,
			},
		}
		expected = generator.RunOptions{
			Override:         false,
			WorkingDirectory: `/dir`,
		}
	)
	got := protoRunRequestToOptions(request)
	s.Equal(expected, got)
}

func (s *convertTestSuite) TestProtoRunRequestToValues() {
	var (
		request = &protobuf.Run_Request{
			Values: &protobuf.Run_Request_Values{
				Application: &protobuf.Run_Request_Values_Application{
					Name:        `application`,
					Description: `description`,
				},
				Golang: &protobuf.Run_Request_Values_Golang{
					Module:  `github.com/username/project`,
					Version: `1.15.6`,
				},
			},
		}
		expected = generator.RunValues{
			Application: generator.RunValuesApplication{
				Name:        `application`,
				Description: `description`,
			},
			Golang: generator.RunValuesGolang{
				Module:  `github.com/username/project`,
				Version: `1.15.6`,
			},
		}
	)
	got := protoRunRequestToValues(request)
	s.Equal(expected, got)
}

func (s *convertTestSuite) TestProtoRunResponseToResult() {
	var (
		list    = []string{`/some/path`}
		request = &protobuf.Run_Response{
			Generated: list,
		}
		expected = generator.RunResult(list)
	)
	got := protoRunResponseToResult(request)
	s.Equal(expected, got)
}
