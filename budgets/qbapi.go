package budgets

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type BudgetDoc struct {
	XMLName xml.Name `xml:"qdbapi"`
	ErrCode int      `xml:"errcode"`
	Budgets []Budget `xml:"record"`
}

type Budget struct {
	BudgetId                 int64   `xml:"record_id_"`
	ProjectName              string  `xml:"project_name"`
	ContractStatus           string  `xml:"project___contract_status"`
	ProjectClientBillable    int     `xml:"project___client_billable"`
	BillMonth                string  `xml:"billing_month"`
	Month                    string  `xml:"month"`
	Initiative               string  `xml:"initiative"`
	ClientBudget             float64 `xml:"client_budget"`
	Partner                  string  `xml:"project___partner"`
	ClientName               string  `xml:"client_name"`
	BillableFees             float64 `xml:"billable_fees"`
	ExternalHourlyHourRate   float64 `xml:"external_hourly_rate"`
	Dbname                   string  `xml:"dbname"`
	InitiativeBillable       bool    `xml:"initiative___billable"`
	Component                string  `xml:"component"`
	ProjectStatus            string  `xml:"project_status"`
	DateMod                  string  `xml:"date_modified"`
	Product                  string  `xml:"project___product"`
	ProjectId                int64   `xml:"related_project"`
	RelatedClient            int64   `xml:"project___related_client"`
	Minimum                  float64 `xml:"minimum"`
	PrebillMedia             bool    `xml:"project___prebill_media"`
	BfoPct                   float64 `xml:"bfo__"`
	ReferralFee              float64 `xml:"___referral_fee"`
	ProjectPrebillTech       float64 `xml:"project___prebill_tech"`
	ProjectManagementFeeCalc string  `xml:"project___mgmt_fee_calc"`
	PrebillManagement        bool    `xml:"project___prebill_management"`
}

// Budget hours is the math/avg of billable fees and external hourly rate

func UnmarshalBudget(raw string) BudgetDoc {
	var data BudgetDoc
	if err := xml.Unmarshal([]byte(raw), &data); err != nil {
		log.Fatalf("An error occurred unamrshalling the budgets %q", err)
	}
	return data
}

func (b *Budget) BillingMonth() pq.NullTime {
	if len(b.BillMonth) == 0 {
		return pq.NullTime{Valid: false}
	}
	return pq.NullTime{Time: EpochToDate(b.BillMonth), Valid: true}
}

func (b *Budget) DateModified() pq.NullTime {
	if len(b.DateMod) == 0 {
		return pq.NullTime{Valid: false}
	}
	return pq.NullTime{Time: EpochToDate(b.DateMod), Valid: true}
}

func (b *Budget) BudgetMonth() pq.NullTime {
	if len(b.Month) == 0 {
		return pq.NullTime{Valid: false}
	}
	return pq.NullTime{Time: EpochToDate(b.Month), Valid: true}
}

type ApiCount struct {
	XMLName xml.Name `xml:qdbapi`
	Total   int      `xml:"numMatches"`
}

func UnmarshalApiCount(raw string) ApiCount {
	var data ApiCount
	xml.Unmarshal([]byte(raw), &data)
	return data
}

func BudgetQuery(ticket string, size int, offset int) string {
	return fmt.Sprintf("<qdbapi><ticket>%s</ticket><apptoken>%s</apptoken><query>{1.IR.'last 5 y '}OR{2.IR. 'this y'}</query><includeRids>1</includeRids><clist>a</clist><slist>3</slist><options>num-%d.skp-%d.sortorder-A</options></qdbapi>", ticket, "bue6purcydb2dwcggn43udiycga7", size, offset)
}

func CountQbBudgets(ticket string, table string) ApiCount {
	client := &http.Client{}
	pcQuery := CntQuery(ticket, "bue6purcydb2dwcggn43udiycga7")
	url := fmt.Sprintf("https://befoundonline.quickbase.com/db/%s", table)
	fmt.Println(bytes.NewBuffer([]byte(pcQuery)))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(pcQuery)))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("content-type", "application/xml")
	req.Header.Add("accept", "application/xml")
	req.Header.Add("QUICKBASE-ACTION", "API_DoQueryCount")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return UnmarshalApiCount(bodyString)
}

func CntQuery(ticket string, token string) string {
	return fmt.Sprintf("<qdbapi><ticket>%s</ticket><apptoken>%s</apptoken><query>{1.IR.'last 5 y '}OR{2.IR. 'this y'}</query></qdbapi>", ticket, token)
}

func EpochToDate(e string) time.Time {
	dsec := e[:10]
	dnsec := e[10:]
	sec, _ := strconv.ParseInt(dsec, 10, 64)
	nsec, _ := strconv.ParseInt(dnsec, 10, 64)
	return time.Unix(sec, nsec)
}

type AuthResult struct {
	XMLName xml.Name `xml:"qdbapi"`
	Errcode int      `xml:"errcode"`
	Errtxt  string   `xml:"errtext"`
	Ticket  string   `xml:"ticket"`
}

func Auth(username string, password string) AuthResult {
	client := &http.Client{}
	body := authBody(username, password)
	url := "https://befoundonline.quickbase.com/db/main"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		fmt.Printf("An error ocurred when creating a request to retrieve a ticket")
		log.Fatal(err)
	}

	req.Header.Add("content-type", "application/xml")
	req.Header.Add("accept", "application/xml")
	req.Header.Add("QUICKBASE-ACTION", "API_Authenticate")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	var result AuthResult
	err = xml.Unmarshal([]byte(bodyString), &result)
	if err != nil {
		log.Printf("Failed ot unmarshal the auth body: %s", bodyString)
		log.Fatal(err)
	}
	return result
}

type ApiAuth struct {
	XMLName xml.Name `xml:"qdbapi"`
	Ticket  string   `xml:"ticket"`
}

func authBody(username string, password string) string {
	return fmt.Sprintf("<qdbapi><username>%s</username><password>%s</password></qdbapi>", username, password)
}

func GetBudgets(ticket string, table string, size int, offset int) BudgetDoc {
	client := &http.Client{}
	pcQuery := BudgetQuery(ticket, size, offset)
	log.Printf("QUERY: %s", pcQuery)
	baseUrl := fmt.Sprintf("https://befoundonline.quickbase.com/db/%s", table)
	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer([]byte(pcQuery)))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("content-type", "application/xml")
	req.Header.Add("accept", "application/xml")
	req.Header.Add("QUICKBASE-ACTION", "API_DoQuery")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return UnmarshalBudget(bodyString)
}
