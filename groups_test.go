package routem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testPathTwo = "test/path/two"
)

func assertDefaultGroup(t *testing.T, group Group) {
	assertDefaultConfig(t, group)
	assertGroup(t, group)
}

func assertTestGroup(t *testing.T, group Group) {
	assertTestConfig(t, group)
	assertGroup(t, group)
}

func assertGroupRelation(t *testing.T, group, subGroup Group) {
	assert.Equal(t, 0, len(subGroup.Routes()), "Wrong subGroup route count")
	assert.Equal(t, 1, len(group.Routes()), "Wrong route count")
}

func assertGroup(t *testing.T, group Group) {
	assert.Equal(t, testPath, group.Path(), "Wrong path")
}

func TestNewGroup(t *testing.T) {
	group := newGroup(defaultConfig(), testPath)

	assertDefaultGroup(t, group)

	group = newGroup(testConfig(), testPath)

	assertTestGroup(t, group)
}

func TestGroupWith(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.With(GetMethod, testPath, testHandler)

	assertTestRoute(t, route)
}

func TestGroupWithHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}

	route := group.WithHTTP(GetMethod, testPath, testHandler)

	assertTestHTTPRoute(t, route, testHandler)
}

func TestTestGroupWithSubGroup(t *testing.T) {
	group := newGroup(testConfig(), testPathTwo)

	subGroup := group.WithGroup(testPath)

	assertTestGroup(t, subGroup)
	assertGroupRelation(t, group, subGroup)
}

func TestDefaultGroupWithSubgroup(t *testing.T) {
	group := newGroup(defaultConfig(), testPathTwo)

	subGroup := group.WithGroup(testPath)

	assertDefaultGroup(t, subGroup)
	assertGroupRelation(t, group, subGroup)
}

func TestNoop(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Noop(testPath, testHandler)
	assertRouteWithMethods(t, route, NoMethod)
}

func TestConnect(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Connect(testPath, testHandler)
	assertRouteWithMethods(t, route, ConnectMethod)
}

func TestDelete(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Delete(testPath, testHandler)
	assertRouteWithMethods(t, route, DeleteMethod)
}

func TestGet(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Get(testPath, testHandler)
	assertRouteWithMethods(t, route, GetMethod)
}

func TestHead(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Head(testPath, testHandler)
	assertRouteWithMethods(t, route, HeadMethod)
}

func TestOptions(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Options(testPath, testHandler)
	assertRouteWithMethods(t, route, OptionsMethod)
}

func TestPatch(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Patch(testPath, testHandler)
	assertRouteWithMethods(t, route, PatchMethod)
}

func TestPut(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Put(testPath, testHandler)
	assertRouteWithMethods(t, route, PutMethod)
}

func TestPost(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Post(testPath, testHandler)
	assertRouteWithMethods(t, route, PostMethod)
}

func TestTrace(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Trace(testPath, testHandler)
	assertRouteWithMethods(t, route, TraceMethod)
}

func TestCrud(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Crud(testPath, testHandler)
	assertRouteWithMethods(t, route, CrudMethod)
}

func TestAny(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	route := group.Any(testPath, testHandler)
	assertRouteWithMethods(t, route, AnyMethod)
}

func TestNoopHttp(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.NoopHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, NoMethod)
}

func TestConnectHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.ConnectHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, ConnectMethod)
}

func TestDeleteHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.DeleteHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, DeleteMethod)
}

func TestGetHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.GetHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, GetMethod)
}

func TestHeadHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.HeadHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, HeadMethod)
}

func TestOptionsHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.OptionsHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, OptionsMethod)
}

func TestPatchHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.PatchHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, PatchMethod)
}

func TestPutHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.PutHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, PutMethod)
}

func TestPostHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.PostHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, PostMethod)
}

func TestTraceHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.TraceHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, TraceMethod)
}

func TestCrudHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.CrudHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, CrudMethod)
}

func TestAnyHTTP(t *testing.T) {
	group := newGroup(testConfig(), testPath)

	var called bool
	testHandler := testHTTPHandler{
		called: &called,
	}
	route := group.AnyHTTP(testPath, testHandler)

	assertTestHTTPRouteWithMethods(t, route, testHandler, AnyMethod)
}
