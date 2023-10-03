package handlers

import (
	"cinemaonline/dto"
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

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

// function get all transactions by user
func (h *handlerTransaction) FindTransactionsByUser(c echo.Context) error {
	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	transactions, err := h.TransactionRepository.FindTransactionsByUser(int(userId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	for i, transaction := range transactions {
		transactions[i].Movie.Thumbnail = path_thumbnail + transaction.Movie.Thumbnail
		transactions[i].Movie.Trailer = path_trailer + transaction.Movie.Trailer
		transactions[i].Movie.FullMovie = path_full_movie + transaction.Movie.FullMovie
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertMultipleTransactionResponse(transactions)})
}

// function get all transactions
func (h *handlerTransaction) FindTransactions(c echo.Context) error {
	transactions, err := h.TransactionRepository.FindTransactions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	for i, transaction := range transactions {
		transactions[i].Movie.Thumbnail = path_thumbnail + transaction.Movie.Thumbnail
		transactions[i].Movie.Trailer = path_trailer + transaction.Movie.Trailer
		transactions[i].Movie.FullMovie = path_full_movie + transaction.Movie.FullMovie
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertMultipleTransactionResponse(transactions)})
}

// function get transaction by id
func (h *handlerTransaction) GetTransaction(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var transaction models.Transaction
	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	transaction.Movie.Thumbnail = path_thumbnail + transaction.Movie.Thumbnail
	transaction.Movie.Trailer = path_trailer + transaction.Movie.Trailer
	transaction.Movie.FullMovie = path_full_movie + transaction.Movie.FullMovie

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: transaction})
}

// function create transaction
func (h *handlerTransaction) CreateTransaction(c echo.Context) error {
	request := new(dto.CreateTransactionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	movie, _ := h.TransactionRepository.GetMovie(request.MovieID)

	request.BuyerID = int(userId)
	request.SellerID = movie.UserID
	request.Status = "pending"

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	var transactionIsMatch = false
	var transactionId int
	for !transactionIsMatch {
		transactionId = int(time.Now().Unix())
		transactionData, _ := h.TransactionRepository.GetTransaction(transactionId)
		if transactionData.ID == 0 {
			transactionIsMatch = true
		}
	}

	transaction := models.Transaction{
		ID:       transactionId,
		MovieID:  request.MovieID,
		BuyerID:  int(userId),
		SellerID: movie.UserID,
		Price:    request.Price,
		Status:   request.Status,
	}

	dataTransactions, err := h.TransactionRepository.CreateTransaction(transaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	transactionAdded, _ := h.TransactionRepository.GetTransaction(dataTransactions.ID)

	// 1. Initiate Snap client
	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY_TRANSACTION_MOVIE"), midtrans.Sandbox)
	// Use to midtrans.Production if you want Production Environment (accept real transaction).

	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(dataTransactions.ID),
			GrossAmt: int64(dataTransactions.Price),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: dataTransactions.Buyer.Username,
			Email: dataTransactions.Buyer.Email,
		},
	}

	// 3. Execute request create Snap transaction to Midtrans Snap API
	snapResp, _ := s.CreateTransaction(req)

	// mengupdate token di database
	updateTokenTransaction, _ := h.TransactionRepository.UpdateTokenTransaction(snapResp.Token, transactionAdded.ID)

	// mengambil data transaction yang baru diupdate
	transactionUpdated, _ := h.TransactionRepository.GetTransaction(updateTokenTransaction.ID)

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertTransactionResponse(transactionUpdated)})
}

// function update transaction by admin
func (h *handlerTransaction) UpdateTransactionByAdmin(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	status := c.FormValue("status")
	request := dto.UpdateTransactionRequest{
		Status: status,
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	_, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResult{Status: http.StatusNotFound, Message: err.Error()})
	}

	// Update the transaction
	transactionUpdated, err := h.TransactionRepository.UpdateTransaction(request.Status, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	// Retrieve the updated transaction
	getTransactionUpdated, err := h.TransactionRepository.GetTransaction(transactionUpdated.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertTransactionResponse(getTransactionUpdated)})
}

// function delete transaction
func (h *handlerTransaction) DeleteTransaction(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.TransactionRepository.DeleteTransaction(transaction, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertTransactionResponse(data)})
}

// function notification
func (h *handlerTransaction) Notification(c echo.Context) error {
	var notificationPayload map[string]interface{}

	if err := c.Bind(&notificationPayload); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)
	order_id, _ := strconv.Atoi(orderId)

	// fmt.Print("ini payloadnya", notificationPayload)
	// fmt.Println("order id", order_id)

	transaction, _ := h.TransactionRepository.GetTransaction(order_id)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			h.TransactionRepository.UpdateTransaction("pending", order_id)
		} else if fraudStatus == "accept" {
			h.TransactionRepository.UpdateTransaction("success", order_id)
			SendMailTransactionMovie("Transaction Success", transaction)
		}
	} else if transactionStatus == "settlement" {
		h.TransactionRepository.UpdateTransaction("success", order_id)
		SendMailTransactionMovie("Transaction Success", transaction)
	} else if transactionStatus == "deny" {
		h.TransactionRepository.UpdateTransaction("failed", order_id)
		SendMailTransactionMovie("Transaction Failed", transaction)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		h.TransactionRepository.UpdateTransaction("failed", order_id)
		SendMailTransactionMovie("Transaction Failed", transaction)
	} else if transactionStatus == "pending" {
		h.TransactionRepository.UpdateTransaction("pending", order_id)
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: notificationPayload})
}

