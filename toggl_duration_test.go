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

func dotestnoseconds(n int64, exp string) {
	td, _ := timenote.NewTogglDuration(n)
	td.OmitSeconds()
	Expect(td.String()).To(Equal(exp))
}

func dotesterror(n int64, exp string) {
	_, err := timenote.NewTogglDuration(n)
	errText := "no error"
	if err != nil {
		errText = err.Error()
	}
	Expect(errText).To(Equal(exp))
}

var _ = Describe("Human readable duration in seconds", func() {
	It("should handle finished tasks", func() {
		dotest(10, "10s")
		dotest(60, "1m 0s")
		dotest(3600, "1h 0m 0s")
		dotest(3666, "1h 1m 6s")
		dotest(3600*24, "1d 0h 0m 0s")
	})
	It("should omit seconds when asked", func() {
		dotestnoseconds(10, "")
		dotestnoseconds(60, "1m")
		dotestnoseconds(3600, "1h 0m")
		dotestnoseconds(3666, "1h 1m")
		dotestnoseconds(3600*24, "1d 0h 0m")
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
