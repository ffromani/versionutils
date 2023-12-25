package e2e

import (
	"path/filepath"

	g "github.com/onsi/ginkgo/v2"
	o "github.com/onsi/gomega"
)

var _ = g.Describe("skiprange", func() {
	g.Context("without arguments", func() {
		g.It("emits correct cooked output", func() {
			got, err := runCmdline([]string{
				filepath.Join(binariesPath, "skiprange"),
				"4.11.7",
			})
			o.Expect(err).ToNot(o.HaveOccurred())
			expected := "'>=4.10.0 <4.11.7'"
			o.Expect(got).To(o.Equal(expected))
		})

		g.It("emits correct raw output", func() {
			got, err := runCmdline([]string{
				filepath.Join(binariesPath, "skiprange"),
				"--raw",
				"4.12.2",
			})
			o.Expect(err).ToNot(o.HaveOccurred())
			expected := ">=4.11.0 <4.12.2"
			o.Expect(got).To(o.Equal(expected))
		})
	})
})
