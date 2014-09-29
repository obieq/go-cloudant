package cloudant_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var dbName string = "ginkgo_database_spec_test"

var _ = Describe("Database", func() {

	It("should create a database", func() {
		cdr, err := testClient.CreateDatabase(dbName)
		Ω(err).NotTo(HaveOccurred())
		Ω(cdr).ShouldNot(BeNil())
	})

	It("should get a database", func() {
		db := testClient.GetDatabase(dbName)
		Ω(db).ShouldNot(BeNil())
		Ω(db.Client()).ShouldNot(BeNil())
		Ω(db.Name()).Should(Equal(dbName))
	})

	It("should list all database names", func() {
		dbs, err := testClient.ListDatabases()
		Ω(err).NotTo(HaveOccurred())
		Ω(dbs).ShouldNot(BeNil())
		Ω(len(dbs)).Should(BeNumerically(">", 0))
		Ω(dbs).Should(ContainElement(dbName))
	})

	It("should delete a database", func() {
		cdr, err := testClient.DeleteDatabase(dbName)
		Ω(err).NotTo(HaveOccurred())
		Ω(cdr).ShouldNot(BeNil())
	})
})
