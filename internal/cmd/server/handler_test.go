package server

import (
	"fmt"
	"slices"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"

	resourcev1 "github.com/otterscale/otterscale/api/resource/v1"
	runtimev1 "github.com/otterscale/otterscale/api/runtime/v1"
)

// TestLongRunningPathsCoversEveryStream derives the expected set from
// the proto descriptors instead of restating it, so adding a
// server-streaming RPC without exempting it from the request timeouts
// fails here rather than as a stream that dies after five minutes.
func TestLongRunningPathsCoversEveryStream(t *testing.T) {
	var want []string
	for _, file := range []protoreflect.FileDescriptor{
		resourcev1.File_resource_v1_resource_proto,
		runtimev1.File_runtime_v1_runtime_proto,
	} {
		services := file.Services()
		for i := range services.Len() {
			service := services.Get(i)
			methods := service.Methods()
			for j := range methods.Len() {
				method := methods.Get(j)
				if method.IsStreamingServer() {
					want = append(want, fmt.Sprintf("/%s/%s", service.FullName(), method.Name()))
				}
			}
		}
	}

	if len(want) == 0 {
		t.Fatal("no server-streaming methods found; the descriptors are probably not what this test assumes")
	}

	got := (&Handler{}).LongRunningPaths()
	slices.Sort(got)
	slices.Sort(want)

	if !slices.Equal(got, want) {
		t.Errorf("LongRunningPaths() =\n\t%q\nwant\n\t%q", got, want)
	}
}
