package oak_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOak(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oak Suite")
}
