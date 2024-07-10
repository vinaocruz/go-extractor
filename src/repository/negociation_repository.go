package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/vinaocruz/go-extractor/src/dto"
	"github.com/vinaocruz/go-extractor/src/model"
)

type NegociationRepository interface {
	BulkImport(toInsert []model.Negociation)
	Find(ticker, transactionAt string) (dto.Negociation, error)
	SetupIndex()
}

type PostgresNegociationRepository struct {
}

func NewNegociationRepository() *PostgresNegociationRepository {
	return &PostgresNegociationRepository{}
}

func (r *PostgresNegociationRepository) BulkImport(toInsert []model.Negociation) {
	db, err := r.getConnection()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	defer db.Close()

	txn, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to start transaction: ", err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("negotiations", "closedat", "transationat", "ticketcode", "price", "quantity"))
	if err != nil {
		log.Fatal("Failed to prepare statement: ", err)
	}

	for _, item := range toInsert {
		_, err := stmt.Exec(item.ClosedAt, item.TransactionAt, item.TicketCode, item.Price, item.Quantity)
		if err != nil {
			log.Fatal("Failed to execute statement: ", err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("Failed to execute statement: ", err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal("Failed to close statement: ", err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal("Failed to commit transaction: ", err)
	}
}

func (r *PostgresNegociationRepository) Find(ticker, transactionAt string) (dto.Negociation, error) {
	result := dto.Negociation{}

	db, err := r.getConnection()
	if err != nil {
		return result, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	query := `WITH tickers AS (
		SELECT transationat, SUM(quantity) as qntTotal, MAX(price) as maxPrice
		FROM negotiations n 
		WHERE ticketcode = $1
		%s
		GROUP BY transationat
	)
	SELECT MAX(maxprice) AS max_range_value, MAX(qntTotal) AS max_daily_volume
	FROM tickers`

	params := []interface{}{ticker}

	if transactionAt != "" {
		query = fmt.Sprintf(query, "AND transationat = $2")
		params = append(params, transactionAt)
	} else {
		query = fmt.Sprintf(query, "")
	}

	var maxRangeValue, maxDailyVolume sql.NullString

	row := db.QueryRow(query, params...)
	err = row.Scan(&maxRangeValue, &maxDailyVolume)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, fmt.Errorf("ticker not found")
		}
		return result, fmt.Errorf("failed to execute query: %v", err)
	}

	if !maxDailyVolume.Valid || !maxRangeValue.Valid {
		return result, fmt.Errorf("not fount")
	}

	if maxRangeValue.Valid {
		var tmp float64
		tmp, err = strconv.ParseFloat(maxRangeValue.String, 32)
		if err != nil {
			return result, fmt.Errorf("failed to parse max_range_value: %v", err)
		}
		result.Max_range_value = float32(tmp)
	}

	if maxDailyVolume.Valid {
		result.Max_daily_volume, err = strconv.Atoi(maxDailyVolume.String)
		if err != nil {
			return result, fmt.Errorf("failed to parse max_daily_volume: %v", err)
		}
	}

	return result, nil
}

func (r *PostgresNegociationRepository) SetupIndex() {
	db, err := r.getConnection()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE INDEX negotiations_transationat_idx ON public.negotiations (transationat);
	CREATE INDEX negotiations_ticketcode_idx ON public.negotiations (ticketcode);
	`)
	if err != nil {
		log.Fatal("Failed to create index: ", err)
	}
}

func (r *PostgresNegociationRepository) getConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST")),
	)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		fmt.Println("Failed to ping")
		return nil, err
	}

	return db, nil
}
