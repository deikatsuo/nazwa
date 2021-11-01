package dbquery

import (
	"fmt"
	"nazwa/wrapper"
	"time"
)

// InstalmentsReceiptPrintedStatus update nama arah
func InstalmentsReceiptPrintedStatus(rid int, printed bool) error {
	db := DB
	query := `UPDATE "order_monthly_credit"
	SET printed=$2
	WHERE id=$1`
	_, err := db.Exec(query, rid, printed)

	return err
}

// InstalmentsReceiptUpdateNotes update notes
func InstalmentsReceiptUpdateNotes(rid int, notes string) error {
	db := DB
	query := `UPDATE "order_monthly_credit"
	SET notes=$2
	WHERE id=$1`
	_, err := db.Exec(query, rid, notes)

	return err
}

// InstalmentsMoneyIn uang angsuran masuk
func InstalmentsMoneyIn(oid, receiver, moneyIn int, notes, mode string) error {
	db := DB

	var order wrapper.Order

	if ord, err := OrderGetOrderByID(oid); err == nil {
		order = ord
	} else {
		return err
	}

	var monthlyQ []wrapper.OrderMonthlyCreditSelect

	query := `SELECT *, TO_CHAR(due_date, 'DD-MM-YYYY') AS due_date
	FROM "order_monthly_credit"
	WHERE order_id=$1 AND done=false ORDER BY nth`

	err := db.Select(&monthlyQ, query, oid)
	if err != nil {
		return err
	}

	// Transaksi
	tx := db.MustBegin()
	moneyInRem := moneyIn
	mIndex := 0

	for moneyInRem > 0 {
		if (len(monthlyQ) - mIndex) > 0 {
			toFill := order.CreditDetail.Monthly - monthlyQ[mIndex].Paid
			toPay := 0
			if moneyInRem >= toFill {
				toPay = toFill
				moneyInRem = moneyInRem - toFill
			} else {
				toPay = moneyInRem
				moneyInRem = 0
			}

			// Kwitansi/angsuran dianggap selesai jika "paid" sudah penuh
			isDone := false
			if (monthlyQ[mIndex].Paid + toPay) == order.CreditDetail.Monthly {
				isDone = true
			}

			query = `UPDATE "order_monthly_credit"
			SET paid=$2, done=$3
			WHERE id=$1`
			_, err := tx.Exec(query, monthlyQ[mIndex].ID, monthlyQ[mIndex].Paid+toPay, isDone)
			if err != nil {
				tx.Rollback()
				return err
			}

			mIndex += 1
		} else {
			//masukan remain
			fmt.Println("Lebih: ", moneyInRem)
			moneyInRem = 0
		}

	}

	isCreditDone := false
	isCreditActive := true
	orderStatus := "aktif"
	if (order.CreditDetail.Remaining - moneyIn) <= 0 {
		isCreditDone = true
		isCreditActive = false
		orderStatus = "lunas"
	}

	// Update kredit detail
	query = `UPDATE "order_credit_detail"
	SET last_paid=$2, remaining=$3, done=$4, active=$5
	WHERE id=$1`
	_, err = tx.Exec(query, order.CreditDetail.ID, time.Now(), order.CreditDetail.Remaining-moneyIn, isCreditDone, isCreditActive)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update order status
	query = `UPDATE "order"
	SET status=$2
	WHERE id=$1`
	_, err = tx.Exec(query, order.ID, orderStatus)
	if err != nil {
		tx.Rollback()
		return err
	}

	isCash := true
	if mode == "transfer" {
		isCash = false
	}
	// Tambahkan histori pembayaran angsuran
	if _, err := db.Exec(`INSERT INTO "order_credit_payment" (order_id, receiver_id, payment_date, cash, notes, amount) VALUES ($1, $2, $3, $4, $5, $6)`, order.ID, receiver, time.Now(), isCash, notes, moneyIn); err != nil {
		log.Warn("dbquery.line.go LineNew() Gagal menambahkan arah")
		tx.Rollback()
		return err
	}

	// Komit
	err = tx.Commit()
	return err
}
