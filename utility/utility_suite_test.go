package utility_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtility(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utility Suite")
}
