//  (C) Copyright 2014 yum-nginx-api Contributors.
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//  http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	//"time"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/access"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/fault"
	"github.com/go-ozzo/ozzo-routing/v2/slash"
	//"github.com/go-ozzo/ozzo-routing/v2/file"
	"github.com/h2non/filetype"
)

func open_lock_file() (*os.File, error) {
	return os.OpenFile(
		filepath.Join(uploadDir, ".createrepo.lck"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
}

func close_lock_file(f *os.File) {
	f.Close()
}

func acquire_lock_file(f *os.File, mode int) {
	syscall.Flock(int(f.Fd()), mode)
}

func release_lock_file(f *os.File) {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}

type repo_chan_struct struct {
	reponame string
}

var repo_chan chan repo_chan_struct
var log *Logger

func update_repo(repodir string) {
	f, err := open_lock_file()
	if err != nil {
		log.Print("Unable to open createrepo lock file ", err.Error())
	} else {
		var crExec *exec.Cmd
		acquire_lock_file(f, syscall.LOCK_EX)
		if strings.HasSuffix(crBin, "_c") {
			crExec = exec.Command(crBin, "--update", "--compress-type", "bzip2", "--general-compress-type", "bzip2", "--workers", createRepo, repodir)
		} else {
			crExec = exec.Command(crBin, "--update", "--workers", createRepo, repodir)
		}
		var out []byte
		out, err = crExec.CombinedOutput()
		release_lock_file(f)
		close_lock_file(f)

		if err != nil {
			log.Print(fmt.Sprintf("Unable to execute createrepo - %v [%s]", err, out))
		}
	}
}

// crRoutine is a simple buffer to not overload the system
// by running too many createrepo system commands, uncompress,
// and sqlite connections at the same time
func crRoutine() {
	var rcd repo_chan_struct

	for {
		select {
		case rcd = <-repo_chan:
			log.Print(fmt.Sprintf("UPDATING REPO [%s]", rcd.reponame))
			update_repo(filepath.Join(uploadDir, rcd.reponame))
		}
	}
}

// api_upload function for handler /upload
func api_upload(c *routing.Context) error {
	var reponame string

	reponame = c.Param("reponame")

	err := c.Request.ParseMultipartForm(maxLength)
	if err != nil {
		c.Response.WriteHeader(http.StatusInternalServerError)
		_ = c.Write("Upload Failed " + err.Error())
		return err
	}
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.Response.WriteHeader(http.StatusInternalServerError)
		_ = c.Write("Upload Failed " + err.Error())
		return err
	}
	defer file.Close()
	filePath := filepath.Join(uploadDir, reponame, handler.Filename)
	dirPath := filepath.Join(uploadDir, reponame)
	os.MkdirAll(dirPath, 0755) // ignore errors
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		c.Response.WriteHeader(http.StatusInternalServerError)
		_ = c.Write("Upload Failed " + err.Error())
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		c.Response.WriteHeader(http.StatusInternalServerError)
		_ = c.Write("Upload Failed " + err.Error())
		return err
	}
	buf, _ := ioutil.ReadFile(filePath)
	if kind, err := filetype.Match(buf); err != nil || kind.MIME.Value != "application/x-rpm" {
		err := os.Remove(filePath)
		if err != nil {
			log.Print("Unable to delete " + filePath)
		}
		c.Response.WriteHeader(http.StatusUnsupportedMediaType)
		return c.Write(handler.Filename + " not RPM")
	}
	// If not in development mode increment create-repo counter
	// for command to be ran by go routine crRoutine
	if !devMode {
		repo_chan <- repo_chan_struct{reponame: reponame}
	}
	c.Response.WriteHeader(http.StatusAccepted)
	return c.Write("Uploaded")
}

// api_health function for handler /health
func api_health(c *routing.Context) error {
	c.Response.Header().Add("Version", commitHash)
	return c.Write("OK")
}

// api_repo function for handler /repo
//func api_repo(c *routing.Context) error {
//	return c.Write(rJSON)
//}

func repo_file_access(c *routing.Context) error {
	var lf *os.File
	var file *os.File
	var fstat os.FileInfo
	var path string
	var err error

	path = filepath.Join(uploadDir, c.Param("reponame"), c.Param(""))

	lf, err = open_lock_file()
	if err != nil {
		return routing.NewHTTPError(http.StatusForbidden, err.Error())
	}

	acquire_lock_file(lf, syscall.LOCK_SH)
	file, err = os.Open(path)
	if err != nil {
		release_lock_file(lf)
		close_lock_file(lf)
		return routing.NewHTTPError(http.StatusNotFound, err.Error())
	}
	fstat, err = file.Stat()
	if err != nil {
		file.Close()
		release_lock_file(lf)
		close_lock_file(lf)
		return routing.NewHTTPError(http.StatusNotFound, err.Error())
	} else if fstat.IsDir() {
		file.Close()
		release_lock_file(lf)
		close_lock_file(lf)
		return routing.NewHTTPError(http.StatusNotFound)
	}
	c.Response.Header().Del("Content-Type")
	http.ServeContent(c.Response, c.Request, path, fstat.ModTime(), file)
	file.Close()
	release_lock_file(lf)
	close_lock_file(lf)

	return nil
}

func main() {
	log = NewLogger()

	if err := configValidate(); err != nil {
		log.Fatalln(err.Error())
	}

	repo_chan = make(chan repo_chan_struct, 10)

	go crRoutine()
	rtr := routing.New()

	api := rtr.Group("/api",
		fault.Recovery(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		content.TypeNegotiator(content.JSON))

	repo := rtr.Group("/repo",
		fault.Recovery(log.Printf),
		slash.Remover(http.StatusMovedPermanently),
		content.TypeNegotiator(content.JSON))

	// Disable logging on health endpoints
	api.Get("/health", api_health)
	api.Use(access.Logger(log.Printf))
	//api.Post("/upload/<reponame>/*", api_upload
	api.Post("/upload/<reponame>", api_upload) // support a single level repository /repo/aaa, /repo/bbb, etc
	//api.Get("/repo", api_repo)

	repo.Use(access.Logger(log.Printf))
	//repo.Get("/<name>", file.Server(file.PathMap{"/repo/<name>": uploadDir}))
	//repo.Get("/*", file.Server(file.PathMap{"/repo": uploadDir}))
	repo.Get("/<reponame>/*", repo_file_access)

	http.Handle("/", rtr)
	log.Printf("built-on %s, version %s started on %s", builtOn, commitHash, port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Panicln(err)
	}

	close(repo_chan)
}
