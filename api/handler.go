package api

import (
	"errors"
	"net/http"
	"strings"

	"API_VentasGO/internal/metadata"
	"API_VentasGO/internal/sale"
	"API_VentasGO/internal/user"

	"github.com/gin-gonic/gin"
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
	m := &metadata.Metadata{
		Quantity:     0,
		Approved:     0,
		Pending:      0,
		Rejected:     0,
		Total_amount: float32(0),
	}

	for {
		err := h.metadataService.Create(m, u.ID)
		if err == nil {
			break
		}
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

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

// handleCreate handles POST /sales
func (h *handler) handleCreateSale(ctx *gin.Context) {
	// request payload
	var req struct {
		UserId string  `json:"user_id"`
		Amount float32 `json:"amount"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": sale.ErrInvalidAmoun})
		return
	}
	sale := &sale.Sale{
		UserId: req.UserId,
		Amount: req.Amount,
	}
	if _, err := h.userService.Get(sale.UserId); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := h.saleService.Create(sale); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for {
		_, err := h.metadataService.Update("new_sale", sale.UserId, sale.Amount)
		if err == nil {
			break
		}
	}

	ctx.JSON(http.StatusCreated, sale)
}

func checkStatus(status string) bool {
	estados := map[string]string{"pending": "", "approved": "", "rejected": "", "": ""}
	_, ok := estados[status]
	return ok
}

// handleRead handles GET /sales/:id
func (h *handler) handleReadSale(ctx *gin.Context) {
	type SaleResponse struct {
		Metadata *metadata.Metadata `json:"metadata"`
		Results  []*sale.Sale       `json:"results"`
	}

	id := ctx.Query("user_id")
	status := strings.ToLower(ctx.Query("status"))

	_, err := h.userService.Get(id)
	if errors.Is(err, user.ErrNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !checkStatus(status) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": sale.ErrInvalidStatus})
		return
	}

	sales := h.saleService.GetUserSales(id, status)

	//salesJson, _ := json.Marshal(sales)
	//fmt.Println("sales: ", salesJson)

	metadata := h.metadataService.Get(id)

	//metadataJson, _ := json.Marshal(metadata)
	//fmt.Println("metadata: ", metadataJson)

	response := SaleResponse{
		Metadata: metadata,
		Results:  sales,
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

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for {
		_, err := h.metadataService.Update(updated_sale.Status, updated_sale.UserId)
		if err == nil {
			break
		}
	}

	ctx.JSON(http.StatusOK, updated_sale)
}
