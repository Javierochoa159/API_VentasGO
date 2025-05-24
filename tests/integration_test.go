package tests

import (
	"API_VentasGO/api"
	"API_VentasGO/internal/sale"
	"API_VentasGO/internal/user"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestIntegrationCreateAndGet(t *testing.T) {
	app := gin.Default()
	api.InitRoutes(app)

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	res := fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.Code)
	require.Contains(t, res.Body.String(), "pong")

	req, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{
		"name":"Ayrton",
		"address": "Pringles",
		"nickname": "Chiche"	
	}`))

	res = fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusCreated, res.Code)

	var resUser *user.User
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &resUser))
	require.Equal(t, "Ayrton", resUser.Name)
	require.Equal(t, "Pringles", resUser.Address)
	require.Equal(t, "Chiche", resUser.NickName)
	require.Equal(t, 1, resUser.Version)
	require.NotEmpty(t, resUser.ID)
	require.NotEmpty(t, resUser.CreatedAt)
	require.NotEmpty(t, resUser.UpdatedAt)

	req, _ = http.NewRequest(http.MethodGet, "/users/"+resUser.ID, nil)

	res = fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusOK, res.Code)
}

func TestIntegrationPostAndPathAndGetSale(t *testing.T) {
	app := gin.Default()
	api.InitRoutes(app)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{
		"name":"Ayrton",
		"address": "Pringles",
		"nickname": "Chiche"	
	}`))

	res := fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusCreated, res.Code)

	var resUser *user.User
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &resUser))
	require.Equal(t, "Ayrton", resUser.Name)
	require.Equal(t, "Pringles", resUser.Address)
	require.Equal(t, "Chiche", resUser.NickName)
	require.Equal(t, 1, resUser.Version)
	require.NotEmpty(t, resUser.ID)
	require.NotEmpty(t, resUser.CreatedAt)
	require.NotEmpty(t, resUser.UpdatedAt)
	/*
		req, _ = http.NewRequest(http.MethodPost, "/sales", bytes.NewBufferString(`{
			"user_id": "`+resUser.ID+`",
			"amount": 15000
		}`))
	*/

	sale1 := map[string]interface{}{
		"user_id": resUser.ID,
		"amount":  15000,
	}
	t.Logf("JSON a enviar: %+v", sale1)
	jsonBody, err := json.Marshal(sale1)
	require.NoError(t, err)

	req, err = http.NewRequest(http.MethodPost, "/sales", bytes.NewReader(jsonBody))
	require.NoError(t, err)
	res = fakeRequest(app, req)

	require.NotNil(t, res)
	require.Equal(t, http.StatusCreated, res.Code)

	var resSale *sale.Sale
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &resSale))
	require.Equal(t, resUser.ID, resSale.UserId)
	require.Equal(t, 15000, resSale.Amount)
	require.NotEmpty(t, resSale.Status)
	require.Equal(t, 1, resSale.Version)
	require.NotEmpty(t, resSale.CreatedAt)
	require.NotEmpty(t, resSale.UpdatedAt)
	status := []string{"Approved", "Pending", "Rejected"}
	require.Contains(t, status, resSale.Status)
	/*
		req, _ = http.NewRequest(http.MethodGet, "/sales?user_id="+resSale.UserId, nil)

		res = fakeRequest(app, req)

		require.NotNil(t, res)
		require.Equal(t, http.StatusOK, res.Code)
	*/

}

func fakeRequest(e *gin.Engine, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	return w
}
