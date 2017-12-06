package bazel

import (
	"fmt"
	"log"
	"os"
	"path"
	"io/ioutil"
	"os/exec"
	"github.com/golang/protobuf/proto"
	"github.com/matttproud/golang_protobuf_extensions/pbutil"
	"github.com/bzl-io/bzl/config"
	build  "github.com/bzl-io/bzl/proto/build_go"
	stream "github.com/bzl-io/bzl/proto/build_event_stream_go"
)

// Default bazel version is whatever is currently in the users' path.
var bazel = "bazel"
var home = ""

type Bazel struct {
	Name string
}

func New() *Bazel {
	return &Bazel{
		Name: bazel,
	}
}

// Set the version of bazel to use.  Given '0.7.0', this looks for .cache/bzl/release/0.7.0/bin/bazel.
func SetVersion(version string) error {
	home, err := config.GetHome()
	if err != nil {
		return err
	}
	
	exe := path.Join(home, "release", version, "bin", "bazel")
	if _, err := os.Stat(exe); os.IsNotExist(err) {
		log.Printf("Error: bazel %s does not exist in the release cache.  Try 'bzl install %s' first.", version, version)
		return err
	}
	log.Printf("Setting bazel version: %s\n", bazel)
	bazel = exe
	return nil
}

// Make Generic invocation to bazel
func (b *Bazel) Invoke(args []string) error {
	fmt.Printf("\n%s %v\n", b.Name, args)
	cmd := exec.Command(b.Name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = ""
	return cmd.Run() 
}


// Make Invocation to bazel and get back the event graph
func (b *Bazel) InvokeWithEvents(args []string) ([]*stream.BuildEvent, error) {
	fmt.Printf("\n%s %v\n", b.Name, args)
	file, err := ioutil.TempFile("/tmp", "bes-")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	args = append(args, "--build_event_binary_file", file.Name())
	cmd := exec.Command(b.Name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = ""
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	events, err := b.readBuildEventStream(file.Name())
	if err != nil {
		return nil, err
	}
	return events, nil
}

// Do a query invocation and get the query result proto back
func (b *Bazel) Query(pattern string) (*build.QueryResult, error) {
	var (
		cmdOut []byte
		err    error
	)
	args := []string{
		"query", pattern,
		"--output", "proto",
	}

	cmd := exec.Command(b.Name, args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Query Error: ", string(out), err, "\n")
		return nil, err
	}

	build := &build.QueryResult{}
	err = proto.Unmarshal(out, build)
	if err != nil {
		fmt.Printf("Query Error: ", string(cmdOut), err, "\n")
		
		return nil, err
	}

	return build, nil
}


func (b *Bazel) readBuildEventStream(filename string) ([]*stream.BuildEvent, error) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to read <%s>: %s\n", filename, err)
		return nil, err
	}
	defer f.Close()

	events := make([]*stream.BuildEvent, 0)
	for {
		event := &stream.BuildEvent{}
		remaining, err := pbutil.ReadDelimited(f, event)
		if remaining == 0 {
			return events, nil
		}
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// From a list of BuildEvents, return the first one typed as
// 'Completed'.  Anecdotally, there is only one per bazel invocation.
// Pointer will be nil if none found.
func FirstTargetComplete(events []*stream.BuildEvent) *stream.TargetComplete {
	for _, event := range events {
		switch event.Payload.(type) {
		case *stream.BuildEvent_Completed:
			return event.GetCompleted();
		}
	}
	return nil
}
