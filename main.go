/*
** MIT License
**
** Copyright (c) 2021 Ryan SVIHLA
**
** Permission is hereby granted, free of charge, to any person obtaining a copy
** of this software and associated documentation files (the "Software"), to deal
** in the Software without restriction, including without limitation the rights
** to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
** copies of the Software, and to permit persons to whom the Software is
** furnished to do so, subject to the following conditions:
**
** The above copyright notice and this permission notice shall be included in all
** copies or substantial portions of the Software.
**
** THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
** IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
** FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
** AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
** LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
** OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
** SOFTWARE.
 */
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func usage() string {
	return "usage: go-init <new project dir>"
}
func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println(usage())
		os.Exit(1)
	}
	dir := args[1]
	org := args[2]
	if err := createProject(dir, org); err != nil {
		fmt.Printf("unable to create project for dir '%v' with error '%v'", dir, err)
		os.Exit(1)
	}
	fmt.Printf("project %v created\ngithub repo is expected to be http://github.com/%v/%v\n", dir, org, dir)
}

type Gen struct {
	Name string
	Perm os.FileMode
	Get  func() string
}

func createProject(dir string, org string) error {
	year := time.Now().Year()
	author, err := getGitName()
	if err != nil {
		return err
	}
	if err := os.Mkdir(dir, 0755); err != nil {
		return fmt.Errorf("unable to create dir '%v' with error '%v'", dir, err)
	}
	scriptsDir := filepath.Join(dir, "scripts")
	if err := os.Mkdir(scriptsDir, 0755); err != nil {
		return fmt.Errorf("unable to create dir '%v' with error '%v'", scriptsDir, err)
	}

	binDir := filepath.Join(dir, "bin")
	if err := os.Mkdir(binDir, 0755); err != nil {
		return fmt.Errorf("unable to create dir '%v' with error '%v'", binDir, err)
	}
	gens := []Gen{
		{
			Name: filepath.Join(dir, "README.md"),
			Perm: 0644,
			Get: func() string {
				return createReadme(year, dir, author, org)
			},
		},
		{
			Name: filepath.Join(dir, "LICENSE"),
			Perm: 0644,
			Get: func() string {
				return createLicense(year, author)
			},
		},
		{
			Name: filepath.Join(dir, "main.go"),
			Perm: 0644,
			Get: func() string {
				return createMain(year, dir, author)
			},
		},
		{
			Name: filepath.Join(dir, "main_test.go"),
			Perm: 0644,
			Get: func() string {
				return createMainTest(year, dir, author)
			},
		},
		{
			Name: filepath.Join(dir, ".gitignore"),
			Perm: 0644,
			Get:  createGitIgnore,
		},
		{
			Name: filepath.Join(dir, "go.mod"),
			Perm: 0644,
			Get: func() string {
				return createGoMod(org, dir)
			},
		},
		{
			Name: filepath.Join(scriptsDir, "all"),
			Perm: 0755,
			Get:  createAllScript,
		},
		{
			Name: filepath.Join(scriptsDir, "bootstrap"),
			Perm: 0755,
			Get:  createBootsrapScript,
		},
		{
			Name: filepath.Join(scriptsDir, "build"),
			Perm: 0755,
			Get: func() string {
				return createBuildScript(dir)
			},
		},
		{
			Name: filepath.Join(scriptsDir, "cibuild"),
			Perm: 0755,
			Get:  createCiScript,
		},
		{
			Name: filepath.Join(scriptsDir, "clean"),
			Perm: 0755,
			Get:  createCleanScript,
		},
		{
			Name: filepath.Join(scriptsDir, "cover-html"),
			Perm: 0755,
			Get:  createCoverHtml,
		},
		{
			Name: filepath.Join(scriptsDir, "install.sh"),
			Perm: 0755,
			Get: func() string {
				return createInstall(org, dir)
			},
		},
		{
			Name: filepath.Join(scriptsDir, "lint"),
			Perm: 0755,
			Get:  createLint,
		},
		{
			Name: filepath.Join(scriptsDir, "package"),
			Perm: 0755,
			Get: func() string {
				return createPackage(dir)
			},
		},
		{
			Name: filepath.Join(scriptsDir, "setup"),
			Perm: 0755,
			Get:  createSetupScript,
		},
		{
			Name: filepath.Join(scriptsDir, "test"),
			Perm: 0755,
			Get:  createTests,
		},
		{
			Name: filepath.Join(scriptsDir, "update"),
			Perm: 0755,
			Get:  createUpdate,
		},
	}
	for _, g := range gens {
		if err := os.WriteFile(g.Name, []byte(g.Get()), g.Perm); err != nil {
			return fmt.Errorf("failure writing '%v' with error '%v'", g.Name, err)
		}
	}

	return nil
}

