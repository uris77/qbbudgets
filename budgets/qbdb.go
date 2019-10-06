package budgets

import (
	"database/sql"
	"fmt"
	"log"
	"math"
)

func byId(db *sql.DB, id int64) bool {
	q := fmt.Sprintf("select id from quickbase_budgets where budget_id = %d", id)
	rows, err := db.Query(q)
	if err != nil {
		log.Fatalf("Error happened while querying quickbase budgets by id(%d). Error: %v", id, err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatalf("Error closing Query %v", err)
		}
	}()

	exists := rows.Next()
	log.Printf("Budget with id %d exists? %v", id, exists)
	return exists
}

func Upsert(db *sql.DB, budget BudgetDoc) {
	insertStmt := `INSERT INTO quickbase_budgets (budget_id, partner, client_name, project_name, initiative,
month, billable_fees, client_budget, external_hourly_rate, dbname, initiative_billable, component,
project_status, date_modified, product, project_id, related_client, budget_hours, billing_month, minimum,
prebill_media, prebill_management, bfo_pct, referral_fee, project_prebill_tech, project_mgmt_fee_calc)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26)`

	updateStmt := `UPDATE quickbase_budgets SET (partner, client_name, project_name, initiative,
month, billable_fees, client_budget, external_hourly_rate, dbname, initiative_billable, component,
project_status, date_modified, product, project_id, related_client, budget_hours, billing_month, minimum,
prebill_media, prebill_management, bfo_pct, referral_fee, project_prebill_tech, project_mgmt_fee_calc) = 
($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25) WHERE budget_id = $26`
	budgets := budget.Budgets
	for i := 0; i < len(budgets); i++ {
		budget := budgets[i]

		if byId(db, budget.BudgetId) == true {
			//Update
			pstmt, err := db.Prepare(updateStmt)
			if err != nil {
				log.Fatal(err)
			}
			var budgetHours float64 = 0
			if budget.ExternalHourlyHourRate > 0 {
				budgetHours = math.Round(budget.BillableFees / budget.ExternalHourlyHourRate)
			}


			res, err := pstmt.Exec(budget.Partner, budget.ClientName, budget.ProjectName,
				budget.Initiative, budget.BudgetMonth(), budget.BillableFees, budget.ClientBudget,
				budget.ExternalHourlyHourRate, budget.Dbname, budget.InitiativeBillable, budget.Component,
				budget.ProjectStatus, budget.DateModified(), budget.Product, budget.ProjectId,
				budget.RelatedClient, budgetHours, budget.BillingMonth(), budget.Minimum,
				budget.PrebillMedia, budget.PrebillManagement, budget.BfoPct, budget.ReferralFee,
				budget.ProjectPrebillTech, budget.ProjectManagementFeeCalc, budget.BudgetId)

			if err != nil {
				log.Fatal(err)
			}

			rowCnt, err := res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Updated Budget affected = %d with id: %d and month: %v and month as epoch: %s\n", rowCnt, budget.BudgetId,  budget.BudgetMonth(), budget.Month)
		} else {
			//Insert
			insertStmt, err := db.Prepare(insertStmt)
			if err != nil {
				log.Fatal("Error Inserting Budget: ", err)
			}
			var budgetHours float64 = 0
			if budget.ExternalHourlyHourRate > 0 {
				budgetHours = math.Round(budget.BillableFees / budget.ExternalHourlyHourRate)
			}

			res, err := insertStmt.Exec(budget.BudgetId, budget.Partner, budget.ClientName, budget.ProjectName,
				budget.Initiative, budget.BudgetMonth(), budget.BillableFees, budget.ClientBudget,
				budgetHours, budget.Dbname, budget.InitiativeBillable, budget.Component,
				budget.ProjectStatus, budget.DateModified(), budget.Product, budget.ProjectId,
				budget.RelatedClient, budgetHours, budget.BillingMonth(), budget.Minimum,
				budget.PrebillMedia, budget.PrebillManagement, budget.BfoPct, budget.ReferralFee,
				budget.ProjectPrebillTech, budget.ProjectManagementFeeCalc)
			if err != nil {
				log.Fatal(err)
			}

			//lastId, err := res.LastInsertId()
			//if err != nil {
			//	log.Fatal(err)
			//}
			rowCnt, err := res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("New Budget affected = %d\n", rowCnt)
		}
	}
}
