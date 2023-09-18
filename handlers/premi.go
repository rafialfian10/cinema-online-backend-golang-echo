package handlers

import (
	dto "cinemaonline/dto"
	"cinemaonline/models"
	"cinemaonline/repositories"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

type handlerPremi struct {
	PremiRepository repositories.PremiRepository
}

func HandlerPremi(PremiRepository repositories.PremiRepository) *handlerPremi {
	return &handlerPremi{PremiRepository}
}

// function get all premiums
func (h *handlerPremi) FindPremis(c echo.Context) error {
	premiums, err := h.PremiRepository.FindPremis()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: premiums})
}

// function get by id premium
func (h *handlerPremi) GetPremi(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	premium, err := h.PremiRepository.GetPremi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: premium})
}

// function update premium
func (h *handlerPremi) UpdatePremiumByUser(c echo.Context) error {
	price, _ := strconv.Atoi(c.FormValue("price"))
	statusForm := c.FormValue("status")
	status, _ := strconv.ParseBool(statusForm)

	request := dto.UpdatePremiRequest{
		Status: status,
		Price:  price,
	}

	validation := validator.New()
	if err := validation.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	id, _ := strconv.Atoi(c.Param("id"))
	premi, err := h.PremiRepository.GetPremi(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResult{Status: http.StatusNotFound, Message: "Transaction not found"})
	}

	var transactionIsMatch = false
	var transactionId int
	for !transactionIsMatch {
		transactionId = int(time.Now().Unix())
		transactionData, _ := h.PremiRepository.GetPremi(transactionId)
		if transactionData.ID == 0 {
			transactionIsMatch = true
		}
	}

	premi.OrderID = transactionId
	premi.Price = request.Price
	premi.Status = request.Status

	// Menyimpan perubahan ke dalam database
	premiUpdated, err := h.PremiRepository.UpdatePremiUser(premi, premi.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	// Mengupdate data premi dengan perubahan ActivatedAt dan ExpiredAt
	_, err = h.PremiRepository.UpdatePremiUser(premiUpdated, premiUpdated.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	// Mengambil data premi yang baru diupdate
	getPremiUpdated, _ := h.PremiRepository.GetPremi(premiUpdated.ID)

	// 1. Initiate Snap client
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY_TRANSACTION_PREMIUM"), midtrans.Sandbox)
	// Use to midtrans.Production if you want Production Environment (accept real transaction).

	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(premiUpdated.OrderID),
			GrossAmt: int64(premiUpdated.Price),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: premiUpdated.User.Username,
			Email: premiUpdated.User.Email,
		},
	}

	// 3. Execute request create Snap transaction to Midtrans Snap API
	snapResp, _ := s.CreateTransaction(req)

	// mengupdate token di database
	updateTokenPremi, _ := h.PremiRepository.UpdateTokenPremi(snapResp.Token, getPremiUpdated.ID)

	// mengambil data transaction yang baru diupdate
	datapremiUpdated, _ := h.PremiRepository.GetPremi(updateTokenPremi.ID)

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertPremiResponse(datapremiUpdated)})
}

// function update premium by admin
func (h *handlerPremi) UpdatePremiByAdmin(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	statusForm := c.FormValue("status")
	status, err := strconv.ParseBool(statusForm)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: "Invalid status value"})
	}

	request := dto.UpdatePremiRequest{
		Status: status,
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	_, err = h.PremiRepository.GetPremi(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResult{Status: http.StatusNotFound, Message: err.Error()})
	}

	// Update premium
	premiumUpdated, err := h.PremiRepository.UpdatePremiUserStatus(request.Status, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	// Retrieve the updated premium
	getPremiumUpdated, err := h.PremiRepository.GetPremi(premiumUpdated.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertPremiResponse(getPremiumUpdated)})
}

// function delete premium
func (h *handlerPremi) DeletePremi(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	premium, err := h.PremiRepository.GetPremi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	dataPremium, err := h.PremiRepository.DeletePremi(premium, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: dataPremium})
}

// function notification
func (h *handlerPremi) NotificationPremi(c echo.Context) error {
	var notificationPayload map[string]interface{}

	if err := c.Bind(&notificationPayload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)
	order_id, _ := strconv.Atoi(orderId)

	premi, err := h.PremiRepository.GetPremiOrderId(order_id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: "Transaksi tidak ditemukan"})
	}

	if transactionStatus == "capture" {
		if fraudStatus == "accept" {
			premi.Status = true
			h.PremiRepository.UpdatePremiUserStatus(premi.Status, premi.ID)
			SendMailPremi("Transaction Success", premi)
		} else if fraudStatus == "challenge" {
			// Transaksi dalam tantangan fraud
			h.PremiRepository.UpdatePremiUserStatus(premi.Status, premi.ID)
		}
	} else if transactionStatus == "settlement" {
		premi.Status = true
		h.PremiRepository.UpdatePremiUserStatus(premi.Status, premi.ID)
		SendMailPremi("Transaksi Success", premi)
	} else if transactionStatus == "deny" {
		// Pembayaran ditolak
		h.PremiRepository.UpdatePremiUserStatus(premi.Status, premi.ID)
		SendMailPremi("Transaction Failed", premi)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		// Pembayaran dibatalkan atau kadaluwarsa
		h.PremiRepository.UpdatePremiUserStatus(premi.Status, premi.ID)
		SendMailPremi("Transaction Canceled / Expired", premi)
	} else if transactionStatus == "pending" {
		// Transaksi masih dalam status pending
		h.PremiRepository.UpdatePremiUserStatus(premi.Status, premi.ID)
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: notificationPayload})
}

// function SendMailPremium
func SendMailPremi(status string, premi models.Premi) {
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = "Cinema Online <rafialfian770@gmail.com>"
	var CONFIG_AUTH_EMAIL = os.Getenv("SYSTEM_EMAIL")
	var CONFIG_AUTH_PASSWORD = os.Getenv("SYSTEM_PASSWORD")

	var price = strconv.Itoa(premi.Price)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", premi.User.Email)
	mailer.SetHeader("Subject", "Transaction Status")
	mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="UTF-8" />
				<meta http-equiv="X-UA-Compatible" content="IE=edge" />
				<meta name="viewport" content="width=device-width, initial-scale=1.0" />
				<title>Document</title>
				<style>
					h1 {
						color: brown;
					}
				</style>
			</head>
			<body>
				<h2>Premium payment :</h2>
				<ul style="list-style-type:none;">
					<li>Total payment: Rp.%s</li>
					<li>Status : <b>%s</b></li>
					<li>Thank you for subscribing to premium</li>
				</ul>
			</body>
		</html>`, price, status))

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent! to " + premi.User.Email)
}

// function convert premium
func ConvertPremiResponse(premi models.Premi) models.PremiResponse {
	return models.PremiResponse{
		ID:          premi.ID,
		OrderID:     premi.OrderID,
		Status:      premi.Status,
		Price:       premi.Price,
		Token:       premi.Token,
		UserID:      premi.UserID,
		ActivatedAt: premi.ActivatedAt,
		ExpiredAt:   premi.ExpiredAt,
	}
}
