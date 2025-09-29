package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateShortBatch(t *testing.T) {
	tests := []struct {
		testName string
		number   string
		isValid  bool
	}{
		{
			testName: "not numeric",
			number:   "foobar",
			isValid:  false,
		},

		{
			testName: "invalid from wiki",
			number:   "4561261212345464",
			isValid:  false,
		},
		{
			testName: "valid from wiki",
			number:   "4561261212345467",
			isValid:  true,
		},

		{
			testName: "invalid from TestGophermart",
			number:   "12345678902",
			isValid:  false,
		},

		{
			testName: "valid from rosettacode",
			number:   "49927398716",
			isValid:  true,
		},
		{
			testName: "invalid from rosettacode",
			number:   "49927398717",
			isValid:  false,
		},
		{
			testName: "invalid from rosettacode",
			number:   "1234567812345678",
			isValid:  false,
		},
		{
			testName: "valid from rosettacode",
			number:   "1234567812345670",
			isValid:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			assert.Equal(t, test.isValid, IsValidLuhnNumber(test.number))
		})
	}
}
