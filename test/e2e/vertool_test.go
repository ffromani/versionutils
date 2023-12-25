package e2e

import (
	"path/filepath"

	g "github.com/onsi/ginkgo/v2"
	o "github.com/onsi/gomega"
)

var _ = g.Describe("vertool skiprange", func() {
	g.DescribeTable("cooked output arguments",
		func(arg, expected string) {
			got, err := runCmdline([]string{
				filepath.Join(binariesPath, "vertool"),
				"skiprange",
				arg,
			})
			o.Expect(err).ToNot(o.HaveOccurred())
			o.Expect(got).To(o.Equal(expected))
		},
		g.Entry("mid-stream", "4.11.7", "'>=4.10.0 <4.11.7'"),
		g.Entry("begin-stream", "4.11.0", "'>=4.10.0 <4.11.0'"),
	)
})

var _ = g.Describe("vertool previous", func() {
	g.DescribeTable("cooked output arguments",
		func(arg, expected string) {
			got, err := runCmdline([]string{
				filepath.Join(binariesPath, "vertool"),
				"previous",
				arg,
			})
			o.Expect(err).ToNot(o.HaveOccurred())
			o.Expect(got).To(o.Equal(expected))
		},
		g.Entry("mid-stream", "4.11.7", "'4.11.6'"),
		g.Entry("begin-stream", "4.11.0", "'4.10.0'"),
	)
})

var _ = g.Describe("vertool is-head", func() {
	g.DescribeTable("cooked output arguments",
		func(arg string, expectSuccess bool) {
			_, err := runCmdline([]string{
				filepath.Join(binariesPath, "vertool"),
				"is-head",
				arg,
			})
			if expectSuccess {
				o.Expect(err).ToNot(o.HaveOccurred())
			} else {
				o.Expect(err).To(o.HaveOccurred())
			}
		},
		g.Entry("mid-stream", "4.11.7", false),
		g.Entry("begin-stream", "4.11.0", true),
		g.Entry("new-major", "1.0.0", true),
		g.Entry("new-minor", "2.0.7", false),
	)
})
