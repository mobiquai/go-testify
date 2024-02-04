package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createResponseRecorder(count, city string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/cafe", nil)
	query := req.URL.Query()
	query.Add("count", count)
	query.Add("city", city)
	req.URL.RawQuery = query.Encode()

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	return responseRecorder
}

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandlerWhenOk(t *testing.T) {
	responseRecorder := createResponseRecorder("2", "moscow")

	require.Equal(t, http.StatusOK, responseRecorder.Code) // код ответа должен быть 200, иначе нет смысла продолжать тест
	assert.NotEmpty(t, responseRecorder.Body.String())     // тело ответа не должно быть пустым
}

// Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerWhenWrongCity(t *testing.T) {
	responseRecorder := createResponseRecorder("3", "london")

	expected := "wrong city value"

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code) // код ответа должен быть 400
	assert.Equal(t, expected, responseRecorder.Body.String())     // если в теле текст "wrong city value", то ошибка
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4 // Необходимое кол-во для сравнения

	responseRecorder := createResponseRecorder("5", "moscow")

	require.Equal(t, http.StatusOK, responseRecorder.Code) // код ответа должен быть 200, иначе нет смысла продолжать тест

	retSlyce := strings.Split(responseRecorder.Body.String(), ",") // получаем слайс из строки по разделителю ","
	assert.Equal(t, totalCount, len(retSlyce))                     // Если указано меньше 4, то ошибка
}
