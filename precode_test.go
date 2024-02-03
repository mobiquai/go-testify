package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, http.StatusOK) // код ответа должен быть 200
	assert.NotEmpty(t, responseRecorder.Body.String())     // тело ответа не должно быть пустым
}

// Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerWhenWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=london", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	expected := "wrong city value"

	assert.NotEqual(t, responseRecorder.Code, http.StatusBadRequest) // если код ответа 400, то ошибка
	assert.NotEqual(t, responseRecorder.Body.String(), expected)     // если в теле текст "wrong city value", то ошибка
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 6
	req := httptest.NewRequest("GET", "/cafe?count=6&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, http.StatusOK) // код ответа должен быть 200

	retSlyce := strings.Split(responseRecorder.Body.String(), ",") // получаем слайс из строки по разделителю ","
	assert.Less(t, totalCount, len(retSlyce))                      // Если указано больше, то ошибка
}
