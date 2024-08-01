package api

import (
	"context"
	"kidl-api/pkg/db"
	"net/http"
	"slices"
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
	transactionAPI.GET("/:year", api.getTransactions)
	transactionAPI.PATCH("", api.updateTransaction)

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
	Accounting  pgtype.Date   `json:"Bogføringsdato"`
	Amount      pgtype.Float8 `json:"Beløb"`
	Sender      pgtype.Int8   `json:"Afsender"`
	Recipient   pgtype.Int8   `json:"Modtager"`
	Name        pgtype.Text   `json:"Navn"`
	Description pgtype.Text   `json:"Beskrivelse"`
	Balance     pgtype.Float8 `json:"Saldo"`
	Currency    pgtype.Text   `json:"Valuta"`
	Aligned     pgtype.Text   `json:"Afstemt"`
}

// region helperFunctions
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

		// Try parse accounting
		accounting, err := time.Parse("2006/01/02", t.Accounting)
		if err != nil {
			log.Info(t.Accounting, " is not a valid date!")
			continue
		}
		// Assign all the values to object and append to the result array
		result = append(result, TransactionParamConverted{
			Accounting:  pgtype.Date{Time: accounting, Valid: true},
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

func timeSortFuncDesc(a, b TransactionParamConverted) int {
	if a.Accounting.Time.Equal(b.Accounting.Time) {
		return 0
	} else if a.Accounting.Time.Before(b.Accounting.Time) {
		return 1
	}
	return -1
}

func timeSortFuncAsc(a, b db.InsertNewTransactionsParams) int {
	if a.Accounting.Time.Equal(b.Accounting.Time) {
		return 0
	} else if a.Accounting.Time.Before(b.Accounting.Time) {
		return -1
	}
	return 1
}

func isTransactionsEqual(paramTrans TransactionParamConverted, latestTrans db.TransactionsTest) bool {
	if paramTrans.Accounting.Time.Format("2006/01/02") == latestTrans.Accounting.Time.Format("2006/01/02") &&
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

// endregion

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

	slices.SortFunc(parameters, timeSortFuncDesc)

	for _, t := range parameters {
		// Check for full match if match break because that transactin is allready added
		if isTransactionsEqual(t, latestTrans) {
			goto InsertPart
		}

		// Add transaction to array of transactions to insert
		transactionsToInsert = append(transactionsToInsert, db.InsertNewTransactionsParams{
			Accounting:  t.Accounting,
			Amount:      t.Amount,
			Sender:      t.Sender,
			Recipient:   t.Recipient,
			Name:        t.Name,
			Description: t.Description,
			Balance:     t.Balance,
			Currency:    t.Currency,
			Aligned:     t.Aligned,
			Reserved:    pgtype.Bool{Bool: false, Valid: false},
		})
	}

InsertPart:
	slices.SortFunc(transactionsToInsert, timeSortFuncAsc)
	addTransactions, err := dbtx.InsertNewTransactions(ctx, transactionsToInsert)
	if err != nil {
		log.Error("Error inserting transactions: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
	}

	ReturnResponseData(c, http.StatusOK, "OK", addTransactions)
}

func (api *Api) getTransactions(c *gin.Context) {
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	yearstr := c.Param("year")

	// parse year to int
	year, err := strconv.Atoi(yearstr)
	if err != nil {
		log.Error("Cant convert year to string")
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing year")
	}

	transactions, err := dbtx.GetTransactionsFromViewByYear(ctx, int32(year))
	if err != nil {
		log.Error("Error fetching transactions by year: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
	}

	ReturnResponseData(c, http.StatusOK, "OK", transactions)
}

func (api *Api) updateTransaction(c *gin.Context) {
	var parameters db.TransactionsTest
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	if err := c.BindJSON(&parameters); err != nil {
		log.Error("Error parsing request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

	err := dbtx.UpdateTransactionCategoryAndAccounting(ctx, db.UpdateTransactionCategoryAndAccountingParams{
		ID:             parameters.ID,
		UserCategory:   parameters.UserCategory,
		UserAccounting: parameters.UserAccounting,
	})
	if err != nil {
		log.Error("Error update transaction: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
	}

	ReturnResponseData(c, http.StatusOK, "OK", nil)
}
