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
	if _, err := tx.Exec(`INSERT INTO "order_credit_payment" (order_id, receiver_id, payment_date, cash, notes, amount) VALUES ($1, $2, $3, $4, $5, $6)`, order.ID, receiver, time.Now(), isCash, notes, moneyIn); err != nil {
		log.Warn("dbquery.instalments.go LineNew() Gagal menambahkan arah")
		tx.Rollback()
		return err
	}

	// Komit
	err = tx.Commit()
	return err
}

// InstalmentsMoneyOut uang angsuran tidak jadi masuk/dihapus
func InstalmentsMoneyOut(pid int) error {
	db := DB

	var payment wrapper.OrderPaymentSelect
	query := `SELECT *, TO_CHAR(payment_date, 'DD-MM-YYYY') AS payment_date
	FROM "order_credit_payment"
	WHERE order_id=$1`

	err := db.Get(&payment, query, pid)
	if err != nil {
		return err
	}

	var order wrapper.Order

	if ord, err := OrderGetOrderByID(payment.OrderID); err == nil {
		order = ord
	} else {
		return err
	}

	var monthlyQ []wrapper.OrderMonthlyCreditSelect

	query = `SELECT *, TO_CHAR(due_date, 'DD-MM-YYYY') AS due_date
	FROM "order_monthly_credit"
	WHERE order_id=$1 AND paid > 1 ORDER BY nth DESC`

	err = db.Select(&monthlyQ, query, payment.OrderID)
	if err != nil {
		return err
	}

	// Transaksi
	tx := db.MustBegin()
	moneyOutRem := payment.Amount
	mIndex := 0

	for moneyOutRem > 0 {
		if (len(monthlyQ) - mIndex) > 0 {
			toPull := monthlyQ[mIndex].Paid
			toCancelPay := 0
			if moneyOutRem >= toPull {
				toCancelPay = toPull
				moneyOutRem = moneyOutRem - toPull
			} else {
				toCancelPay = moneyOutRem
				moneyOutRem = 0
			}

			// Kwitansi/angsuran dianggap selesai jika "paid" sudah penuh
			isDone := true
			if toCancelPay > 0 {
				isDone = false
			}

			query = `UPDATE "order_monthly_credit"
			SET paid=$2, done=$3
			WHERE id=$1`
			_, err := tx.Exec(query, monthlyQ[mIndex].ID, monthlyQ[mIndex].Paid-toCancelPay, isDone)
			if err != nil {
				tx.Rollback()
				return err
			}

			mIndex += 1
		} else {
			//masukan remain
			fmt.Println("Lebih: ", moneyOutRem)
			moneyOutRem = 0
		}

	}

	isCreditDone := true
	isCreditActive := false
	orderStatus := "lunas"
	if (order.CreditDetail.Remaining + payment.Amount) > 0 {
		isCreditDone = false
		isCreditActive = true
		orderStatus = "aktif"
	}

	// Update kredit detail
	query = `UPDATE "order_credit_detail"
	SET last_paid=$2, remaining=$3, done=$4, active=$5
	WHERE id=$1`
	_, err = tx.Exec(query, order.CreditDetail.ID, time.Now(), order.CreditDetail.Remaining+payment.Amount, isCreditDone, isCreditActive)
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

	// Hapus histori pembayaran angsuran
	query = `DELETE FROM "order_credit_payment"
	WHERE id=$1`
	if _, err := tx.Exec(query, pid); err != nil {
		log.Warn("dbquery.instalments.go InstalmentsMoneyOut() gagal menghapus record pembayaran")
		tx.Rollback()
		return err
	}

	// Komit
	err = tx.Commit()
	return err
}
