package cloudant_test

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"time"

	. "github.com/onsi/gomega"
)

// Create an index and assert it was successful.
// sleepInMilliseconds should be set if the index will be queried with the same test
// b/c it can take some time before Cloudant's GET (all indices) will return the new index
func CreateIndexAndAssert(fields []string, indexName string, ddocName string) {
	opts := make(map[string]string)

	if indexName != "" {
		opts["index_name"] = indexName
	}

	if ddocName != "" {
		opts["ddoc_name"] = ddocName
	}

	cdr, err := testDb.CreateIndex(fields, opts)
	Ω(err).NotTo(HaveOccurred())
	Ω(cdr.Result).Should(Equal("created"))

	time.Sleep(indexCreateSleep * time.Millisecond)
}

// Create a Document and assert it was successful.
func CreateDocumentAndAssert(doc interface{}, id string, isBatch bool) {
	cdr, err := testDb.CreateDocument(doc, isBatch)
	Ω(err).NotTo(HaveOccurred())
	Ω(cdr).ShouldNot(BeNil())
	Ω(cdr.Revision).ShouldNot(BeNil())
	if len(id) > 0 {
		Ω(cdr.Id).Should(Equal(id))
	}
}

// Generate a random string
func GenerateRandomString(size int) string {
	rb := make([]byte, size)
	rand.Read(rb)

	return base64.URLEncoding.EncodeToString(rb)
}

// Generate a random UUID
func GenerateRandomUUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	uuid[8] = 0x80
	uuid[4] = 0x40

	return hex.EncodeToString(uuid)
}

// Returns a map that fails json marshalling
func GenerateInvalidJson() map[string]interface{} {
	outer := make(map[string]interface{})
	inner := make(map[int]interface{})
	inner[6] = 6 // json.Marshal will fail b/c key can't be an int
	outer["random-key"] = inner

	return outer
}
