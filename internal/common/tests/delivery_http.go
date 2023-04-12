package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type Wrapper func(req *http.Request) *http.Request

func DeliveryTestPost(t *testing.T, r *chi.Mux, target string, requestBody string, expectedStatus int, expectedJSONResponse string,
	wrapper Wrapper) *httptest.ResponseRecorder {
	
	t.Helper()
	req := httptest.NewRequest("POST", target, bytes.NewBufferString(requestBody))
	return deliveryTest(t, r, req, expectedStatus, expectedJSONResponse, wrapper)
}

func DeliveryTestGet(t *testing.T, r *chi.Mux, target string, expectedStatus int, expectedJSONResponse string,
	wrapper Wrapper) *httptest.ResponseRecorder {
	
	t.Helper()
	req := httptest.NewRequest("GET", target, nil)
	return deliveryTest(t, r, req, expectedStatus, expectedJSONResponse, wrapper)
}

func DeliveryTestDelete(t *testing.T, r *chi.Mux, target string, expectedStatus int, expectedJSONResponse string,
	wrapper Wrapper) *httptest.ResponseRecorder {
	
	t.Helper()	
	req := httptest.NewRequest("DELETE", target, nil)
	return deliveryTest(t, r, req, expectedStatus, expectedJSONResponse, wrapper)
}

func deliveryTest(t *testing.T, r *chi.Mux, req *http.Request, expectedStatus int, expectedJSONResponse string,
	wrapper Wrapper) *httptest.ResponseRecorder {

	w := httptest.NewRecorder()

	r.ServeHTTP(w, wrapper(req))

	assert.Equal(t, expectedStatus, w.Code)
	assert.JSONEq(t, expectedJSONResponse, w.Body.String())
	return w
}
