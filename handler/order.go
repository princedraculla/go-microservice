package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/princedraculla/go-microservice/model"
	repo "github.com/princedraculla/go-microservice/repository/order"
)

type Order struct {
	Repo *repo.RedisRepo
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID uuid.UUID        `json:"customer_id"`
		LineItem   []model.LineItem `json:"line_items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	order := model.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItem,
		CreatedAt:  &now,
	}

	err := o.Repo.Insert(r.Context(), order)

	if err != nil {
		fmt.Println("failed to Insert order:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)

	if err != nil {
		fmt.Println("fialed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")

	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50

	res, err := o.Repo.FindAll(r.Context(), repo.FindAllPage{
		Size:   size,
		Offset: cursor,
	})

	if err != nil {
		fmt.Println("fialed to find all: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}

	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one Order By ID")

	orderIDStr := chi.URLParam(r, "id")

	if orderIDStr == "" {
		orderIDStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	orderId, err := strconv.ParseUint(orderIDStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := o.Repo.FindByID(r.Context(), orderId)
	if errors.Is(err, repo.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find order:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(order); err != nil {
		fmt.Println("failed to marshal:", err)
		return
	}
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update one Order by ID")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one Order By ID")
}
