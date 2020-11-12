package gheap_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoHeap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoHeap Suite")
}
