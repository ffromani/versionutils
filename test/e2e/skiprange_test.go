package e2e

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	g "github.com/onsi/ginkgo/v2"
	o "github.com/onsi/gomega"
)

var _ = g.Describe("skiprange", func() {
	g.Context("without arguments", func() {
		g.It("emits correct cooked output", func() {
			cmdline := []string{
				filepath.Join(binariesPath, "skiprange"),
				"4.11.7",
			}
			fmt.Fprintf(g.GinkgoWriter, "running: %v\n", cmdline)

			cmd := exec.Command(cmdline[0], cmdline[1:]...)
			cmd.Stderr = g.GinkgoWriter

			out, err := cmd.Output()
			o.Expect(err).ToNot(o.HaveOccurred())

			got := strings.TrimSpace(string(out))
			expected := "'>=4.10.0 <4.11.7'"

			o.Expect(got).To(o.Equal(expected))
		})

		g.It("emits correct raw output", func() {
			cmdline := []string{
				filepath.Join(binariesPath, "skiprange"),
				"--raw",
				"4.12.2",
			}
			fmt.Fprintf(g.GinkgoWriter, "running: %v\n", cmdline)

			cmd := exec.Command(cmdline[0], cmdline[1:]...)
			cmd.Stderr = g.GinkgoWriter

			out, err := cmd.Output()
			o.Expect(err).ToNot(o.HaveOccurred())

			got := strings.TrimSpace(string(out))
			expected := ">=4.11.0 <4.12.2"

			o.Expect(got).To(o.Equal(expected))
		})
	})
})
