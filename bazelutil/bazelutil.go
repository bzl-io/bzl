package bazelutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"

	"github.com/bzl-io/bzl/config"
	"github.com/golang/protobuf/proto"
	"github.com/matttproud/golang_protobuf_extensions/pbutil"

	bes "github.com/bzl-io/bzl/proto/bes"
	build "github.com/bzl-io/bzl/proto/build"
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

// Set the version of bazel to use.  Given '0.7.0', this looks for .cache/bzl/release/0.7.0/bin/bazelutil.
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
	// log.Printf("Setting bazel version: %s\n", bazel)
	bazel = exe
	return nil
}

// Make Generic invocation to bazel
func (b *Bazel) Invoke(args []string, dir string) (error, int) {

	//fmt.Printf("\n%s %v\n", b.Name, args)
	cmd := exec.Command(b.Name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir

	err := cmd.Run()

	var exitCode int

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH, in
			// this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and
			// format err to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", b.Name, args)
			exitCode = -1
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	return err, exitCode
}

// Make Invocation to bazel and get back the event graph
func (b *Bazel) InvokeWithEvents(args []string) ([]*bes.BuildEvent, error) {
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
		fmt.Printf("Query Error %q: %v\n", string(out), err)
		return nil, err
	}

	build := &build.QueryResult{}
	err = proto.Unmarshal(out, build)
	if err != nil {
		fmt.Printf("Query Error %q: %v\n", string(cmdOut), err)

		return nil, err
	}

	return build, nil
}

func (b *Bazel) readBuildEventStream(filename string) ([]*bes.BuildEvent, error) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to read <%s>: %s\n", filename, err)
		return nil, err
	}
	defer f.Close()

	events := make([]*bes.BuildEvent, 0)
	for {
		event := &bes.BuildEvent{}
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
func FirstTargetComplete(events []*bes.BuildEvent) *bes.TargetComplete {
	for _, event := range events {
		switch event.Payload.(type) {
		case *bes.BuildEvent_Completed:
			return event.GetCompleted()
		}
	}
	return nil
}

// https://gist.github.com/elazarl/5507969
func CopyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
