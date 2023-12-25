/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2023 Red Hat, Inc.
 */

package main

import (
	"flag"
	"fmt"
	"os"

	goversion "github.com/hashicorp/go-version"
)

func main() {
	commands := map[string]func() int{
		"skiprange": skipRange,
		"previous":  previous,
		"is-head":   isHead,
	}
	if len(os.Args) <= 1 { // defensive
		usage("", commands)
		os.Exit(0)
	}
	cmd, ok := commands[os.Args[1]]
	if !ok {
		usage(os.Args[1], commands)
		os.Exit(99)
	}
	os.Args = os.Args[1:]
	os.Exit(cmd())
}

func usage(cmdName string, commands map[string]func() int) {
	if cmdName == "" {
		fmt.Fprintf(os.Stderr, "missing command\n")
	} else {
		fmt.Fprintf(os.Stderr, "unsupported command: %q\n", cmdName)
	}
	fmt.Fprintf(os.Stderr, "available commands:\n")
	for name := range commands {
		fmt.Fprintf(os.Stderr, "- %s\n", name)
	}
}

func isHead() int {
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "too many arguments (%d)\n", len(os.Args))
		return 1
	}
	v := os.Args[1]
	ver, err := goversion.NewVersion(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "not a version: %q (%v)\n", v, err)
		return 1
	}
	segs := ver.Segments()
	if len(segs) != 3 {
		fmt.Fprintf(os.Stderr, "only 3-components versions are supported, got %v\n", segs)
		return 2
	}
	patch := segs[2]
	if patch != 0 {
		fmt.Fprintf(os.Stderr, "not head: patch=%d (%v)\n", patch, segs)
		return 4
	}
	return 0
}

func previous() int {
	rawOutput := false
	flag.BoolVar(&rawOutput, "raw", false, "emit unquoted output")
	flag.Parse()
	for _, arg := range flag.Args() {
		ret, err := previousVersion(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error computing previous verion: %v\n", err)
			return 2
		}
		if rawOutput {
			fmt.Printf("%s\n", ret)
		} else {
			fmt.Printf("'%s'\n", ret)
		}
	}
	return 0
}

func skipRange() int {
	rawOutput := false
	skewBack := 0
	flag.IntVar(&skewBack, "skew-back", 1, "skew versions back")
	flag.BoolVar(&rawOutput, "raw", false, "emit unquoted output")
	flag.Parse()

	if skewBack < 0 {
		fmt.Fprintf(os.Stderr, "skew-back must be > 0 (got %d)\n", skewBack)
		return 1
	}

	for _, arg := range flag.Args() {
		ret, err := previousMinor(arg, skewBack)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error computing previous minor: %v\n", err)
			return 2
		}
		if rawOutput {
			fmt.Printf(">=%s <%s\n", ret, arg)
		} else {
			fmt.Printf("'>=%s <%s'\n", ret, arg)
		}
	}
	return 0
}

func previousVersion(v string) (string, error) {
	ver, err := goversion.NewVersion(v)
	if err != nil {
		return "", err
	}
	segs := ver.Segments()
	if len(segs) != 3 {
		return "", fmt.Errorf("only 3-components versions are supported, got %v", segs)
	}
	major, minor, patch := segs[0], segs[1], segs[2]
	if patch > 0 {
		return fmt.Sprintf("%d.%d.%d", major, minor, patch-1), nil
	}
	if minor > 0 {
		return fmt.Sprintf("%d.%d.%d", major, minor-1, patch), nil
	}
	if major > 0 {
		return fmt.Sprintf("%d.%d.%d", major-1, minor, patch), nil
	}
	return "", fmt.Errorf("cannot compute the previous version of %q", ver)
}

func previousMinor(v string, skew int) (string, error) {
	ver, err := goversion.NewVersion(v)
	if err != nil {
		return "", err
	}
	segs := ver.Segments()
	if len(segs) != 3 {
		return "", fmt.Errorf("only 3-components versions are supported, got %v", segs)
	}
	major, minor := segs[0], segs[1]
	if minor < skew {
		return "", fmt.Errorf("excessive skew=%d minor=%d", skew, minor)
	}
	return fmt.Sprintf("%d.%d.0", major, minor-skew), nil
}
