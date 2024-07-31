package api

import (
	"context"
	"kidl-api/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

func (api *Api) transactionsSubrouter(r *gin.RouterGroup) {
	transactionAPI := r.Group("/transaction")

	transactionAPI.POST("", api.uploadTransactions)

}

type TransactionParamRaw struct {
	Accounting  string `json:"Bogføringsdato"`
	Amount      string `json:"Beløb"`
	Sender      string `json:"Afsender"`
	Recipient   string `json:"Modtager"`
	Name        string `json:"Navn"`
	Description string `json:"Beskrivelse"`
	Balance     string `json:"Saldo"`
	Currency    string `json:"Valuta"`
	Aligned     string `json:"Afstemt"`
}

type TransactionParamConverted struct {
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

func convertTransactionParameters(parameters []TransactionParamRaw) ([]TransactionParamConverted, error) {
	result := []TransactionParamConverted{}

	for _, t := range parameters {
		// replace , with . for numbers
		dotAmount := strings.Replace(t.Amount, ",", ".", -1)

		dotBalance := strings.Replace(t.Balance, ",", ".", -1)

		// Parse float values
		amount, err := strconv.ParseFloat(dotAmount, 64)
		if err != nil {
			return nil, err
		}
		balance, err := strconv.ParseFloat(dotBalance, 64)
		if err != nil {
			return nil, err
		}

		// Try parse int values this is allowed to be null
		sender, sendererr := strconv.ParseInt(t.Sender, 10, 64)
		recipient, recipienterr := strconv.ParseInt(t.Recipient, 10, 64)

		// Assign all the values to object and append to the result array
		result = append(result, TransactionParamConverted{
			Accounting:  t.Accounting,
			Amount:      pgtype.Float8{Float64: amount, Valid: true},
			Sender:      pgtype.Int8{Int64: sender, Valid: sendererr == nil},
			Recipient:   pgtype.Int8{Int64: recipient, Valid: recipienterr == nil},
			Name:        pgtype.Text{String: t.Name, Valid: true},
			Description: pgtype.Text{String: t.Description, Valid: true},
			Balance:     pgtype.Float8{Float64: balance, Valid: true},
			Currency:    pgtype.Text{String: t.Currency, Valid: true},
			Aligned:     pgtype.Text{String: t.Aligned, Valid: true},
		})
	}
	return result, nil
}
func (api *Api) uploadTransactions(c *gin.Context) {
	var parametersraw []TransactionParamRaw
	var parameters []TransactionParamConverted

	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	if err := c.BindJSON(&parametersraw); err != nil {
		log.Error("Error parsing request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

	parameters, err := convertTransactionParameters(parametersraw)
	if err != nil {
		log.Error("Error converting request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

	// Get latest transaction
	latestTrans, err := dbtx.GetLatestTransaction(ctx)
	if err != nil {
		log.Error("Error getting latest transactions: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
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
	addTransactions, err := dbtx.InsertNewTransactions(ctx, transactionsToInsert)
	if err != nil {
		log.Error("Error inserting transactions: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
	}

	ReturnResponseData(c, http.StatusOK, "OK", addTransactions)
}

func isTransactionsEqual(paramTrans TransactionParamConverted, latestTrans db.TransactionsTest) bool {
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
