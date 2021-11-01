package dbquery

import (
	"fmt"
	"nazwa/wrapper"
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
func InstalmentsMoneyIn(oid int, moneyIn int) error {
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
			fmt.Println("Topay: ", toPay)

			mIndex += 1
		} else {
			//todo
			//masukan remain
			fmt.Println("Lebih: ", moneyInRem)
			moneyInRem = 0
		}
	}

	return nil
}
