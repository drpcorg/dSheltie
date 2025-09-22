package flow_test

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/drpcorg/nodecore/internal/protocol"
	"github.com/drpcorg/nodecore/internal/upstreams/flow"
	"github.com/drpcorg/nodecore/pkg/test_utils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMultiMatcher(t *testing.T) {
	method := flow.NewMethodMatcher("eth_getBalance")
	status := flow.NewStatusMatcher()
	multiMatcher := flow.NewMultiMatcher(method, status)
	methods := mocks.NewMethodsMock()
	methods.On("HasMethod", "eth_getBalance").Return(true)
	state := protocol.UpstreamState{Status: protocol.Available, UpstreamMethods: methods}

	resp := multiMatcher.Match("1", &state)

	assert.Equal(t, flow.SuccessResponse{}, resp)
	assert.Equal(t, flow.SuccessType, resp.Type())
}

func TestMultiMatcherResponseMethodType(t *testing.T) {
	method := flow.NewMethodMatcher("no-method")
	status := flow.NewStatusMatcher()
	multiMatcher := flow.NewMultiMatcher(method, status)
	methods := mocks.NewMethodsMock()
	methods.On("HasMethod", "no-method").Return(false)
	state := protocol.UpstreamState{Status: protocol.Unavailable, UpstreamMethods: methods}

	resp := multiMatcher.Match("1", &state)

	assert.IsType(t, flow.MethodResponse{}, resp)
	assert.Equal(t, flow.MethodType, resp.Type())
	assert.Equal(t, "method no-method is not supported", resp.Cause())
}

func TestMethodMatcher(t *testing.T) {
	matcher := flow.NewMethodMatcher("eth_getBalance")
	methods := mocks.NewMethodsMock()
	methods.On("HasMethod", "eth_getBalance").Return(true)
	state := protocol.UpstreamState{UpstreamMethods: methods}

	resp := matcher.Match("1", &state)

	assert.Equal(t, flow.SuccessResponse{}, resp)
	assert.Equal(t, flow.SuccessType, resp.Type())
}

func TestMethodMatcherNoMethod(t *testing.T) {
	matcher := flow.NewMethodMatcher("no-method")
	methods := mocks.NewMethodsMock()
	methods.On("HasMethod", "no-method").Return(false)
	state := protocol.UpstreamState{UpstreamMethods: methods}

	resp := matcher.Match("1", &state)

	assert.IsType(t, flow.MethodResponse{}, resp)
	assert.Equal(t, flow.MethodType, resp.Type())
	assert.Equal(t, "method no-method is not supported", resp.Cause())
}

func TestStatusMatcher(t *testing.T) {
	matcher := flow.NewStatusMatcher()
	state := protocol.UpstreamState{Status: protocol.Available}

	resp := matcher.Match("1", &state)

	assert.Equal(t, flow.SuccessResponse{}, resp)
	assert.Equal(t, flow.SuccessType, resp.Type())
}

func TestStatusMatcherNotAvailable(t *testing.T) {
	matcher := flow.NewStatusMatcher()
	state := protocol.UpstreamState{Status: protocol.Unavailable}

	resp := matcher.Match("1", &state)

	assert.Equal(t, flow.AvailabilityResponse{}, resp)
	assert.Equal(t, "upstream is not available", resp.Cause())
	assert.Equal(t, flow.AvailabilityType, resp.Type())
}

func TestWsCapMatcher(t *testing.T) {
	matcher := flow.NewWsCapMatcher("sub")
	state := protocol.UpstreamState{Caps: mapset.NewThreadUnsafeSet[protocol.Cap](protocol.WsCap)}

	resp := matcher.Match("1", &state)

	assert.Equal(t, flow.SuccessResponse{}, resp)
	assert.Equal(t, flow.SuccessType, resp.Type())
}

func TestWsCapMatcherNotAvailable(t *testing.T) {
	matcher := flow.NewWsCapMatcher("sub")
	state := protocol.UpstreamState{Caps: mapset.NewThreadUnsafeSet[protocol.Cap]()}

	resp := matcher.Match("1", &state)

	assert.IsType(t, flow.MethodResponse{}, resp)
	assert.Equal(t, flow.MethodType, resp.Type())
	assert.Equal(t, "method sub is not supported", resp.Cause())
}

func TestUpstreamIndexMatcher(t *testing.T) {
	matcher := flow.NewUpstreamIndexMatcher("index")
	state := protocol.UpstreamState{UpstreamIndex: "index"}

	resp := matcher.Match("1", &state)

	assert.Equal(t, flow.SuccessResponse{}, resp)
	assert.Equal(t, flow.SuccessType, resp.Type())
}

func TestUpstreamIndexMatcherNotExist(t *testing.T) {
	matcher := flow.NewUpstreamIndexMatcher("not-exist")
	state := protocol.UpstreamState{UpstreamIndex: "index"}

	resp := matcher.Match("1", &state)

	assert.IsType(t, flow.UpstreamIndexResponse{}, resp)
	assert.Equal(t, flow.UpstreamIndexType, resp.Type())
	assert.Equal(t, "no upstream with index not-exist", resp.Cause())
}
