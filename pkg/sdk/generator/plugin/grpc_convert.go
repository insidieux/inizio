package plugin

import (
	"github.com/insidieux/inizio/pkg/api/protobuf"
	"github.com/insidieux/inizio/pkg/sdk/generator"
)

func generatorRunArgumentsToProto(options generator.RunOptions, values generator.RunValues) *protobuf.Run_Request {
	return &protobuf.Run_Request{
		Options: &protobuf.Run_Request_Options{
			Override:         options.Override,
			WorkingDirectory: options.WorkingDirectory,
		},
		Values: &protobuf.Run_Request_Values{
			Application: &protobuf.Run_Request_Values_Application{
				Name:        values.Application.Name,
				Description: values.Application.Description,
			},
			Golang: &protobuf.Run_Request_Values_Golang{
				Module:  values.Golang.Module,
				Version: values.Golang.Version,
			},
		},
	}
}

func generatorRunResultToProto(result generator.RunResult) *protobuf.Run_Response {
	return &protobuf.Run_Response{Generated: result}
}

func protoRunRequestToOptions(request *protobuf.Run_Request) generator.RunOptions {
	options := generator.RunOptions{
		Override:         request.Options.Override,
		WorkingDirectory: request.Options.WorkingDirectory,
	}
	return options
}

func protoRunRequestToValues(request *protobuf.Run_Request) generator.RunValues {
	values := generator.RunValues{
		Application: generator.RunValuesApplication{
			Name:        request.Values.Application.Name,
			Description: request.Values.Application.Description,
		},
		Golang: generator.RunValuesGolang{
			Module:  request.Values.Golang.Module,
			Version: request.Values.Golang.Version,
		},
	}
	return values
}

func protoRunResponseToResult(response *protobuf.Run_Response) generator.RunResult {
	return response.Generated
}