func getGitName() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("unable to get git config author name with error %v", err)
	}
	return string(stdout), nil
}

func createLicense(year int, author string) string {
	title := fmt.Sprintf("MIT License\n\nCopyright (c) %v %v", year, author)
	body := `
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`
	return strings.Join([]string{title, body}, "")
}

func createReadme(year int, projectName string, author string, org string) string {
	title := fmt.Sprintf("# %v\n\nPLACEHOLDER\n", projectName)
	use := fmt.Sprintf("## how to use\n\n````sh\n%v\n```\n\n", projectName)
	install := fmt.Sprintf("## how to install\n\n```sh\ncurl -s https://raw.githubusercontent.com/%v/%v/scripts/install.sh\n```\n", org, projectName)
	build := fmt.Sprintf("## how to build and run\n\n````sh\n./scripts/build\n./bin/%v\n```\n\n", projectName)
	test := "## how to build and run\n\n````sh\n./scripts/test\n```\n\n##license\n\n"
	return strings.Join([]string{title, use, install, build, test, createLicense(year, author)}, "")
}

func createGitIgnore() string {
	return `# ---> Go
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

bin/
# Dependency directories (remove the comment below to include it)
# vendor/
.DS_Store

dist/
`
}

func createMain(year int, proj string, author string) string {
	title := `/*
** MIT License
**
`
	copywrite := fmt.Sprintf("** Copyright (c) %v %v", year, author)
	body := `** 
** Permission is hereby granted, free of charge, to any person obtaining a copy
** of this software and associated documentation files (the "Software"), to deal
** in the Software without restriction, including without limitation the rights
** to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
** copies of the Software, and to permit persons to whom the Software is
** furnished to do so, subject to the following conditions:
** 
** The above copyright notice and this permission notice shall be included in all
** copies or substantial portions of the Software.
** 
** THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
** IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
** FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
** AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
** LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
** OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
** SOFTWARE.
*/
package main

import (
    "os"
    "fmt"
)

func usage() string {

`
	usage := fmt.Sprintf("    return \"usage: %v\"", proj)
	footer := `
}

func main(){
    args := os.Args
    if len(args) < 2 {
        fmt.Println(usage())
        os.Exit(1)
    }
}
`
	return strings.Join([]string{title, copywrite, body, usage, footer}, "")
}

func createMainTest(year int, proj string, author string) string {
	title := `/*
** MIT License
**
`
	copywrite := fmt.Sprintf("** Copyright (c) %v %v", year, author)
	body := `** 
** Permission is hereby granted, free of charge, to any person obtaining a copy
** of this software and associated documentation files (the "Software"), to deal
** in the Software without restriction, including without limitation the rights
** to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
** copies of the Software, and to permit persons to whom the Software is
** furnished to do so, subject to the following conditions:
** 
** The above copyright notice and this permission notice shall be included in all
** copies or substantial portions of the Software.
** 
** THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
** IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
** FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
** AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
** LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
** OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
** SOFTWARE.
*/
package main

import (
    "testing"
)

func TestFunc(t * testing.T){
    if 1 != 2 {
        t.Errorf("update test to be useful")
    }
}
`
	return strings.Join([]string{title, copywrite, body}, "")
}

func createGoMod(org string, projectName string) string {
	return fmt.Sprintf("module github.com/%v/%v\n\ngo 1.17", org, projectName)
}

func createAllScript() string {
	return `#!/bin/bash
    # scripts/all: Runs several things together: lint, test, clean and build
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

$DIR/lint
$DIR/test
$DIR/clean
$DIR/build
`
}

