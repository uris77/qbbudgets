package main

import (
	"database/sql"
	"fmt"
	"math"

	"qbbudgets/budgets"
)

func main() {
	budgets.SetupLogger()
	defer budgets.Logger.Sync()
	cnf := budgets.ReadConf("prod")
	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", cnf.DbUser, cnf.DbName, cnf.DbPassword, cnf.DbHost)

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		budgets.Logger.Fatalw("Failed to open database connection", "error", err)
	}

	defer db.Close()
	ticket := budgets.Auth(cnf.QbUsername, cnf.QbPassword)

	cnt := budgets.CountQbBudgets(ticket.Ticket, "bg44v3kd9")
	sz := 100
	pages := int(math.Ceil(float64(cnt.Total) / float64(sz)))

	budgets.Logger.Debugw("Number of Pages", "pages", pages)
	i := 0
	var budgetDoc budgets.BudgetDoc
	for i < pages {
		budgetDoc = budgets.GetBudgets(ticket.Ticket, "bg44v3kd9", sz, i*sz)
		budgets.Upsert(db, budgetDoc)
		budgets.Logger.Debugw("Iteration", "iteration", i)
		i++
	}

}
