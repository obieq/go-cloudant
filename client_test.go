package cloudant // needs to be same namespace as code files b/c it tests unexported functions

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type nopCloser struct {
	io.Reader
}

type nopReader struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func (nopReader) Close() error { return nil }

func (nopReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error stix")
}

var _ = Describe("Client", func() {

	var (
		emptyClient  *Client        = NewClient("", os.Getenv("invalid-key"), os.Getenv("invalid-password"))
		googleClient *Client        = NewClient("https://www.google.com", os.Getenv("invalid-key"), os.Getenv("invalid-password"))
		cError       *CloudantError = &CloudantError{
			Status:     "400 Bad Request",
			StatusCode: 400,
			Code:       "bad_request",
			Detail:     "Invalid rev format",
		}
	)

	It("should set request headers", func() {
		headers := make(map[string]string)
		headers["HeaderKey"] = "HeaderValue"

		_, err := emptyClient.doRequest("GET", "/", headers, nil)
		Ω(err).To(HaveOccurred())
	})

	It("should decode without printing responses", func() {
		var cdr CloudantDocumentResponse = CloudantDocumentResponse{}

		googleClient.PrintResponses = false
		resp, err := googleClient.doRequest("GET", "/", nil, nil)
		Ω(err).NotTo(HaveOccurred())
		err = googleClient.handleResponse(resp, err, 200, &cdr)
		Ω(err).NotTo(HaveOccurred())
	})

	It("should decode and print responses", func() {
		var cdr CloudantDocumentResponse = CloudantDocumentResponse{}

		googleClient.PrintResponses = true
		resp, err := googleClient.doRequest("GET", "/", nil, nil)
		Ω(err).NotTo(HaveOccurred())
		err = googleClient.handleResponse(resp, err, 200, &cdr)
		Ω(err).NotTo(HaveOccurred())
	})

	Describe("Error Handling", func() {
		Context("doRequest", func() {
			It("should return an error if http.NewRequest fails", func() {
				_, err := emptyClient.doRequest("PUT", "%gh&%ij", nil, nil)
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).Should(Equal(`parse %gh&%ij: invalid URL escape "%gh"`))
			})
		})

		Context("handleResponse", func() {
			It("should return the same error if an error is provided for the function call", func() {
				err := emptyClient.handleResponse(nil, cError, 400, nil)
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).Should(Equal("400 Bad Request (400): bad_request Invalid rev format"))
			})

			It("should return an error if the http status code doesn't match what's expected", func() {
				response, err := googleClient.doRequest("GET", "/", nil, nil)
				err = googleClient.handleResponse(response, nil, 999, nil)
				Ω(err).To(HaveOccurred())
			})
		})

		Context("newCloudantError", func() {
			It("should return an error if ioutil.ReadAll fails", func() {
				response := &http.Response{}
				response.Body = nopReader{bytes.NewBufferString("")}

				err := newCloudantError(response)
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring("error stix"))
			})

			It("should return an error if json.Unmarshal fails", func() {
				response := &http.Response{}
				response.Body = nopCloser{bytes.NewBufferString("phablet")}

				err := newCloudantError(response)
				Ω(err).To(HaveOccurred())
				Ω(err.Error()).Should(ContainSubstring("Can not unmarshal JSON response '''phablet'''"))
			})
		})

		Context("CloudantError", func() {
			It("should have an Error() method", func() {
				Ω(cError.Error()).Should(Equal("400 Bad Request (400): bad_request Invalid rev format"))
			})
		})
	})
})