func createBootsrapScript() string {
	return `
#!/bin/bash
# scripts/bootstrap: Resolve all dependencies that the application requires to
#                   run.



GOV="1.17.1"
if ! command -v go &> /dev/null
then
    echo "os $(uname -s) arch $(uname -m)"
    if [ "$(uname -s)" = "Darwin" ]; then
        echo "install via homebrew"
        brew update
        brew install go
    fi

    if [ "$(uname -s)" = "Linux" ] && [ "$(uname -m)" = "armv7l" ]; then
        echo "arm found installing go"
        curl -L -O https://golang.org/dl/go$GOV.linux-armv6l.tar.gz
        sudo tar -C /usr/local -xzf go$GOV.linux-armv6l.tar.gz
        echo "add ‘export PATH=\$PATH:/usr/local/go/bin’ to your .bashrc"
        rm go$GOV.linux-armv6l.tar.gz
    fi

    if [ "$(uname -s)" = "Linux" ] && [ "$(uname -m)" = "amd64" ]; then
        echo "amd64 found installing go"
        curl -L -O https://golang.org/dl/go$GOV.linux-amd64.tar.gz
        sudo tar -C /usr/local -xzf go$GOV.linux-amd64.tar.gz
        echo "add 'export PATH=\$PATH:/usr/local/go/bin' to your .bashrc"
        rm go$GOV.linux-arm64.tar.gz
    fi

else
    echo "go installed skipping"
fi

if ! command -v golangci-lint &> /dev/null
then
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.36.0
else
    echo "golangci-lint installed skipping"
fi
`
}

func createBuildScript(projectName string) string {
	return strings.Join([]string{`#!/bin/bash
# scripts/build: Compiles binary and outputs it to the bin folder

rm -fr ./bin
mkdir ./bin
`,
		fmt.Sprintf("go build -o bin/%v .\n", projectName)}, "")
}

func createCiScript() string {
	return `#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# scripts/cibuild: Setup environment for CI to run tests. This is primarily
#                 designed to run on the continuous integration server.


$DIR/setup && \
$DIR/lint && \
$DIR/test && \
$DIR/build
`
}

func createCleanScript() string {
	return `#!/bin/bash
# scripts/clean: Remove build binary files

rm -fr ./bin
mkdir ./bin
`
}

func createCoverHtml() string {
	return `#!/bin/bash
# scripts/cover-html: See the coverage report in a webpage

t="/tmp/go-cover.$$.tmp"
go test -race -covermode=atomic -coverprofile=$t ./... && go tool cover -html=$t && unlink $t
`

}

func createInstall(org string, projectName string) string {
	return strings.Join([]string{
		fmt.Sprintf("#!/usr/bin/env bash\norig_dir=$(pwd)\ncd /tmp\ncurl -O https://github.com/%v/%v/archive/refs/heads/main.zip\n", org, projectName),
		`# scripts/install.sh: install script for others to use, install.sh is a convention and why the name is different
unzip main.zip
rm main.zip
cd main
./scripts/build
echo "copying binary to /usr/local/bin/go-init need sudo permissions to write"
sudo cp ./bin/go-init /usr/local/bin/
cd ..
rm -fr main
cd $orig_dir
`}, "")
}

func createLint() string {
	return `#!/bin/bash
# scripts/lint: verify no obvious bugs or layout problems are found

gofmt -s -w . && \
golangci-lint run
`
}

func createPackage(projectName string) string {
	return strings.Join([]string{`#!/bin/bash
# scripts/package: build and tgz all supported platforms and architectures
`,
		fmt.Sprintf("BIN=%v", projectName),
		`DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
$DIR/clean
VERSION=$(git describe --abbrev=0 --tags)
ORIG=$(git branch --show-current)
echo "packaging $VERSION$"
git checkout $VERSION
GOOS=darwin GOARCH=amd64 go build -o bin/$BIN .
tar czvf ./bin/$BIN-$VERSION-darwin-amd64.tgz ./bin/$BIN
GOOS=darwin GOARCH=arm64 go build -o bin/$BIN .
tar czvf ./bin/$BIN-$VERSION-darwin-arm64.tgz ./bin/$BIN
GOOS=linux GOARCH=amd64 go build -o bin/$BIN .
tar czvf ./bin/$BIN-$VERSION-linux-amd64.tgz ./bin/$BIN
GOOS=linux GOARCH=arm64 go build -o bin/$BIN .
tar czvf ./bin/$BIN-$VERSION-linux-arm64.tgz ./bin/$BIN
git checkout $ORIG
`}, "")
}

func createSetupScript() string {
	return `#!/bin/bash
# scripts/setup: Set up application for the first time after cloning, or set it
#               back to the initial first unused state.

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
$DIR/bootstrap && \
go mod verify
`
}

func createTests() string {
	return `#!/bin/bash
# scripts/test: Run test suite for application. 

t="/tmp/go-cover.$$.tmp"
go test -race -covermode=atomic  -coverprofile=$t ./... && go tool cover -func=$t 
last=$?
unlink $t || true
if [ "$last" = "0" ]; then
    echo "successfully ran"
else
    (exit 1)
fi
`
}

func createUpdate() string {
	return `#!/bin/bash
# scripts/update: Update application to run for its current checkout.

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
$DIR/bootstrap
go mod tidy
`
}
