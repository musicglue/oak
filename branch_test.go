package oak_test

import (
	. "github.com/musicglue/oak"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Branch", func() {
	type Options struct {
		Path  []string
		Value interface{}
	}

	var (
		branch *Branch
	)

	BeforeEach(func() {
		branch = NewBranch()
	})

	Describe("Constuctor", func() {
		It("Should build a proper Branch", func() {
			Expect(branch.Value).To(BeNil())
			Expect(len(branch.Branches)).To(Equal(0))
		})
	})

	Describe("Get", func() {
		var (
			req1 = Options{Path: []string{}, Value: "Home"}
			req2 = Options{Path: []string{"employees"}, Value: "Employees Home"}
			req3 = Options{Path: []string{"employees", "1"}, Value: "Bob Smith"}
			reqs = []Options{req1, req2, req3}
		)

		BeforeEach(func() {
			for _, req := range reqs {
				branch.Set(req.Path, req.Value)
			}
		})

		It("Should return the base node value for an empty query", func() {
			result, ok := branch.Get(req1.Path)
			Expect(result).To(Equal(req1.Value))
			Expect(ok).To(BeTrue())
		})

		It("Should return the value for a single nesting level", func() {
			result, ok := branch.Get(req2.Path)
			Expect(result).To(Equal(req2.Value))
			Expect(ok).To(BeTrue())
		})

		It("Should return the value for the deepest path", func() {
			result, ok := branch.Get(req3.Path)
			Expect(result).To(Equal(req3.Value))
			Expect(ok).To(BeTrue())
		})

		It("Should return a nil/false response for non-matching paths", func() {
			result, ok := branch.Get([]string{"nope"})
			Expect(result).To(BeNil())
			Expect(ok).To(BeFalse())
		})

		It("Should return a nil/false response for non-matching deep nestings", func() {
			result, ok := branch.Get([]string{"not", "here", "either"})
			Expect(result).To(BeNil())
			Expect(ok).To(BeFalse())
		})
	})

	Describe("Longest Match", func() {
		var (
			path1 = []string{"categories", "news", "latest", "headlines"}
			req1  = Options{Path: path1, Value: "Headlines"}
			path2 = []string{"categories"}
			req2  = Options{Path: path2, Value: "Categories"}
			reqs  = []Options{req1, req2}
		)

		BeforeEach(func() {
			for _, req := range reqs {
				branch.Set(req.Path, req.Value)
			}
		})

		It("Should return the right score for a direct match", func() {
			result, ok := branch.Match(path1)
			Expect(result).To(Equal("Headlines"))
			Expect(ok).To(BeTrue())
		})

		It("Should return the parent value on a missing key", func() {
			result, ok := branch.Match([]string{"categories", "news", "latest"})
			Expect(result).To(Equal("Categories"))
			Expect(ok).To(BeTrue())
		})

		It("Should return the first only matching key if one match found", func() {
			result, ok := branch.Match([]string{"categories"})
			Expect(result).To(Equal("Categories"))
			Expect(ok).To(BeTrue())
		})
	})

	Describe("Set", func() {
		var (
			path = []string{"home"}
			req1 = Options{Path: path, Value: "Home"}
			req2 = Options{Path: path, Value: "New Home"}
		)

		BeforeEach(func() {
			branch.Set(req1.Path, req1.Value)
		})

		Context("Simple trees", func() {
			It("Should return the initial value", func() {
				result, ok := branch.Get(path)
				Expect(result).To(Equal(req1.Value))
				Expect(ok).To(BeTrue())
			})

			It("Should return the new value if updated", func() {
				branch.Set(req2.Path, req2.Value)
				result, ok := branch.Get(path)
				Expect(result).To(Equal(req2.Value))
				Expect(result).NotTo(Equal(req1.Value))
				Expect(ok).To(BeTrue())
			})
		})

		Context("Deep trees, with intermediaries", func() {
			var (
				path1 = []string{"categories", "news"}
				path2 = []string{"categories", "news", "today"}
				req1  = Options{Path: path1, Value: "Exciting"}
				req2  = Options{Path: path2, Value: "Depressing"}
				reqs  = []Options{req1, req2}
			)

			BeforeEach(func() {
				for _, req := range reqs {
					branch.Set(req.Path, req.Value)
				}
			})

			Context("Overwriting the deeper node leaves the shallower one intact", func() {
				var (
					req3 = Options{Path: path2, Value: "Overwritten"}
				)

				BeforeEach(func() {
					branch.Set(req3.Path, req3.Value)
				})

				It("Returns the new value", func() {
					result, ok := branch.Get(path2)
					Expect(result).To(Equal(req3.Value))
					Expect(ok).To(BeTrue())
				})

				It("Doesn't change the shallower value", func() {
					result, ok := branch.Get(path1)
					Expect(result).To(Equal(req1.Value))
					Expect(ok).To(BeTrue())
				})
			})

			Context("Overwriting a shallower node leaves the deeper ones intact", func() {
				var (
					req3 = Options{Path: path1, Value: "Overwritten"}
				)

				BeforeEach(func() {
					branch.Set(req3.Path, req3.Value)
				})

				It("Returns the new value", func() {
					result, ok := branch.Get(path1)
					Expect(result).To(Equal(req3.Value))
					Expect(ok).To(BeTrue())
				})

				It("Returns the nested value unaltered", func() {
					result, ok := branch.Get(path2)
					Expect(result).To(Equal(req2.Value))
					Expect(ok).To(BeTrue())
				})
			})
		})

	})

	Describe("Remove", func() {
		var (
			req1 = Options{Path: []string{}, Value: "Home"}
			req2 = Options{Path: []string{"employees"}, Value: "Employees Home"}
			req3 = Options{Path: []string{"employees", "1"}, Value: "Bob Smith"}
			reqs = []Options{req1, req2, req3}
		)

		BeforeEach(func() {
			for _, req := range reqs {
				branch.Set(req.Path, req.Value)
			}
		})

		It("Removes the named node", func() {
			outcome := branch.Remove(req2.Path)
			Expect(outcome).To(Equal(true))
			result, ok := branch.Get(req2.Path)
			Expect(result).To(BeNil())
			Expect(ok).To(BeFalse())
		})

		It("Removes nested nodes if there are any", func() {
			outcome := branch.Remove(req2.Path)
			Expect(outcome).To(Equal(true))
			result, ok := branch.Get(req3.Path)
			Expect(result).To(BeNil())
			Expect(ok).To(BeFalse())
		})

		It("Returns false if there are no matching nodes", func() {
			outcome := branch.Remove([]string{"nothing here"})
			Expect(outcome).To(Equal(false))
		})

		It("Returns false if you try and remove the root node", func() {
			outcome := branch.Remove([]string{})
			Expect(outcome).To(Equal(false))
		})
	})

	Describe("Replace", func() {
		var (
			path1   = []string{"categories", "news"}
			path2   = []string{"categories", "news", "today"}
			req1    = Options{Path: path1, Value: "Exciting"}
			req2    = Options{Path: path2, Value: "Depressing"}
			reqs    = []Options{req1, req2}
			newnode *Branch
			newval  = "New Value"
		)

		BeforeEach(func() {
			for _, req := range reqs {
				branch.Set(req.Path, req.Value)
			}
			newnode = NewBranch()
			newnode.Set([]string{}, newval)
		})

		Context("Replacing a shallow node replaces the sub-tree too", func() {
			BeforeEach(func() {
				branch.Replace(path1, newnode)
			})

			It("Sets the value at path1", func() {
				result, ok := branch.Get(path1)
				Expect(result).To(Equal(newval))
				Expect(ok).To(BeTrue())
			})

			It("Wipes the value at path2", func() {
				result, ok := branch.Get(path2)
				Expect(result).To(BeNil())
				Expect(ok).To(BeFalse())
			})
		})
	})
})
