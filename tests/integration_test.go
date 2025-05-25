package tests

import (
	"API_VentasGO/api"
	"API_VentasGO/internal/sale"
	"API_VentasGO/internal/user"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	go func() {
		//gin.setMode(gin.TestMode)
		r := gin.Default()
		api.InitRoutes(r)
		r.Run(":9090") // Servidor real
	}()

	time.Sleep(1 * time.Second) // Dale tiempo al servidor
	os.Exit(m.Run())
}

func TestIntegrationCreateAndGet(t *testing.T) {
	app := gin.Default()
	api.InitRoutes(app)

	resp, err := http.Post("http://localhost:9090/users", "application/json",
		bytes.NewBufferString(`{
		"name":"Ayrton",
		"address":"Pringles",
		"nickname":"Chiche"
	}`))

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var resUser user.User
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&resUser))
	defer resp.Body.Close()
	require.Equal(t, "Ayrton", resUser.Name)
	require.Equal(t, "Pringles", resUser.Address)
	require.Equal(t, "Chiche", resUser.NickName)
	require.Equal(t, 1, resUser.Version)
	require.NotEmpty(t, resUser.ID)
	require.NotEmpty(t, resUser.CreatedAt)
	require.NotEmpty(t, resUser.UpdatedAt)

	resp2, err := http.Get("http://localhost:9090/users/" + resUser.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var getUser user.User
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&getUser))
	defer resp2.Body.Close()

	require.Equal(t, resUser.ID, getUser.ID)
}

func engineMiddleware(engine *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("engine", engine)
		c.Next()
	}
}

func TestIntegrationPostAndPathAndGetSale(t *testing.T) {
	app := gin.Default()
	app.Use(engineMiddleware(app))
	api.InitRoutes(app)

	reqUser := map[string]interface{}{
		"name":     "Ayrton",
		"address":  "Pringles",
		"nickname": "Chiche",
	}

	jsonUser, _ := json.Marshal(reqUser)
	resq, err := http.NewRequest(http.MethodPost, "http://localhost:9090/users", bytes.NewReader(jsonUser))
	resp := httptest.NewRecorder()
	app.ServeHTTP(resp, resq)

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.Code)

	var resUser user.User
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&resUser))
	defer resq.Body.Close()
	fmt.Println("USUARIO: ", resUser)
	require.Equal(t, "Ayrton", resUser.Name)
	require.Equal(t, "Pringles", resUser.Address)
	require.Equal(t, "Chiche", resUser.NickName)
	require.Equal(t, 1, resUser.Version)
	require.NotEmpty(t, resUser.ID)
	require.NotEmpty(t, resUser.CreatedAt)
	require.NotEmpty(t, resUser.UpdatedAt)

	os.Setenv("MODO", "testing")

	saleData := map[string]interface{}{
		"user_id": resUser.ID,
		"amount":  15000,
	}
	jsonSale, _ := json.Marshal(saleData)
	resq, err = http.NewRequest(http.MethodPost, "http://localhost:9090/sales", bytes.NewReader(jsonSale))
	resp = httptest.NewRecorder()
	app.ServeHTTP(resp, resq)
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.Code)

	var resSale sale.Sale
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&resSale))
	fmt.Println("VENTA ACTUAL: ", resSale)
	defer resq.Body.Close()
	require.Equal(t, resUser.ID, resSale.UserId)
	require.Equal(t, float32(15000), resSale.Amount)
	require.NotEmpty(t, resSale.Status)
	require.Equal(t, "pending", resSale.Status)
	require.Equal(t, 1, resSale.Version)
	require.NotEmpty(t, resSale.CreatedAt)
	require.NotEmpty(t, resSale.UpdatedAt)

	saleDataToUpdate := map[string]string{
		"status": "approved",
	}

	reqJson, _ := json.Marshal(saleDataToUpdate)
	resq, err = http.NewRequest(http.MethodPatch, "/sales/"+resSale.ID, bytes.NewReader(reqJson))
	resp = httptest.NewRecorder()
	app.ServeHTTP(resp, resq)

	var resSaleUpdated sale.Sale
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&resSaleUpdated))
	fmt.Println("VENTA ACTUALIZADA: ", resSaleUpdated)
	defer resq.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Code)
	require.NotEmpty(t, resSaleUpdated.Status)
	require.NotEqual(t, "pending", resSaleUpdated.Status)
	status := []string{"approved", "rejected"}
	require.Contains(t, status, resSaleUpdated.Status)
}
