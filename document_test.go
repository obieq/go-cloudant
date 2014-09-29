package cloudant_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var autoYear int = 1960
var autoMake string = "Plymouth"
var autoModel string = "XNR"
var doc CloudantAutomobile
var docId string

var _ = Describe("Document", func() {
	BeforeEach(func() {
		doc = CloudantAutomobile{Year: autoYear, Make: autoMake, Model: autoModel}
		docId = GenerateRandomUUID()
		doc.SetId(docId)
	})

	Context("Creating", func() {
		It("should create a document", func() {
			CreateDocumentAndAssert(&doc, docId, false)
		})

		It("should create a document in batch mode", func() {
			CreateDocumentAndAssert(&doc, docId, true)
		})

		//It("should create 600 documents in batch mode", func() {
		//for i := 1; i <= 600; i++ {
		//auto := CloudantAutomobile{Year: autoYear, Make: autoMake, Model: autoModel}
		//auto.SetId(autoId + strconv.Itoa(i))
		//go testDb.CreateDocument(&auto, true)
		//}
		//time.Sleep(8000 * time.Millisecond)
		//})
	})

	Context("Getting", func() {
		It("should get a document", func() {
			CreateDocumentAndAssert(&doc, docId, true)

			err := testDb.GetDocument(docId, &doc)
			Ω(err).NotTo(HaveOccurred())
			Ω(doc).ShouldNot(BeNil())
			Ω(doc.Id()).Should(Equal(docId))
			Ω(doc.Revision()).ShouldNot(BeNil())
			Ω(doc.Year).Should(Equal(autoYear))
			Ω(doc.Make).Should(Equal(autoMake))
			Ω(doc.Model).Should(Equal(autoModel))
		})
	})

	Context("Updating", func() {
		It("should update a document", func() {
			CreateDocumentAndAssert(&doc, docId, true)

			err := testDb.GetDocument(docId, &doc)
			Ω(err).NotTo(HaveOccurred())
			Ω(doc.Revision).ShouldNot(BeNil())

			origRevision := doc.Revision()
			newAutoYear := 1936
			newAutoMake := "Bugatti"
			newAutoModel := "Type 57 Atalante"

			doc.Year = newAutoYear
			doc.Make = newAutoMake
			doc.Model = newAutoModel

			cdr, err := testDb.UpdateDocument(&doc, false)
			Ω(err).NotTo(HaveOccurred())
			Ω(cdr).ShouldNot(BeNil())
			Ω(cdr.Id).Should(Equal(docId))
			Ω(cdr.Revision).ShouldNot(BeNil())
			Ω(cdr.Revision).ShouldNot(Equal(origRevision))

			// query cloudant again in order to verify year was updated as expected
			err = testDb.GetDocument(docId, &doc)
			Ω(err).NotTo(HaveOccurred())
			Ω(doc.Year).Should(Equal(newAutoYear))
		})

		It("should update a document in batch mode", func() {
			CreateDocumentAndAssert(&doc, docId, false)

			err := testDb.GetDocument(docId, &doc)
			Ω(err).NotTo(HaveOccurred())

			origRevision := doc.Revision()
			newAutoYear := 1950

			doc.Year = newAutoYear

			cdr, err := testDb.UpdateDocument(&doc, true)
			Ω(err).NotTo(HaveOccurred())
			Ω(cdr).ShouldNot(BeNil())
			Ω(cdr.Id).Should(Equal(docId))
			Ω(cdr.Revision).ShouldNot(BeNil())
			Ω(cdr.Revision).ShouldNot(Equal(origRevision))
		})
	})

	Context("Updating", func() {
		It("should delete a document", func() {
			CreateDocumentAndAssert(&doc, docId, true)

			err := testDb.GetDocument(docId, &doc)
			cdr, err := testDb.DeleteDocument(docId, doc.Revision())
			Ω(err).NotTo(HaveOccurred())
			Ω(cdr).ShouldNot(BeNil())
			Ω(cdr.Id).Should(Equal(docId))
			Ω(cdr.Revision).ShouldNot(BeNil())
			Ω(cdr.Revision).ShouldNot(Equal(doc.Revision()))
		})
	})

	Describe("Error Handling", func() {
		Context("Creating", func() {
			It("should return an error if the document fails json.Marshal", func() {
				_, err := testDb.CreateDocument(GenerateInvalidJson(), false)
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).Should(Equal("json: unsupported type: map[int]interface {}"))
			})
		})

		Context("Updating", func() {
			It("should return an error if the document fails json.Marshal", func() {
				CreateDocumentAndAssert(&doc, docId, true)

				doc.FluxCapacitor = GenerateInvalidJson()
				_, err := testDb.UpdateDocument(&doc, false)
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).Should(Equal("json: unsupported type: map[int]interface {}"))
			})
		})
	})
})
