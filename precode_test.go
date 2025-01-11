package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	reqCount := 10

	url := fmt.Sprintf("/cafe?count=%d&city=moscow", reqCount)
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	respCafeList := strings.Split(responseRecorder.Body.String(), ",")

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code 200")
	assert.Len(t, respCafeList, totalCount, "Expected a full list of cafes")
}

func TestMainHandlerWhenOkAndBodyNotEmpty(t *testing.T) {

	url := "/cafe?count=4&city=moscow"
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status Ok (200)")
	assert.NotEmpty(t, responseRecorder.Body, "Response body should not be empty")

}

func TestMainHandlerWhenWrongCity(t *testing.T) {

	url := "/cafe?count=4&city=novocherkassk"
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status BadRequest (400)")
	assert.Equal(t, "wrong city value", responseRecorder.Body.String(), "Expected error: 'wrong city value'")
}

func TestMainHandlerWhenCountWrong(t *testing.T) {

	url := "/cafe?count=яндексвозьмитенаработупжпж&city=moscow"
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected BadRequest status (400)")
	assert.Equal(t, "wrong count value", responseRecorder.Body.String(), "Expected error: 'wrong count value'")
}

func TestMainHandlerWhenCountMissing(t *testing.T) {

	url := "/cafe?city=moscow"
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected BadRequest status (400)")
	assert.Equal(t, "count missing", responseRecorder.Body.String(), "Expected error: 'count missing'")
}

func TestMainHandlerWhenCityMissing(t *testing.T) {

	url := "/cafe?count=4"
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected BadRequest status (400)")
	assert.Equal(t, "wrong city value", responseRecorder.Body.String(), "Expected error: 'wrong city value'")
}
