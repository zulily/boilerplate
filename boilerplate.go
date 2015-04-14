/*
 * Copyright (C) 2015 zulily, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func copyFile(src, dst string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	// do the actual work
	_, err = io.Copy(w, r)
	return err
}

// deployScaffold creates the directory structure for a new Go project, copies
// any required non-template files into it, and returns the path to the root
// dir of the newly-created project
func deployScaffold(t Target) (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", errors.New("$GOPATH is not set")
	}

	if ex, err := exists(gopath); err != nil {
		return "", err
	} else if !ex {
		return "", fmt.Errorf("GOPATH does not exist at: %s", gopath)
	}

	fmt.Printf("GOPATH is: %s\n", gopath)
	// the "root" dir is at: $GOPATH/src/github.com/zulily/fizzbuzz
	root := path.Join(gopath, "src", t.Repository, t.Namespace, t.Project)

	if ex, err := exists(root); err != nil {
		return "", err
	} else if ex {
		fmt.Printf("%s already exists. Overwrite existing files? [y/n]: ", root)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return "", err
		} else if !strings.EqualFold(strings.TrimSpace(text), "y") {
			return "", fmt.Errorf("%s already exists", root)
		}
	}

	fmt.Printf("Boilerplating the project at: %s\n", root)

	if err := os.MkdirAll(path.Join(root, "build"), 0755); err != nil {
		return "", err
	}

	if err := copyFile("build/Dockerfile", path.Join(root, "build", "Dockerfile")); err != nil {
		return "", err
	}

	return root, nil
}

// deployTemplate parses and executues a template to a new file under the
// specified `root` dir. Template files are assumed to end with ".template".
// A templated file named `foo.template` will be placed at `root/foo`.
func deployTemplate(root, tmpl string, target Target) error {
	fm := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}
	fname := tmpl[:len(tmpl)-len(".template")]
	fmt.Printf("Creating new: %s\n", fname)
	t := template.Must(template.New(tmpl).Funcs(fm).ParseFiles(tmpl))

	if f, err := os.Create(path.Join(root, fname)); err != nil {
		return err
	} else {
		defer f.Close()
		return t.Execute(f, target)
	}
}

// Target represents a Go build target; typically a binary
type Target struct {

	// Respository is the name of the source control repository (e.g. github.com)
	Repository string

	// Namespace is the name of the organization/group in the repository (e.g. zulily)
	Namespace string

	// Project is the name of the binary or package (e.g. fizzbuzz)
	Project string
}

var opts struct {
	verbose bool
	Target
}

func main() {

	flag.StringVar(&opts.Repository, "repository", "", "the name of the git repository (e.g. github.com)")
	flag.StringVar(&opts.Namespace, "namespace", "", "the name of the organization/group in the repository (e.g. zulily)")
	flag.StringVar(&opts.Project, "project", "", "the name of the project (e.g. fizzbuzz)")
	flag.BoolVar(&opts.verbose, "verbose", false, "toggles verbose output")
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)

	if opts.Repository == "" {
		fmt.Printf("Enter the name of git repository (e.g. github.com): ")
		if scanner.Scan() {
			opts.Repository = scanner.Text()
		}
	}

	if opts.Namespace == "" {
		fmt.Printf("Enter the namespace in the repository (e.g. zulily): ")
		if scanner.Scan() {
			opts.Namespace = scanner.Text()
		}
	}

	if opts.Project == "" {
		fmt.Printf("Enter the name of the project (e.g. fizzbuzz): ")
		if scanner.Scan() {
			opts.Project = scanner.Text()
		}
	}

	root, err := deployScaffold(opts.Target)
	if err != nil {
		panic(err)
	}

	if files, err := filepath.Glob("*.template"); err == nil {
		for _, templ := range files {
			if err = deployTemplate(root, templ, opts.Target); err != nil {
				panic(err)
			}
		}
	} else {
		panic(err)
	}

	out := ioutil.Discard
	if opts.verbose {
		out = os.Stdout
	}

	fmt.Println("Initializing git repo")
	c := exec.Command("git", "init")
	c.Dir = root
	c.Stdout, c.Stderr = out, out
	if err = c.Run(); err != nil {
		panic(err)
	}

	fmt.Println("Initializing godeps")
	c = exec.Command("make", "godep")
	c.Dir = root
	c.Stdout, c.Stderr = out, out
	if err = c.Run(); err != nil {
		panic(err)
	}

	fmt.Println("Done")
}
