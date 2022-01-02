package utility_test

import (
	"fmt"

	"github.com/imaskm/getir/utility"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("utility function tests", func() {
	Context("Testing GetSum", func() {
		It("sum of non-empty array", func() {
			arr := [5]int{1, 2, 4, 5, 3}
			s := utility.GetSum(arr[:])
			Expect(s).To(Equal(15))
		})

		It("sum of empty array", func() {
			arr := [0]int{}
			s := utility.GetSum(arr[:])
			Expect(s).To(Equal(0))
		})
	})

	Context("Testing convertString function)", func() {
		It("pass a valid string date", func() {
			s := "2016-02-08"
			_, err := utility.ConvertStringDateToTime(s)
			Expect(err).To(BeNil())
		})

		It("pass an invalid string date", func() {
			s := "2016-02*04"
			d, err := utility.ConvertStringDateToTime(s)
			fmt.Println(d)
			Expect(err).ToNot(BeNil())
		})
	})
})
