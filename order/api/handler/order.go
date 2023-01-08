package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bosamatheus/order-book/order/api/presenter"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
)

type OrderHandler struct {
	producerConn *kafka.Conn
}

func NewOrderHandler(conn *kafka.Conn) *OrderHandler {
	return &OrderHandler{
		producerConn: conn,
	}
}

func (h *OrderHandler) SendOrder(c echo.Context) error {
	o := &presenter.Order{}
	if err := c.Bind(o); err != nil {
		return err
	}
	err := h.produceMessage(o)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusAccepted, o)
}

func (h *OrderHandler) produceMessage(o *presenter.Order) error {
	key := []byte(strconv.FormatInt(o.WalletID, 10))
	val, err := json.Marshal(o)
	if err != nil {
		return err
	}
	_, err = h.producerConn.WriteMessages(
		kafka.Message{Key: key, Value: val},
	)
	return err
}
