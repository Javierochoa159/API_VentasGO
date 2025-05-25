package api

import (
	"API_VentasGO/internal/metadata"
	"API_VentasGO/internal/sale"
	"API_VentasGO/internal/user"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// handler holds the user service and implements HTTP handlers for user CRUD.
type handler struct {
	userService     *user.Service
	saleService     *sale.Service
	metadataService *metadata.Service
}

// handleCreate handles POST /users
func (h *handler) handleCreateUser(ctx *gin.Context) {
	// request payload
	var req struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		NickName string `json:"nickname"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := &user.User{
		Name:     req.Name,
		Address:  req.Address,
		NickName: req.NickName,
	}
	if err := h.userService.Create(u); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, u)
}

// handleRead handles GET /users/:id
func (h *handler) handleReadUser(ctx *gin.Context) {
	id := ctx.Param("id")

	u, err := h.userService.Get(id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, u)
}

// handleUpdate handles PUT /users/:id
func (h *handler) handleUpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	// bind partial update fields
	var fields *user.UpdateFields
	if err := ctx.ShouldBindJSON(&fields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.userService.Update(id, fields)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, u)
}

// handleDelete handles DELETE /users/:id
func (h *handler) handleDeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.userService.Delete(id); err != nil {
		if errors.Is(err, user.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func validarUsuario(id string, engine *gin.Engine) (int, error) {
	if os.Getenv("MODO") != "testing" {
		resp, err := http.Get("http://localhost:9090/users/" + id)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		return resp.StatusCode, nil
	} else {
		resp, err := http.NewRequest(http.MethodGet, "http://localhost:9090/users/"+id, bytes.NewReader(nil))
		req := httptest.NewRecorder()
		engine.ServeHTTP(req, resp)
		defer resp.Body.Close()
		return req.Code, err
	}
}

// handleCreate handles POST /sales
func (h *handler) handleCreateSale(ctx *gin.Context) {
	// request payload
	var req struct {
		UserId string  `json:"user_id"`
		Amount float32 `json:"amount"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.saleService.Logger.Error("error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Amount <= 0 {
		h.saleService.Logger.Error("error", zap.Error(sale.ErrInvalidAmoun))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": sale.ErrInvalidAmoun})
		return
	}
	sale := &sale.Sale{
		UserId: req.UserId,
		Amount: req.Amount,
	}
	var errCode int
	var err error
	if os.Getenv("MODO") == "testing" {
		engine, exists := ctx.Get("engine")
		if !exists {
			ctx.JSON(500, gin.H{"error": "Invalid engine type"})
		}

		eng, ok := engine.(*gin.Engine)
		if !ok {
			ctx.JSON(500, gin.H{"error": "Invalid engine type"})
		}
		errCode, err = validarUsuario(sale.UserId, eng)
	} else {
		errCode, err = validarUsuario(sale.UserId, nil)
	}

	if err != nil {
		h.saleService.Logger.Error("error", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errCode == http.StatusNotFound {
		h.saleService.Logger.Error("error", zap.Error(errors.New("user not found")))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return

	}

	if err := h.saleService.Create(sale); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, sale)
}

func checkStatus(status string) bool {
	estados := map[string]string{"pending": "", "approved": "", "rejected": "", "": ""}
	_, ok := estados[status]
	return ok
}

// handleRead handles GET /sales?user_id
func (h *handler) handleReadSale(ctx *gin.Context) {
	type SaleResponse struct {
		Metadata *metadata.Metadata `json:"metadata"`
		Results  []*sale.Sale       `json:"results"`
	}

	id := ctx.Query("user_id")
	status := strings.ToLower(ctx.Query("status"))

	if !checkStatus(status) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": sale.ErrInvalidStatus})
		return
	}
	_, err := h.userService.Get(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	m := &metadata.Metadata{
		Quantity:     0,
		Approved:     0,
		Pending:      0,
		Rejected:     0,
		Total_amount: float32(0),
	}

	for {
		err := h.metadataService.Create(m, id)
		if err == nil {
			break
		}
	}

	sales, meta := h.saleService.GetUserSales(id, status)
	m.Quantity = int(meta["quantity"])
	m.Approved = int(meta["approved"])
	m.Rejected = int(meta["rejected"])
	m.Pending = int(meta["pending"])
	m.Total_amount = meta["total_amount"]
	response := SaleResponse{
		Metadata: m,
	}
	if sales != nil {
		response.Results = sales
	} else {
		response.Results = make([]*sale.Sale, 0)
	}

	ctx.JSON(http.StatusOK, response)
}

// handleUpdate handles PATH /sale/:id
func (h *handler) handleUpdateSale(ctx *gin.Context) {
	id := ctx.Param("id")

	// bind partial update fields
	var fields *sale.UpdateFields
	if err := ctx.ShouldBindJSON(&fields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated_sale, err := h.saleService.Update(id, fields)
	if err != nil {
		if errors.Is(err, sale.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, sale.ErrInvalidStatus) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	/*for {
		_, err := h.metadataService.Update(updated_sale.Status, updated_sale.UserId)
		if err != nil {
			continue
		}
		break
	}*/

	ctx.JSON(http.StatusOK, updated_sale)
}

// handleRead handles GET /sale/:id
func (h *handler) handleReadOneSale(ctx *gin.Context) {
	id := ctx.Param("id")

	s, err := h.saleService.Get(id)
	if err != nil {
		if errors.Is(err, sale.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, s)
}
