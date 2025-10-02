package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAccrualService(t *testing.T) {
	logger := zap.NewExample()

	tests := []struct {
		testName           string
		accrualRawResponse string
		expectedAccrual    string
	}{
		{
			testName: "accrual int",
			accrualRawResponse: `{
                "order": "<number>",
                "status": "PROCESSED",
                "accrual": 500
            }`,
			expectedAccrual: "500",
		},
		{
			testName: "accrual float 1",
			accrualRawResponse: `{
                "order": "<number>",
                "status": "PROCESSED",
                "accrual": 729.9
            }`,
			expectedAccrual: "729.9",
		},
		{
			testName: "accrual float 2",
			accrualRawResponse: `{
                "order": "<number>",
                "status": "PROCESSED",
                "accrual": 729.98
            }`,
			expectedAccrual: "729.98",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			srv := httptest.NewServer(
				http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
					res.Write([]byte(test.accrualRawResponse))
				}),
			)
			defer srv.Close()

			client := srv.Client()

			accrualService := NewAccrualService(client, logger, srv.URL)

			accrual, err := accrualService.GetAccrual(context.Background(), "123")

			require.NoError(t, err)
			assert.Equal(t, test.expectedAccrual, accrual.String())
		})
	}
}
