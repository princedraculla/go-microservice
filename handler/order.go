package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

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
	data, err := o.Repo.FindAll(r.Context(), repo.FindAllPage{
		Size:   10,
		Offset: 15,
	})

	if err != nil {
		err := fmt.Sprintf("somthing not right... %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err))
	}

	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf("orders: %v, and the cursor right now, %v", data.Orders, data.Cursor)
	w.Write([]byte(response))

}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get one Order By ID")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update one Order by ID")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one Order By ID")
}
