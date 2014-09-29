package cloudant_test

import (
	"time"

	. "github.com/obieq/go-cloudant"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Index", func() {
	Describe("Getting", func() {
		var (
			idx_index_name        string   = "idx_get_month_day"
			idx_index_name_fields []string = []string{"GetMonth", "GetDay"}

			idx_ddoc_name        string   = "idx_get_hour"
			idx_ddoc_name_fields []string = []string{"GetHour"}

			idx_ddoc_and_index_name        string   = "idx_get_minute_second"
			idx_ddoc_and_index_name_fields []string = []string{"GetMinute", "GetSecond"}
		)

		It("should create indices for all Getting tests", func() {
			CreateIndexAndAssert(idx_index_name_fields, idx_index_name, "")
			CreateIndexAndAssert(idx_ddoc_name_fields, "", idx_ddoc_name)
			CreateIndexAndAssert(idx_ddoc_and_index_name_fields, idx_ddoc_and_index_name, idx_ddoc_and_index_name)
		})

		Context("With no criteria", func() {
			It("should get all indices", func() {
				indices, err := testDb.GetIndices()
				Ω(err).NotTo(HaveOccurred())
				Ω(indices).ShouldNot(BeNil())
				Ω(len(indices)).Should(BeNumerically(">", 0))
			})
		})

		Context("With index name", func() {
			It("should get index", func() {
				index, err := testDb.GetIndexByName(idx_index_name)
				Ω(err).NotTo(HaveOccurred())
				Ω(index).ShouldNot(BeNil())
				Ω(index.Name).Should(Equal(idx_index_name))
				Ω(index.DDocId).ShouldNot(BeNil())
				Ω(index.DDocId).ShouldNot(Equal(idx_index_name))
				Ω(index.Type).Should(Equal("json"))

				// verify indexed fields
				// NOTE: Cloudant preserves the order in which they were specified when the index was created
				Ω(len(index.Definition.Fields)).Should((Equal(len(idx_index_name_fields))))
				Ω(index.Definition.Fields[0]).Should(HaveKey(idx_index_name_fields[0]))
				Ω(index.Definition.Fields[1]).Should(HaveKey(idx_index_name_fields[1]))
			})
		})

		Context("With ddoc name", func() {
			It("should get index", func() {
				ddoc, err := testDb.GetIndexByDDocName(idx_ddoc_and_index_name)
				Ω(err).NotTo(HaveOccurred())
				Ω(ddoc).ShouldNot(BeNil())
				Ω(ddoc.Id()).Should(Equal("_design/" + idx_ddoc_and_index_name))

				// TODO: move this validation logic to the design_document_test.go file
				Ω(ddoc.Language).Should(Equal("query"))
				Ω(ddoc.Views).ShouldNot(BeNil())
				Ω(len(ddoc.Views)).Should(Equal(1))
				view := ddoc.Views[idx_ddoc_and_index_name]
				Ω(view).ShouldNot(BeNil())
				Ω(view.Reduce).Should(Equal("_count"))
				Ω(view.Map).ShouldNot(BeNil())
				Ω(len(view.Map.Fields)).Should(Equal(2))
				Ω(view.Options).ShouldNot(BeNil())
				Ω(view.Options.W).Should(Equal(2))
				Ω(view.Options.Definition).ShouldNot(BeNil())
				Ω(view.Options.Definition.Fields).ShouldNot(BeNil())
				Ω(len(view.Options.Definition.Fields)).Should(Equal(2))
			})
		})

		Context("With non-exsistent index name", func() {
			It("should return a 404 error", func() {
				index, err := testDb.GetIndexByName("does_not_exist")
				Ω(err).To(HaveOccurred())
				Ω(err.(*CloudantError).StatusCode).Should(Equal(404))
				Ω(index).Should(BeNil())
			})
		})
	})

	Describe("Creating", func() {
		var (
			index_no_name_fields []string = []string{"CreateYear"}

			idx_index_name        string   = "idx_create_month_day"
			idx_index_name_fields []string = []string{"CreateMonth", "CreateDay"}

			idx_ddoc_name        string   = "idx_create_hour"
			idx_ddoc_name_fields []string = []string{"CreateHour"}

			idx_ddoc_and_index_name        string   = "idx_create_minute_second"
			idx_ddoc_and_index_name_fields []string = []string{"CreateMinute", "CreateSecond"}
		)

		Context("With default ddoc name and index name", func() {
			It("should create", func() {
				CreateIndexAndAssert(index_no_name_fields, "", "")
			})
		})

		Context("With index name", func() {
			It("should create", func() {
				CreateIndexAndAssert(idx_index_name_fields, idx_index_name, "")

				// verify
				index, err := testDb.GetIndexByName(idx_index_name)
				Ω(err).NotTo(HaveOccurred())
				Ω(index.Name).Should(Equal(idx_index_name))
				Ω(index.DDocId).ShouldNot(Equal(idx_index_name))
				Ω(len(index.Definition.Fields)).Should(Equal(len(idx_index_name_fields)))
			})
		})

		Context("With ddoc name", func() {
			It("should create", func() {
				CreateIndexAndAssert(idx_ddoc_name_fields, "", idx_ddoc_name)

				// verify
				ddoc, err := testDb.GetDesignDocument(idx_ddoc_name)
				Ω(err).NotTo(HaveOccurred())
				Ω(ddoc.Id()).Should(Equal("_design/" + idx_ddoc_name))

				// getting view by key = ddoc name should fail
				view := ddoc.Views[idx_ddoc_name]
				Ω(len(view.Map.Fields)).Should(Equal(0)) // should not find view by index name

				// get view by randomly generated index name and verify number of fields
				keys := ddoc.ViewKeys()
				Ω(len(ddoc.Views[keys[0]].Map.Fields)).Should(Equal(len(idx_ddoc_name_fields)))
			})
		})

		Context("With ddoc name and index name", func() {
			It("should create", func() {
				CreateIndexAndAssert(idx_ddoc_and_index_name_fields, idx_ddoc_and_index_name, idx_ddoc_and_index_name)

				// verify
				index, err := testDb.GetIndexByName(idx_ddoc_and_index_name)
				Ω(err).NotTo(HaveOccurred())
				Ω(index.Name).Should(Equal(idx_ddoc_and_index_name))
				Ω(index.DDocId).Should(Equal("_design/" + idx_ddoc_and_index_name))
				Ω(len(index.Definition.Fields)).Should(Equal(len(idx_ddoc_and_index_name_fields)))
			})
		})

		It("should not create an index if an index by the same name exists", func() {
			opts := make(map[string]string)
			opts["index_name"] = idx_ddoc_and_index_name
			opts["ddoc_name"] = idx_ddoc_and_index_name
			cdr, err := testDb.CreateIndex(idx_ddoc_and_index_name_fields, opts)
			Ω(err).NotTo(HaveOccurred())
			Ω(cdr.Result).Should(Equal("exists"))
		})
	})

	Describe("Deleting", func() {
		var (
			idx_index_name        string   = "idx_delete"
			idx_index_name_fields []string = []string{"DeleteYear"}
		)

		Context("With index name", func() {
			It("should delete index", func() {
				// create index
				CreateIndexAndAssert(idx_index_name_fields, idx_index_name, "")

				// delete index
				cdr, err := testDb.DeleteIndexByName(idx_index_name)
				Ω(err).NotTo(HaveOccurred())
				Ω(cdr).ShouldNot(BeNil())

				// verify index was deleted
				time.Sleep(200 * time.Millisecond) // wait for delete operation to complete on Cloudant
				_, err = testDb.GetIndexByName(idx_index_name)
				Ω(err).To(HaveOccurred())
			})
		})

		Context("With non-existent index name", func() {
			It("should return a 404 error", func() {
				_, err := testDb.DeleteIndexByName("does_not_exist")
				Ω(err).To(HaveOccurred())
				Ω(err.(*CloudantError).StatusCode).Should(Equal(404))
			})
		})
	})

	Describe("Error Handling", func() {
		Context("Getting index by name", func() {
			It("should return an error if the http request fails", func() {
				db := errClientRequest.GetDatabase("non-existent-db-name")
				_, err := db.GetIndexByName("non-existent-index-name")
				Ω(err).To(HaveOccurred())
			})
		})
	})
})
