package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
	"moscow": {"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	reqCount := 10

	url := fmt.Sprintf("/cafe?count=%d&city=moscow", reqCount)
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	respCafeList := strings.Split(responseRecorder.Body.String(), ",")

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code 200")
	assert.Equal(t, cafeList["moscow"], respCafeList, "Response cafe list should contain all available cafes")

}

func TestMainHandlerWhenOkAndBodyNotEmpty(t *testing.T) {

	url := "/cafe?count=4&city=moscow"
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status Ok (200)")
	assert.NotEmpty(t, responseRecorder.Body, "Response body should not be empty")

}

func TestMainHandlerWhenWrongCity(t *testing.T) {

	url := "/cafe?count=4&city=novocherkassk"
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status BadRequest (400)")
	assert.Equal(t, "wrong city value", responseRecorder.Body.String(), "Expected error: 'wrong city value'")
}

func TestMainHandlerWhenCountWrong(t *testing.T) {

	url := "/cafe?count=яндексвозьмитенаработупжпж&city=moscow"
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected BadRequest status (400)")
	assert.Equal(t, "wrong count value", responseRecorder.Body.String(), "Expected error: 'wrong count value'")
}

func TestMainHandlerWhenCountMissing(t *testing.T) {

	url := "/cafe?city=moscow"
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected BadRequest status (400)")
	assert.Equal(t, "count missing", responseRecorder.Body.String(), "Expected error: 'count missing'")
}

func TestMainHandlerWhenCityMissing(t *testing.T) {

	url := "/cafe?count=4"
	req := httptest.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected BadRequest status (400)")
	assert.Equal(t, "wrong city value", responseRecorder.Body.String(), "Expected error: 'wrong city value'")
}
