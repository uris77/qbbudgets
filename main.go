package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"qbbudgets/budgets"
)

func main() {
	cnf := budgets.ReadConf("prod")
	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", cnf.DbUser, cnf.DbName, cnf.DbPassword, cnf.DbHost)

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	ticket := budgets.Auth(cnf.QbUsername, cnf.QbPassword)

	cnt := budgets.CountQbBudgets(ticket.Ticket, "bg44v3kd9")
	sz := 100
	pages := int(math.Ceil(float64(cnt.Total) / float64(sz)))

	log.Printf("Number of Pages: %d", pages)
	i := 0
	var budgetDoc budgets.BudgetDoc
	for i < pages {
		budgetDoc = budgets.GetBudgets(ticket.Ticket, "bg44v3kd9", sz, i*sz)
		budgets.Upsert(db, budgetDoc)
		log.Printf("Iteration: %d", i)
		i++
	}

}
