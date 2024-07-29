package api

import (
	"context"
	"kidl-api/pkg/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

func (api *Api) transactionsSubrouter(r *gin.RouterGroup) {
	transactionAPI := r.Group("/transaction")

	transactionAPI.POST("", api.uploadTransactions)

}

type TransactionParam struct {
	Accounting  string        `json:"Bogføringsdato"`
	Amount      pgtype.Float8 `json:"Beløb"`
	Sender      pgtype.Int8   `json:"Afsender"`
	Recipient   pgtype.Int8   `json:"Modtager"`
	Name        pgtype.Text   `json:"Navn"`
	Description pgtype.Text   `json:"Beskrivelse"`
	Balance     pgtype.Float8 `json:"Saldo"`
	Currency    pgtype.Text   `json:"Valuta"`
	Aligned     pgtype.Text   `json:"Afstemt"`
}

func (api *Api) uploadTransactions(c *gin.Context) {
	var parameters []TransactionParam

	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	if err := c.BindJSON(&parameters); err != nil {
		log.Error("Error parsing request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

	// Get lates transaction
	latestTrans, err := dbtx.GetLatestTransaction(ctx)
	if err != nil {
		log.Error("Error getting latest transactions: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Internal server error")
		return
	}
	var transactionsToInsert []db.InsertNewTransactionsParams

	for _, t := range parameters {
		accounting, err := time.Parse("2006/01/02", t.Accounting)
		if err != nil {
			log.Info(t.Accounting, " is not a valid date!")
			continue
		}

		// Check for full match if match break because that transactin is allready added
		if isTransactionsEqual(t, latestTrans) {
			goto InsertPart
		}

		// Add transaction to array of transactions to insert
		transactionsToInsert = append(transactionsToInsert, db.InsertNewTransactionsParams{
			Accounting:  pgtype.Date{Time: accounting, Valid: true},
			Amount:      t.Amount,
			Sender:      t.Sender,
			Recipient:   t.Sender,
			Name:        t.Name,
			Description: t.Description,
			Balance:     t.Balance,
			Currency:    t.Currency,
			Aligned:     t.Aligned,
			Reserved:    pgtype.Bool{Bool: false, Valid: false},
		})
	}

InsertPart:
	_, err = dbtx.InsertNewTransactions(ctx, transactionsToInsert)
	if err != nil {
		log.Error("Error inserting transactions: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Internal server error")
		return
	}
}

func isTransactionsEqual(paramTrans TransactionParam, latestTrans db.Transaction) bool {
	accounting, err := time.Parse("2006/01/02", paramTrans.Accounting)
	if err != nil {
		return false
	}

	if accounting == latestTrans.Accounting.Time &&
		paramTrans.Currency.String == latestTrans.Currency.String &&
		paramTrans.Description.String == latestTrans.Description.String &&
		paramTrans.Name.String == latestTrans.Name.String &&
		paramTrans.Aligned.String == latestTrans.Aligned.String &&
		paramTrans.Amount.Float64 == latestTrans.Amount.Float64 &&
		paramTrans.Balance.Float64 == latestTrans.Balance.Float64 &&
		paramTrans.Currency.String == latestTrans.Currency.String &&
		paramTrans.Recipient.Int64 == latestTrans.Recipient.Int64 &&
		paramTrans.Sender.Int64 == latestTrans.Sender.Int64 {

		return true
	}
	return false
}
