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
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

type handlerPremium struct {
	PremiumRepository repositories.PremiumRepository
}

func HandlerPremium(PremiumRepository repositories.PremiumRepository) *handlerPremium {
	return &handlerPremium{PremiumRepository}
}

// function get all premiums
func (h *handlerPremium) FindPremiums(c echo.Context) error {
	premiums, err := h.PremiumRepository.FindPremiums()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: premiums})
}

// function get by id premium
func (h *handlerPremium) GetPremium(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	premium, err := h.PremiumRepository.GetPremium(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: premium})
}

// function create premium
func (h *handlerPremium) CreatePremium(c echo.Context) error {
	request := new(dto.CreatePremiumRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	request.Status = true

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	var transactionIsMatch = false
	var transactionId int
	for !transactionIsMatch {
		transactionId = int(time.Now().Unix())
		transactionData, _ := h.PremiumRepository.GetPremium(transactionId)
		if transactionData.ID == 0 {
			transactionIsMatch = true
		}
	}

	premium := models.Premium{
		ID:      transactionId,
		Status:  request.Status,
		Price:   request.Price,
		BuyerID: int(userId),
	}

	dataPremium, err := h.PremiumRepository.CreatePremium(premium)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	premiumAdded, _ := h.PremiumRepository.GetPremium(dataPremium.ID)

	// 1. Initiate Snap client
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY_PREMIUM"), midtrans.Sandbox)
	// Use to midtrans.Production if you want Production Environment (accept real transaction).

	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(dataPremium.ID),
			GrossAmt: int64(dataPremium.Price),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: dataPremium.Buyer.Username,
			Email: dataPremium.Buyer.Email,
		},
	}

	// 3. Execute request create Snap transaction to Midtrans Snap API
	snapResp, _ := s.CreateTransaction(req)

	// mengupdate token di database
	updateTokenTransaction, _ := h.PremiumRepository.UpdateTokenPremium(snapResp.Token, premiumAdded.ID)

	// mengambil data transaction yang baru diupdate
	premiumUpdated, _ := h.PremiumRepository.GetPremium(updateTokenTransaction.ID)

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertPremiumResponse(premiumUpdated)})
}

// function update premium by admin
func (h *handlerPremium) UpdatePremiumByAdmin(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	statusForm := c.FormValue("status")
	status, err := strconv.ParseBool(statusForm)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: "Invalid status value"})
	}

	request := dto.UpdatePremiumRequest{
		Status: status,
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	_, err = h.PremiumRepository.GetPremium(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResult{Status: http.StatusNotFound, Message: err.Error()})
	}

	// Update premium
	premiumUpdated, err := h.PremiumRepository.UpdatePremium(request.Status, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	// Retrieve the updated premium
	getPremiumUpdated, err := h.PremiumRepository.GetPremium(premiumUpdated.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertPremiumResponse(getPremiumUpdated)})
}

// function delete premium
func (h *handlerPremium) DeletePremium(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	premium, err := h.PremiumRepository.GetPremium(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	dataPremium, err := h.PremiumRepository.DeletePremium(premium, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: dataPremium})
}

// function notification
func (h *handlerPremium) NotificationPremium(c echo.Context) error {
	var notificationPayload map[string]interface{}

	if err := c.Bind(&notificationPayload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)
	order_id, _ := strconv.Atoi(orderId)

	premium, _ := h.PremiumRepository.GetPremium(order_id)

	fmt.Println("ini payload", notificationPayload)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			h.PremiumRepository.UpdatePremium(false, order_id)
		} else if fraudStatus == "accept" {
			h.PremiumRepository.UpdatePremium(true, order_id)
			SendMailPremium("Transaction Success", premium)
		}
	} else if transactionStatus == "settlement" {
		h.PremiumRepository.UpdatePremium(true, order_id)
		SendMailPremium("Transaction Success", premium)
	} else if transactionStatus == "deny" {
		h.PremiumRepository.UpdatePremium(false, order_id)
		SendMailPremium("Transaction Failed", premium)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		h.PremiumRepository.UpdatePremium(false, order_id)
		SendMailPremium("Transaction Failed", premium)
	} else if transactionStatus == "pending" {
		h.PremiumRepository.UpdatePremium(false, order_id)
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: notificationPayload})
}

// function SendMailPremium
func SendMailPremium(status string, premium models.Premium) {
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = "Cinema Online <rafialfian770@gmail.com>"
	var CONFIG_AUTH_EMAIL = os.Getenv("SYSTEM_EMAIL")
	var CONFIG_AUTH_PASSWORD = os.Getenv("SYSTEM_PASSWORD")

	var price = strconv.Itoa(premium.Price)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", premium.Buyer.Email)
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

	log.Println("Mail sent! to " + premium.Buyer.Email)
}

// function convert premium
func ConvertPremiumResponse(premium models.Premium) models.PremiumResponse {
	return models.PremiumResponse{
		ID:      premium.ID,
		Status:  premium.Status,
		Price:   premium.Price,
		Token:   premium.Token,
		BuyerID: premium.BuyerID,
		Buyer:   premium.Buyer,
	}
}
