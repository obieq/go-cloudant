package cloudant_test

import (
	"os"

	. "github.com/obieq/go-cloudant"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"
)

var (
	testClient        *Client
	testDb            *Database
	testDbName        string        = "golang_suite_test"
	errClientRequest  *Client       = NewClient(os.Getenv("invalid-url"), os.Getenv("invalid-key"), os.Getenv("invalid-password"))
	errClientResponse *Client       = NewClient(os.Getenv("CLOUDANT_URL"), os.Getenv("invalid-key"), os.Getenv("invalid-password"))
	unauthorizedError string        = "401 Unauthorized (401): unauthorized Name or password is incorrect"
	indexCreateSleep  time.Duration = 300 // in milliseconds
	queryDataCreated  bool          = false
)

type CloudantAutomobile struct {
	CloudantDocument
	Year          int
	Make          string
	Model         string
	Trim          string
	FluxCapacitor interface{}
}

func TestGoCloudant(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cloudant Suite")
}

// NOTE: only create query data once per test run
func createTestQueryDataWithIndices() {
	if queryDataCreated {
		return
	}

	// create indices
	CreateIndexAndAssert([]string{"Model"}, "", "")
	CreateIndexAndAssert([]string{"Year", "Make"}, "", "")

	lambo1 := CloudantAutomobile{Year: 2000, Make: "Lamborghini", Model: "Diablo"}
	lambo2 := CloudantAutomobile{Year: 2005, Make: "Lamborghini", Model: "Murciélago"}
	lambo3 := CloudantAutomobile{Year: 2010, Make: "Lamborghini", Model: "Gallardo", Trim: "basic"}
	lambo4 := CloudantAutomobile{Year: 2011, Make: "Lamborghini", Model: "Gallardo", Trim: "basic"}
	ferrari1 := CloudantAutomobile{Year: 2000, Make: "Ferrari", Model: "550 Maranello", Trim: "sport"}
	ferrari2 := CloudantAutomobile{Year: 2005, Make: "Ferrari", Model: "F430"}
	ferrari3 := CloudantAutomobile{Year: 2010, Make: "Ferrari", Model: "California"}
	ferrari4 := CloudantAutomobile{Year: 2010, Make: "Ferrari", Model: "458"}

	autos := []CloudantAutomobile{lambo1, lambo2, lambo3, lambo4, ferrari1, ferrari2, ferrari3, ferrari4}

	for _, auto := range autos {
		CreateDocumentAndAssert(&auto, "", false)
	}

	queryDataCreated = true
}

func NewTestClient(printResponses bool) *Client {
	c := NewClient(os.Getenv("CLOUDANT_URL"), os.Getenv("CLOUDANT_API_KEY"), os.Getenv("CLOUDANT_API_PASSWORD"))
	c.PrintResponses = printResponses // print all Cloudant responses... helps with debugging new features

	return c
}

var _ = BeforeSuite(func() {
	testClient = NewTestClient(true)

	// delete database from prior test
	_, err := testClient.DeleteDatabase(testDbName)

	// create a new database to be used during all tests
	_, err = testClient.CreateDatabase(testDbName)
	Expect(err).NotTo(HaveOccurred())

	// get an instance of the new database
	testDb = NewDatabase(testDbName, testClient)
	Ω(testDb).ShouldNot(BeNil())
	Ω(testDb.Client()).ShouldNot(BeNil())
	Ω(testDb.Name()).Should(Equal(testDbName))
})

var _ = AfterSuite(func() {
})
