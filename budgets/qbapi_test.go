package budgets

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func init() {
	sugar = zap.NewExample().Sugar()
	defer sugar.Sync()
}

func TestUnmarshalBudget(t *testing.T) {
	sample := `<?xml version="1.0" ?>
<qdbapi>
	<action>API_DoQuery</action>
	<errcode>0</errcode>
	<errtext>No error</errtext>
<dbinfo>
<name>Budgets</name>
<desc></desc>
</dbinfo>
<variables>
</variables>
<chdbids>
</chdbids>
  <record rid="10345">
    <project_name>Motorola- GAP - 8/31/2019</project_name>
    <project___contract_status>Signed</project___contract_status>
    <related_project>1068</related_project>
    <project___client_billable_>1</project___client_billable_>
    <billing_month>1501545600000</billing_month>
    <month>1504224000000</month>
    <component>sept, oct, nov</component>
    <initiative>Motorola - Analytics 360 Management</initiative>
    <comment/>
    <client_budget>1000.10</client_budget>
    <bfo__/>
    <___referral_fee/>
    <minimum/>
    <billable_fees>1000.00</billable_fees>
    <external_hourly_rate>195</external_hourly_rate>
    <budget_hours_by_external_rate>5.13</budget_hours_by_external_rate>
    <budget_hour_override/>
    <adjusted_budget_hours>5.13</adjusted_budget_hours>
    <monthly_budgets>10345</monthly_budgets>
    <add_monthly_budget>https://befoundonline.quickbase.com/db/bj22i7dry?a=API_GenAddRecordForm&amp;_fid_56=10345&amp;z=</add_monthly_budget>
    <dbname>Project Central</dbname>
    <client_name/>
    <date_modified>1557773980205</date_modified>
    <initiative___billable>1</initiative___billable>
    <project___client_name>Motorola</project___client_name>
    <project___partner>BFO Direct</project___partner>
    <project_status>Active</project_status>
    <record_id_>10345</record_id_>
    <___actual_costs/>
    <___cogs>0.00</___cogs>
    <___est__costs>0.00</___est__costs>
    <___labor_budget>1000.00</___labor_budget>
    <___labor_budget__ccbm_>0.00</___labor_budget__ccbm_>
    <___overhead>0.00</___overhead>
    <comments2/>
    <___partner_fees/>
    <___target_profit_margin>0.00</___target_profit_margin>
    <___est__costs/>
    <___overhead/>
    <___target_profit_margin/>
    <actual_spend__ccbm_/>
    <actual_spend_fees___total/>
    <add_rollover_to_specified_month__ccbm_>1509494400000</add_rollover_to_specified_month__ccbm_>
    <billable_fees_based_on_actual_spend__ccbm_/>
    <budget_display_name>Motorola- GAP - 8/31/2019 - 09-01-2017 - Motorola - Analytics 360 Management</budget_display_name>
    <budgeted_hours>15.38</budgeted_hours>
    <client___partner_name/>
    <client___vertical/>
    <date_created>1502985937178</date_created>
    <initiative___project___hours_spent_last_month/>
    <initiative___total_hours/>
    <initiative___total_hours_spent/>
    <initiative_name>Analytics 360 Management</initiative_name>
    <initiative_project_name>Motorola- GAP - 8/31/2019</initiative_project_name>
    <last_modified_by>beth.spiegel@befoundonline.com</last_modified_by>
    <month_eom>1506729600000</month_eom>
    <month_fom>1504224000000</month_fom>
    <non_labor_expenses>0.00</non_labor_expenses>
    <operating_budget___do_not_use>1000.00</operating_budget___do_not_use>
    <project____>&lt;div&gt;&lt;img src=&quot;https://images.quickbase.com/si/16/240-triang_green.png&quot;&gt;&lt;/div&gt;</project____>
    <project___account_directors>steve.kozma@befoundonline.com</project___account_directors>
    <project___analysts/>
    <project___billable_fees_ytd>25000.00</project___billable_fees_ytd>
    <project___client_full_name>Motorola Solutions</project___client_full_name>
    <project___managers/>
    <project___ops_level_access>accounting@befoundonline.com;beth.spiegel@befoundonline.com;dan.golden@befoundonline.com;julia.ebner@befoundonline.com;steve.krull@befoundonline.com</project___ops_level_access>
    <project___product>GAC</project___product>
    <project___reason_closed/>
    <project___related_client>212</project___related_client>
    <project___staffing_cost_total/>
    <project___staffing_costs_this_month/>
    <project___start_date>1504224000000</project___start_date>
    <project___start_date_fom>1504224000000</project___start_date_fom>
    <project_client_name>Motorola</project_client_name>
    <project_end_date>1567209600000</project_end_date>
    <project_exec_access>accounting@befoundonline.com;beth.spiegel@befoundonline.com;dan.golden@befoundonline.com;julia.ebner@befoundonline.com;steve.krull@befoundonline.com</project_exec_access>
    <record_owner>beth.spiegel@befoundonline.com</record_owner>
    <related_client/>
    <related_initiative>4651</related_initiative>
    <rollover__ccbm_>1000.00</rollover__ccbm_>
    <sem_prebill_cc__ccbm_>NA</sem_prebill_cc__ccbm_>
    <tech_fee____ccbm_/>
    <tech_fee__ccbm_>0.00</tech_fee__ccbm_>
    <related_initiative___project_id___initiative/>
    <related_initiative___project_id___initiative_name/>
    <project___prebill_management>0</project___prebill_management>
    <project___prebill_media>0</project___prebill_media>
    <project___prebill_tech>0</project___prebill_tech>
    <project___mgmt_fee_calc/>
    <project___deal_type>Renewal</project___deal_type>
    <records>10345</records>
    <add_record/>
    <invoice_nbr>8514</invoice_nbr>
    <po_nbr/>
    <update_id>1557773980205</update_id>
  </record>
  <record rid="10346">
    <project_name>Motorola- GAP - 8/31/2019</project_name>
    <project___contract_status>Signed</project___contract_status>
    <related_project>1068</related_project>
    <project___client_billable_>1</project___client_billable_>
    <billing_month>1512086400000</billing_month>
    <month>1512086400000</month>
    <component>dec, jan, feb</component>
    <initiative>Motorola - Analytics 360 Management</initiative>
    <comment/>
    <client_budget>1000.00</client_budget>
    <bfo__/>
    <___referral_fee/>
    <minimum/>
    <billable_fees>1000.00</billable_fees>
    <external_hourly_rate>195</external_hourly_rate>
    <budget_hours_by_external_rate>5.13</budget_hours_by_external_rate>
    <budget_hour_override/>
    <adjusted_budget_hours>5.13</adjusted_budget_hours>
    <monthly_budgets>10346</monthly_budgets>
    <add_monthly_budget>https://befoundonline.quickbase.com/db/bj22i7dry?a=API_GenAddRecordForm&amp;_fid_56=10346&amp;z=</add_monthly_budget>
    <dbname>Project Central</dbname>
    <client_name/>
    <date_modified>1557773980205</date_modified>
    <initiative___billable>1</initiative___billable>
    <project___client_name>Motorola</project___client_name>
    <project___partner>BFO Direct</project___partner>
    <project_status>Active</project_status>
    <record_id_>10346</record_id_>
    <___actual_costs/>
    <___cogs>0.00</___cogs>
    <___est__costs>0.00</___est__costs>
    <___labor_budget>1000.00</___labor_budget>
    <___labor_budget__ccbm_>0.00</___labor_budget__ccbm_>
    <___overhead>0.00</___overhead>
    <comments2/>
    <___partner_fees/>
    <___target_profit_margin>0.00</___target_profit_margin>
    <___est__costs/>
    <___overhead/>
    <___target_profit_margin/>
    <actual_spend__ccbm_/>
    <actual_spend_fees___total/>
    <add_rollover_to_specified_month__ccbm_>1517443200000</add_rollover_to_specified_month__ccbm_>
    <billable_fees_based_on_actual_spend__ccbm_/>
    <budget_display_name>Motorola- GAP - 8/31/2019 - 12-01-2017 - Motorola - Analytics 360 Management</budget_display_name>
    <budgeted_hours>15.38</budgeted_hours>
    <client___partner_name/>
    <client___vertical/>
    <date_created>1502985937178</date_created>
    <initiative___project___hours_spent_last_month/>
    <initiative___total_hours/>
    <initiative___total_hours_spent/>
    <initiative_name>Analytics 360 Management</initiative_name>
    <initiative_project_name>Motorola- GAP - 8/31/2019</initiative_project_name>
    <last_modified_by>beth.spiegel@befoundonline.com</last_modified_by>
    <month_eom>1514678400000</month_eom>
    <month_fom>1512086400000</month_fom>
    <non_labor_expenses>0.00</non_labor_expenses>
    <operating_budget___do_not_use>1000.00</operating_budget___do_not_use>
    <project____>&lt;div&gt;&lt;img src=&quot;https://images.quickbase.com/si/16/240-triang_green.png&quot;&gt;&lt;/div&gt;</project____>
    <project___account_directors>steve.kozma@befoundonline.com</project___account_directors>
    <project___analysts/>
    <project___billable_fees_ytd>25000.00</project___billable_fees_ytd>
    <project___client_full_name>Motorola Solutions</project___client_full_name>
    <project___managers/>
    <project___ops_level_access>accounting@befoundonline.com;beth.spiegel@befoundonline.com;dan.golden@befoundonline.com;julia.ebner@befoundonline.com;steve.krull@befoundonline.com</project___ops_level_access>
    <project___product>GAC</project___product>
    <project___reason_closed/>
    <project___related_client>212</project___related_client>
    <project___staffing_cost_total/>
    <project___staffing_costs_this_month/>
    <project___start_date>1504224000000</project___start_date>
    <project___start_date_fom>1504224000000</project___start_date_fom>
    <project_client_name>Motorola</project_client_name>
    <project_end_date>1567209600000</project_end_date>
    <project_exec_access>accounting@befoundonline.com;beth.spiegel@befoundonline.com;dan.golden@befoundonline.com;julia.ebner@befoundonline.com;steve.krull@befoundonline.com</project_exec_access>
    <record_owner>beth.spiegel@befoundonline.com</record_owner>
    <related_client/>
    <related_initiative>4651</related_initiative>
    <rollover__ccbm_>1000.00</rollover__ccbm_>
    <sem_prebill_cc__ccbm_>NA</sem_prebill_cc__ccbm_>
    <tech_fee____ccbm_/>
    <tech_fee__ccbm_>0.00</tech_fee__ccbm_>
    <related_initiative___project_id___initiative/>
    <related_initiative___project_id___initiative_name/>
    <project___prebill_management>0</project___prebill_management>
    <project___prebill_media>0</project___prebill_media>
    <project___prebill_tech>0</project___prebill_tech>
    <project___mgmt_fee_calc/>
    <project___deal_type>Renewal</project___deal_type>
    <records>10346</records>
    <add_record/>
    <invoice_nbr>8804</invoice_nbr>
    <po_nbr/>
    <update_id>1557773980205</update_id>
  </record>
</qdbapi>`
	budgets := UnmarshalBudget(sample)
	budget := budgets.Budgets[0]
	fmt.Printf("\n\n Client Budget: %v\n", budget.ClientBudget)
	sugar.Infow("Client Budget", "ClientBudget", budget.ClientBudget)
	if budget.InitiativeBillable != true {
		t.Errorf("Initiative Billable should be true, but got %v", budget.InitiativeBillable)
	}

}

