package install

import (
	"bytes"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"time"	
)

type InstallError struct {
	Message string
	ExitCode int
}

const ghUrl = "github.com/bazelbuild/bazel"

var (
	DOWNLOAD_ERROR = InstallError{ "Failed to download resource <%s>", 1 }
)


func Execute(c *cli.Context) error {
	version := c.Args().First()
	fmt.Printf("INSTALL %s\n", version)
	os := "linux"
	arch := "x86_64"
	// https://github.com/bazelbuild/bazel/releases/download/0.7.0/bazel-0.7.0-installer-darwin-x86_64.sh.sha256
	shaUrl := fmt.Sprintf("https://%s/releases/download/%s/bazel-%s-installer-%s-%s.sh.sha256", ghUrl, version, version, os, arch)
	url := fmt.Sprintf("https://%s/releases/download/%s/bazel-%s-installer-%s-%s.sh", ghUrl, version, version, os, arch)

	err := downloadFile(url, "/tmp")
	if err != nil {
		return cli.NewExitError(fmt.Sprintf(DOWNLOAD_ERROR.Message + ".\nUnderlying cause: %v", url, err), DOWNLOAD_ERROR.ExitCode)
	}
	err = downloadFile(shaUrl, "/tmp")
	if err != nil {
		return cli.NewExitError(fmt.Sprintf(DOWNLOAD_ERROR.Message + ".\nUnderlying cause: %v", shaUrl, err), DOWNLOAD_ERROR.ExitCode)
	}
	return nil
}

func downloadFile(url string, dest string) error {

	file := path.Base(url)

	log.Printf("Downloading file %s from %s\n", file, url)

	var path bytes.Buffer
	path.WriteString(dest)
	path.WriteString("/")
	path.WriteString(file)

	start := time.Now()

	out, err := os.Create(path.String())
	if err != nil {
		return err
	}

	defer out.Close()

	client := &http.Client{}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("user-agent", "bzl")
	req.Header.Add("accept", "application/octet-stream")
	headResp, err := client.Do(req)
	defer headResp.Body.Close()

	dump, err := httputil.DumpRequest(req, true)
	fmt.Println(string(dump), err)
	
	dump, err = httputil.DumpResponse(headResp, true)
	fmt.Println(string(dump), err)
	
	if headResp.StatusCode != 200 {
		return cli.NewExitError(fmt.Sprintf("%d HTTP status while downloading <%s>", headResp.StatusCode, url), 1)
	}
	
	size := headResp.ContentLength

	done := make(chan int64)

	go printDownloadPercent(done, path.String(), int64(size))

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("user-agent", "bzl")
	req.Header.Add("accept", "application/octet-stream")

	dump, err = httputil.DumpRequest(req, true)
	fmt.Println(string(dump), err)
	
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	done <- n

	elapsed := time.Since(start)

	log.Printf("Download completed in %s", elapsed)
	return nil
}

func printDownloadPercent(done chan int64, path string, total int64) {

	var stop bool = false

	for {
		select {
		case <-done:
			stop = true
		default:

			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

			size := fi.Size()

			if size == 0 {
				size = 1
			}

			var percent float64 = float64(size) / float64(total) * 100

			fmt.Printf("%.0f", percent)
			fmt.Println("%")
		}

		if stop {
			break
		}

		time.Sleep(time.Second)
	}
}
