package plugin

import (
	"context"

	pluginSDK "github.com/hashicorp/go-plugin"

	"github.com/insidieux/inizio/pkg/api/protobuf"
	"github.com/insidieux/inizio/pkg/sdk/generator"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// --- Mocks ---

// mockGeneratorClient is an autogenerated mock type for the generatorClient type
type mockGeneratorClient struct {
	mock.Mock
}

// Run provides a mock function with given fields: ctx, in, opts
func (_m *mockGeneratorClient) Run(ctx context.Context, in *protobuf.Run_Request, opts ...grpc.CallOption) (*protobuf.Run_Response, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *protobuf.Run_Response
	if rf, ok := ret.Get(0).(func(context.Context, *protobuf.Run_Request, ...grpc.CallOption) *protobuf.Run_Response); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*protobuf.Run_Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *protobuf.Run_Request, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockGenerator is an autogenerated mock type for the Generator type
type mockGenerator struct {
	mock.Mock
}

// Run provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockGenerator) Run(_a0 context.Context, _a1 generator.RunOptions, _a2 generator.RunValues) (generator.RunResult, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 generator.RunResult
	if rf, ok := ret.Get(0).(func(context.Context, generator.RunOptions, generator.RunValues) generator.RunResult); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(generator.RunResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, generator.RunOptions, generator.RunValues) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockSdkClient is an autogenerated mock type for the sdkClient type
type mockSdkClient struct {
	mock.Mock
}

// Client provides a mock function with given fields:
func (_m *mockSdkClient) Client() (pluginSDK.ClientProtocol, error) {
	ret := _m.Called()

	var r0 pluginSDK.ClientProtocol
	if rf, ok := ret.Get(0).(func() pluginSDK.ClientProtocol); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pluginSDK.ClientProtocol)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockClientProtocol is an autogenerated mock type for the pluginSDK.ClientProtocol type
type mockClientProtocol struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *mockClientProtocol) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Dispense provides a mock function with given fields: _a0
func (_m *mockClientProtocol) Dispense(_a0 string) (interface{}, error) {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ping provides a mock function with given fields:
func (_m *mockClientProtocol) Ping() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
