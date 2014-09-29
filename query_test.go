package cloudant_test

import (
	. "github.com/obieq/go-cloudant"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query", func() {
	var (
		query        *Query
		sortingQuery *Query
		operator     map[string]interface{}
		results      []CloudantAutomobile
	)

	BeforeEach(func() {
		// create test query data (will only be performed once per test run)
		createTestQueryDataWithIndices()

		query = NewQuery()
		operator = make(map[string]interface{})
		results = []CloudantAutomobile{}
	})

	Describe("Querying", func() {
		Context("With a non-indexed field", func() {
			It("should return a 400 error", func() {
				query.Selector["non-existent-field-name"] = 6
				err := testDb.Query(query, &results)
				Ω(err).To(HaveOccurred())
				Ω(err.(*CloudantError).StatusCode).Should(Equal(400))
				Ω(err.(*CloudantError).Code).Should(Equal("no_usable_index"))
			})
		})

		Context("With one indexed field", func() {
			It("should find zero documents", func() {
				query.Selector["Model"] = "does-not-exist"
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(0))
			})

			It("should find one document", func() {
				query.Selector["Model"] = "Diablo"
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(1))
			})

			It("should find two documents", func() {
				query.Selector["Model"] = "Gallardo"
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(2))
			})

			It("should specify explicit operators", func() {
				operator["$eq"] = "Gallardo"
				query.Selector["Model"] = operator
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(2))
			})

			It("should limit", func() {
				operator["$gt"] = "1111"
				query.Selector["Model"] = operator
				query.Limit = 6
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(6))
			})

			It("should skip", func() {
				operator["$gt"] = "1111"
				query.Selector["Model"] = operator
				query.Skip = 6
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(2))
			})
		})

		Context("With one indexed field and one non-indexed field", func() {
			It("should find zero documents", func() {
				query.Selector["Model"] = "458"
				query.Selector["Trim"] = "sport"
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(0))
			})
			It("should find one document", func() {
				query.Selector["Model"] = "550 Maranello"
				query.Selector["Trim"] = "sport"
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(1))
			})

			It("should find two documents", func() {
				query.Selector["Model"] = "Gallardo"
				query.Selector["Trim"] = "basic"
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(2))
			})
		})

		Context("With two indexed fields", func() {
			It("should find zero documents", func() {
				query.Selector["Year"] = 2000
				query.Selector["Make"] = "Packard"
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(0))
			})

			It("should find one document", func() {
				query.Selector["Year"] = 2000
				query.Selector["Make"] = "Lamborghini"
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(1))
			})

			It("should find two documents", func() {
				query.Selector["Year"] = 2010
				query.Selector["Make"] = "Ferrari"
				err := testDb.Query(query, &results)
				Ω(err).NotTo(HaveOccurred())
				Ω(len(results)).Should(Equal(2))
			})
		})
	})

	Describe("Sorting", func() {
		BeforeEach(func() {
			sortingQuery = NewQuery()
			operator2 := make(map[string]interface{})

			operator["$gt"] = 1800
			operator2["$gt"] = "AA"
			sortingQuery.Selector["Year"] = operator
			sortingQuery.Selector["Make"] = operator2
		})

		It("should sort one field ascending", func() {
			sortingQuery.Sort("Year", true)
			err := testDb.Query(sortingQuery, &results)
			Ω(err).NotTo(HaveOccurred())
			Ω(len(results)).Should(Equal(8))
			Ω(results[0].Year).Should(Equal(2000))
			Ω(results[7].Year).Should(Equal(2011))
		})

		It("should sort one field descending", func() {
			sortingQuery.Sort("Year", false)
			err := testDb.Query(sortingQuery, &results)
			Ω(err).NotTo(HaveOccurred())
			Ω(len(results)).Should(Equal(8))
			Ω(results[0].Year).Should(Equal(2011))
			Ω(results[7].Year).Should(Equal(2000))
		})
	})

	Describe("Error Handling", func() {
		It("should return an error if selector fails json.Marshal", func() {
			query.Selector = GenerateInvalidJson()
			err := testDb.Query(query, &results)
			Ω(err).To(HaveOccurred())
			Ω(err.Error()).Should(Equal("json: unsupported type: map[int]interface {}"))
		})

		It("should return an error if the http request fails", func() {
			db := errClientRequest.GetDatabase("non-existent-db-name")
			query.Selector["Model"] = "Diablo"
			err := db.Query(query, &results)
			Ω(err).To(HaveOccurred())
		})
	})
})
