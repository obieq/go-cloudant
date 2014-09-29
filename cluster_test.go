package cloudant_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cluster", func() {

	It("should get cluster info", func() {
		ci, err := testClient.GetClusterInfo()
		Ω(err).NotTo(HaveOccurred())
		Ω(ci.CloudantBuild).ShouldNot(BeNil())
		Ω(ci.Couchdb).ShouldNot(BeNil())
		Ω(ci.Version).ShouldNot(BeNil())
	})

})
