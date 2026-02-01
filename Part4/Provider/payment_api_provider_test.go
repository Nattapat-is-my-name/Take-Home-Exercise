package provider_test

import (
	"bytes"
	"io"
	"net/http"

	provider "example.com/m/Provider"

	mock_http "example.com/m/spec/support/fake"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

var _ = Describe("CallPaymentAPI", func() {
	var (
		mockCtrl           *gomock.Controller
		mockHTTPClient     *mock_http.MockHTTPClient
		paymentApiProvider provider.PaymentAPIProvider
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockHTTPClient = mock_http.NewMockHTTPClient(mockCtrl)
		paymentApiProvider = provider.NewPaymentProvider(mockHTTPClient)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Fail CallPaymentAPI", func() {
		It("should return error when call external api was timeout", func() {

			mockHTTPClient.EXPECT().Do(gomock.Any()).Return(nil, http.ErrHandlerTimeout)

			paymentResponse, sysErr := paymentApiProvider.CallPaymentAPI(provider.PaymentRequest{})

			Expect(paymentResponse).To(BeNil())
			Expect(sysErr).ToNot(BeNil())
			Expect(sysErr.ErrorCode).To(Equal("ERR2"))
			Expect(sysErr.ErrorMessage).To(Equal("HTTP request failed"))
		})

		It("should return error when external api return is not ok", func() {

			mockResp := &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       http.NoBody,
			}

			mockHTTPClient.EXPECT().Do(gomock.Any()).Return(mockResp, nil)

			paymentResponse, sysErr := paymentApiProvider.CallPaymentAPI(provider.PaymentRequest{})

			Expect(paymentResponse).To(BeNil())
			Expect(sysErr).ToNot(BeNil())
			Expect(sysErr.ErrorCode).To(Equal("ERR3"))
			Expect(sysErr.ErrorMessage).To(Equal("external api returned status 401"))
		})

	})

	Context("Success CallPaymentAPI", func() {
		It("should return payment response when call external api success", func() {

			mockResp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"transaction_id": "12345", "status": "success" }`)),
			}

			mockHTTPClient.EXPECT().Do(gomock.Any()).Return(mockResp, nil)

			paymentResponse, sysErr := paymentApiProvider.CallPaymentAPI(provider.PaymentRequest{})

			Expect(sysErr).To(BeNil())
			Expect(paymentResponse).ToNot(BeNil())
			Expect(paymentResponse.TransactionID).To(Equal("12345"))
			Expect(paymentResponse.Status).To(Equal("success"))
		})
	})

})
