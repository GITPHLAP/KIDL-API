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
	transactionAPI.GET("/is_balance_correct", api.IsBalanceCorrect)

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

func isTransactionsEqual(paramTrans TransactionParamConverted, latestTrans db.Transaction) bool {
	if paramTrans.Accounting.Time.Format("2006/01/02") == latestTrans.Accounting.Time.Format("2006/01/02") {
		return true
	}
	return false
}

// endregion

// The order of the transactions should be desc for accounting else it will not work
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
	// Reverse to insert it in the correct order so the latest transactions also is the latest transaction in db
	slices.Reverse(transactionsToInsert)
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
	var parameters db.TransactionsView
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	if err := c.BindJSON(&parameters); err != nil {
		log.Error("Error parsing request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

	// If categoryId is equal -1 its "no category"
	if parameters.CategoryID.Int32 == -1 {
		parameters.CategoryID.Valid = false
	}

	err := dbtx.UpdateTransactionCategoryAndAccounting(ctx, db.UpdateTransactionCategoryAndAccountingParams{
		ID:             parameters.ID,
		UserCategoryID: parameters.CategoryID,
		UserAccounting: parameters.Accounting,
	})
	if err != nil {
		log.Error("Error update transaction: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
	}

	ReturnResponseData(c, http.StatusOK, "OK", nil)
}

func (api *Api) IsBalanceCorrect(c *gin.Context) {
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	match, err := dbtx.IsBalanceMatchSumOfAmount(ctx)
	if err != nil {
		log.Error("Error checking balance: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
	}

	ReturnResponseData(c, http.StatusOK, "OK", match)
}
