package application_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimasmir03/Calc_GO/internal/application"
	"github.com/dimasmir03/Calc_GO/pkg/calculation"
)

func TestCalculateHandleSuccess(t *testing.T) {
	successCases := []struct {
		name               string
		method             string
		expressionRequest  *application.Request
		expectedResult     float64
		expectedStatusCode int
	}{
		{
			name:               "simple",
			method:             http.MethodPost,
			expressionRequest:  &application.Request{Expression: "2+2"},
			expectedResult:     4,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "priority 1",
			method:             http.MethodPost,
			expressionRequest:  &application.Request{Expression: "2*(2+2)"},
			expectedResult:     8,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "priority 2",
			method:             http.MethodPost,
			expressionRequest:  &application.Request{Expression: "2+2*2"},
			expectedResult:     6,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "division",
			method:             http.MethodPost,
			expressionRequest:  &application.Request{Expression: "1/2"},
			expectedResult:     0.5,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range successCases {
		t.Run(testCase.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, err := json.Marshal(testCase.expressionRequest)
			if err != nil {
				fmt.Printf("Error try json decode: %s", err.Error())
				return
			}
			req := httptest.NewRequest(testCase.method, "/calculate", bytes.NewBuffer(data))
			application.CalculationHandler(w, req)
			if w.Result().StatusCode != testCase.expectedStatusCode {
				t.Fatalf("successful case \"%s\" returns invalid status code: expected=%d, but got=%d", testCase.name, http.StatusOK, w.Result().StatusCode)
			}

			data, err = io.ReadAll(w.Result().Body)
			if err != nil {
				t.Fatalf("err while read Body: %s", err.Error())
			}
			var result application.SuccessResponse
			err = json.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("err while parse Body: %s", err.Error())
			}
			if result.Result != testCase.expectedResult {
				t.Fatalf("successful case \"%s\" returns invalid result: expected=%f, but got=%f", testCase.name, testCase.expectedResult, result.Result)
			}
		})
	}
}

func TestCalculateHandleFailed(t *testing.T) {
	failedCases := []struct {
		name               string
		method             string
		expressionRequest  *application.Request
		expectedError      string
		expectedStatusCode int
	}{
		{
			name:               "invalid method",
			method:             http.MethodGet,
			expressionRequest:  &application.Request{Expression: "2+2**"},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError:      application.ErrInvalidMethod.Error(),
		},
		{
			name:               "invalid expression",
			method:             http.MethodPost,
			expressionRequest:  &application.Request{Expression: "2+2**"},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedError:      calculation.ErrInvalidExpression.Error(),
		},
		{
			name:               "mismatch parentheses",
			method:             http.MethodPost,
			expressionRequest:  &application.Request{Expression: "2*(2+2"},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedError:      calculation.ErrMismatchParentheses.Error(),
		},
		{
			name:               "division by zero",
			method:             http.MethodPost,
			expressionRequest:  &application.Request{Expression: "10/0"},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedError:      calculation.ErrDivisionByZero.Error(),
		},
		{
			name:               "invalid character",
			method:             http.MethodPost,
			expressionRequest:  &application.Request{Expression: "2a(2+2"},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedError:      calculation.ErrInvalidCharacter.Error(),
		},
	}

	for _, testCase := range failedCases {
		t.Run(testCase.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, err := json.Marshal(testCase.expressionRequest)
			if err != nil {
				t.Fatalf("Error try json decode: %s", err.Error())
			}
			req := httptest.NewRequest(testCase.method, "/calculate", bytes.NewBuffer(data))
			application.CalculationHandler(w, req)
			if w.Result().StatusCode != testCase.expectedStatusCode {
				t.Fatalf("failed case \"%s\" returns invalid status code: expected=%d, but got=%d", testCase.name, testCase.expectedStatusCode, w.Result().StatusCode)
			}

			data, err = io.ReadAll(w.Result().Body)
			if err != nil {
				t.Fatalf("err while read Body: %s", err.Error())
			}
			var result application.FailedResponse
			err = json.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("err while parse Body: %s", err.Error())
			}
			if result.Error != testCase.expectedError {
				t.Fatalf("failed case \"%s\" returns invalid err: expected=%s, but got=%s", testCase.name, testCase.expectedError, result.Error)
			}
		})
	}
}