// function sendmail
func SendMailTransactionMovie(status string, transaction models.Transaction) {
	// if status != transaction.Status && (status == "success") {}
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = "Cinema Online <rafialfian770@gmail.com>"
	var CONFIG_AUTH_EMAIL = os.Getenv("SYSTEM_EMAIL")
	var CONFIG_AUTH_PASSWORD = os.Getenv("SYSTEM_PASSWORD")

	var movieTitle = transaction.Movie.Title
	var price = strconv.Itoa(transaction.Price)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", transaction.Buyer.Email)
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
				<h2>Movie payment :</h2>
				<ul style="list-style-type:none;">
					<li>Title : %s</li>
					<li>Total payment: Rp.%s</li>
					<li>Status : <b>%s</b></li>
				</ul>
			</body>
		</html>`, movieTitle, price, status))

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

	log.Println("Mail sent! to " + transaction.Buyer.Email)
}

// function convert transaction
func ConvertTransactionResponse(transaction models.Transaction) models.TransactionResponse {
	return models.TransactionResponse{
		ID:       transaction.ID,
		BuyerID:  transaction.BuyerID,
		Buyer:    transaction.Buyer,
		SellerID: transaction.SellerID,
		Seller:   transaction.Seller,
		Price:    transaction.Price,
		Status:   transaction.Status,
		Token:    transaction.Token,
		MovieID:  transaction.MovieID,
		Movie: models.MovieResponse{
			ID:          transaction.Movie.ID,
			Title:       transaction.Movie.Title,
			CategoryID:  transaction.Movie.CategoryID,
			ReleaseDate: transaction.Movie.ReleaseDate,
			Price:       transaction.Movie.Price,
			Link:        transaction.Movie.Link,
			Description: transaction.Movie.Description,
			Thumbnail:   transaction.Movie.Thumbnail,
			Trailer:     transaction.Movie.Trailer,
			FullMovie:   transaction.Movie.FullMovie,
			UserID:      transaction.Movie.UserID,
			User: models.UserResponse{
				ID:       transaction.Movie.User.ID,
				Username: transaction.Movie.User.Username,
				Email:    transaction.Movie.User.Email,
				Password: transaction.Movie.User.Password,
				Gender:   transaction.Movie.User.Gender,
				Phone:    transaction.Movie.User.Phone,
				Address:  transaction.Movie.User.Address,
				Photo:    transaction.Movie.User.Photo,
				Premi:    transaction.Movie.User.Premi,
			},
		},
	}
}

// function convert multiple transaction
func ConvertMultipleTransactionResponse(transaction []models.Transaction) []models.TransactionResponse {
	var result []models.TransactionResponse

	for _, trans := range transaction {
		var ratings []models.RatingResponse
		for _, rating := range trans.Movie.Rating {
			r := models.RatingResponse{
				ID:      rating.ID,
				Star:    rating.Star,
				MovieID: rating.MovieID,
				UserID:  rating.UserID,
				User:    rating.User,
			}
			ratings = append(ratings, r)
		}

		transaction := models.TransactionResponse{
			ID:       trans.ID,
			BuyerID:  trans.BuyerID,
			Buyer:    trans.Buyer,
			SellerID: trans.SellerID,
			Seller:   trans.Seller,
			Price:    trans.Price,
			Status:   trans.Status,
			Token:    trans.Token,
			MovieID:  trans.MovieID,
			Movie: models.MovieResponse{
				ID:          trans.Movie.ID,
				Title:       trans.Movie.Title,
				CategoryID:  trans.Movie.CategoryID,
				ReleaseDate: trans.Movie.ReleaseDate,
				Price:       trans.Movie.Price,
				Link:        trans.Movie.Link,
				Description: trans.Movie.Description,
				Thumbnail:   trans.Movie.Thumbnail,
				Trailer:     trans.Movie.Trailer,
				FullMovie:   trans.Movie.FullMovie,
				UserID:      trans.Movie.UserID,
				User: models.UserResponse{
					ID:       trans.Movie.User.ID,
					Username: trans.Movie.User.Username,
					Email:    trans.Movie.User.Email,
					Password: trans.Movie.User.Password,
					Gender:   trans.Movie.User.Gender,
					Phone:    trans.Movie.User.Phone,
					Address:  trans.Movie.User.Address,
					Photo:    trans.Movie.User.Photo,
					Premi:    trans.Movie.User.Premi,
				},
				Rating: ratings,
			},
		}

		for _, cat := range trans.Movie.Category {
			category := models.CategoryResponse{
				ID:   cat.ID,
				Name: cat.Name,
			}
			transaction.Movie.Category = append(transaction.Movie.Category, category)
		}

		result = append(result, transaction)
	}

	return result
}
