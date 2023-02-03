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
	skewBack := 0
	flag.IntVar(&skewBack, "skew-back", 1, "skew versions back")
	flag.Parse()

	if skewBack < 0 {
		fmt.Fprintf(os.Stderr, "skew-back must be > 0 (got %d)\n", skewBack)
		os.Exit(1)
	}

	for _, arg := range flag.Args() {
		ret, err := previousMinor(arg, skewBack)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error computing previous minor: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("'>=%s <%s'\n", ret, arg)
	}
}

func previousMinor(v string, skew int) (string, error) {
	ver, err := goversion.NewVersion(v)
	if err != nil {
		return "", nil
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
