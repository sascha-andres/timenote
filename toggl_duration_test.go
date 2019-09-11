package timenote_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"livingit.de/code/timenote"
	"testing"
)

func dotest(n int64, exp string) {
	td, _ := timenote.NewTogglDuration(n)
	Expect(td.String()).To(Equal(exp))
}

func dotesterror(n int64, exp string) {
	_, err := timenote.NewTogglDuration(n)
	errText := "no error"
	if err != nil {
		errText = err.Error()
	}
	Expect(errText).To(Equal("negative values not allowed"))
}

var _ = Describe("Human readable duration in seconds", func() {
	It("should handle finished tasks", func() {
		dotest(10, "10s")
		dotest(60, "1m0s")
		dotest(3600, "1h0m0s")
		dotest(3666, "1h1m6s")
		dotest(3600*24, "1d0h0m0s")
	})
	It("should handle open tasks", func() {
		dotesterror(-10, "negative values not allowed")
		dotesterror(-60, "negative values not allowed")
		dotesterror(-3600, "negative values not allowed")
		dotesterror(-3666, "negative values not allowed")
		dotesterror(-3600*24, "negative values not allowed")
	})
})

func TestFormat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Human readable duration in seconds")
}
