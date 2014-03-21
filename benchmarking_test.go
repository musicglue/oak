package oak_test

import (
	. "github.com/musicglue/oak"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Performance", func() {
	type Options struct {
		Path  []string
		Value interface{}
	}

	var (
		branch *Branch
		req    = Options{Path: []string{"a", "b", "c", "d", "e"}, Value: "Winning"}
	)

	BeforeEach(func() {
		branch = NewBranch()
		branch.Set(req.Path, req.Value)
	})

	Measure("It should be quick to update records", func(b Benchmarker) {
		runtime := b.Time("writing", func() {
			for i := 0; i < 10000; i++ {
				branch.Set(req.Path, i)
			}
			_, ok := branch.Get(req.Path)
			Expect(ok).To(BeTrue())
		})

		Ω(runtime.Seconds()).Should(BeNumerically("<", 0.01), "branch.Set() shouldn't take too long.")
	}, 100)

	Measure("It should be really quick to read records", func(b Benchmarker) {
		runtime := b.Time("reading", func() {
			var ok bool
			for i := 0; i < 10000; i++ {
				_, ok = branch.Get(req.Path)
				if !ok {
					break
				}
			}
			Expect(ok).To(BeTrue())
		})

		Ω(runtime.Seconds()).Should(BeNumerically("<", 0.01), "branch.Get() really shouldn't take too long.")
	}, 100)
})