func TestPushToDb(t *testing.T) {
	cnf := ReadConf("dev")
	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", cnf.DbUser, cnf.DbName, cnf.DbPassword, cnf.DbHost)

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	ticket := Auth("uris77@gmail.com", "d43m0nd0g")
	if ticket.Ticket == "" {
		t.Errorf("Ticket should not be empty")
	}
	budgetDoc := GetBudgets(ticket.Ticket, "bg44v3kd9", 2, 0)

	Upsert(db, budgetDoc)
}

func TestCountQbBudgets(t *testing.T) {
	cnf := ReadConf("dev")
	ticket := Auth(cnf.QbUsername, cnf.QbPassword)

	cnt := CountQbBudgets(ticket.Ticket, "bg44v3kd9")

	if cnt.Total < 1 {
		t.Errorf("expected number of budgets to be greater than 0")
	}

	pages := math.Ceil(float64(cnt.Total) / 100)
	log.Printf("Pages: %f | Total: %d", pages, cnt.Total)
}

func TestEverything(t *testing.T) {
	cnf := ReadConf("dev")
	connstr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", cnf.DbUser, cnf.DbName, cnf.DbPassword, cnf.DbHost)

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	ticket := Auth(cnf.QbUsername, cnf.QbPassword)

	cnt := CountQbBudgets(ticket.Ticket, "bg44v3kd9")
	log.Printf("Total number of records: %d", cnt.Total)
	pages := int(math.Ceil(float64(cnt.Total) / 100))
	sz := 100

	log.Printf("Number of Pages: %d", pages)
	i := 0
	var budgetDoc BudgetDoc
	for i < pages {
		budgetDoc = GetBudgets(ticket.Ticket, "bg44v3kd9", sz, i*sz)
		Upsert(db, budgetDoc)
		log.Printf("Iteration: %d", i)
		i++
	}
}

func TestEpochToDate(t *testing.T) {
	d := EpochToDate("1567296000000")
	if d != time.Now() {
		t.Errorf("time failed %v", d)
	}
}
