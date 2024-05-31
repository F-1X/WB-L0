package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"wb/backend/internal/domain/entity"
)

func (s *APITestSuite) TestGetOrder() {
	req, err := http.NewRequest("GET", "/order?id=order2", nil)
	req.Header.Set("Content-type", "application/json")
	s.NoError(err)

	resp := httptest.NewRecorder()
	s.handler.Mux.ServeHTTP(resp, req)

	r := s.Require()
	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var result entity.Order

	respData, err := io.ReadAll(resp.Body)

	s.NoError(err)
	err = json.Unmarshal(respData, &result)
	r.NoError(err)

	r.Equal(order2, result)
}
