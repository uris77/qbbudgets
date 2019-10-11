package budgets

import (
	"database/sql"
	"fmt"
	"math"
)

func byId(db *sql.DB, id int64) bool {
	q := fmt.Sprintf("select id from quickbase_budgets where budget_id = %d", id)
	rows, err := db.Query(q)
	if err != nil {
		Logger.Fatalw("Error happened while querying quickbase budgets", "id", id, "error", err)
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			Logger.Fatalw("Error closing Query", "error", err, "query", q)
		}
	}()

	exists := rows.Next()
	Logger.Debugw("Budget exists?", "id", id, "exists", exists)
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
				Logger.Fatalw("Preparing upsert statement failed", "error", err, "statement", updateStmt)
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
				Logger.Fatalw("Executing upsert statement failed", "error", err, "statement", updateStmt)
			}

			rowCnt, err := res.RowsAffected()
			if err != nil {
				Logger.Fatalw("Retrieve RowsAffected Failed", "error", err, "statement", updateStmt)
			}
			Logger.Debugw("Updated Budget", "rowCount", rowCnt, "budgetId", budget.BudgetId, "budgetMonth", budget.BudgetMonth(), "month", budget.Month)
		} else {
			//Insert
			insertStmt, err := db.Prepare(insertStmt)
			if err != nil {
				Logger.Fatalw("Error Inserting Budget", "error", err)
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
				Logger.Fatalw("Inserting Statment Failed", "error", err, "insertStatement", insertStmt)
			}

			rowCnt, err := res.RowsAffected()
			if err != nil {
				Logger.Fatalw("Failed to retrieve RowsAffected", "error", err, "statement", insertStmt)
			}
			Logger.Debugw("New Budget affected", "rowCnt", rowCnt)
		}
	}
}
