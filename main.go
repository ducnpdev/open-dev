package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"text/template"
	"time"

	"github.com/dustin/go-humanize"
)

var str = "HDBank: {{.TranDateTime}}\nTK: {{.AccountNoSrc}}\nSố tiền GD: {{.CrDrSign}}{{commaf .TransactionAmount}}{{.Currency}}\nSố dư cuối: {{numerizeBal .ActualBalSrc}}{{.Currency}}\nNội dung GD: {{.Narrative}}\nTG Giao Dịch: {{.OpTs}}\nTG Nhận: {{.CurrentTs}}\nTG Xử lý: {{.ProcessMessageTs}}\nTG Gửi: @SendNotifyTs@"

func main() {

	fmt.Printf("That file is %s.", humanize.Bytes(12000000)) // That file is 83 MB.
	fmt.Println()
	fmt.Println()
	fmt.Println()
	data := NotificationInfo{
		SeqNo:             1,
		TfrSeqNo:          1,
		TransactionAmount: 10,
		TranDateTime:      "2024-05-11 18:46:07",
		TranType:          "123",
		Narrative:         "narrative",
		// OpTs:              "",
		Currency:   "VND",
		TerminalID: "1",
		// Cif:               fmt.Sprint(msg.InternalKey),
		CrDrSign:          "",
		ActualBal:         1000,
		PreviousActualBal: 2000,
		TfrInternalKey:    1,
		ActualBalSrc:      100,
		AccountNoSrc:      "123123",
	}
	funcMap := template.FuncMap{
		"commaf":      humanize.Commaf,
		"numerizeBal": numerizeBal,
	}
	fmt.Println(funcMap)
	notificationTemplate := template.New("notification_template")

	t, err := notificationTemplate.Funcs(funcMap).Parse(str)
	if err != nil {
		fmt.Println(err)
	}

	var output bytes.Buffer
	err = t.Execute(&output, data)
	if err != nil {
		fmt.Println(err)
	}

	content := strings.ReplaceAll(output.String(), "\\n", "\n")
	fmt.Println("content:", content)

	fmt.Println("main", data)
}
func numerizeBal(bal float64) string {
	return humanize.Commaf(math.Abs(bal))
}

type NotificationInfo struct {
	SeqNo             int     `json:"seq_no"`
	ClientNoSrc       string  `json:"client_no_src"`
	ClientNoDes       string  `json:"client_no_des"`
	AccountNoSrc      string  `json:"account_no_src"`
	AccountNoDes      string  `json:"account_no_des"`
	Currency          string  `json:"currency"`
	TransactionAmount float64 `json:"transaction_amount"`
	Narrative         string  `json:"narrative"`      //Nội dung
	TranDateTime      string  `json:"tran_date_time"` //Ngày chuyển khoản
	TranType          string  `json:"tran_type"`
	Mobile            string  `json:"mobile"`
	Email             string  `json:"email"`
	EmailSubject      string  `json:"email_subject"`
	TerminalID        string  `json:"terminal_id"`
	Channel           string  `json:"channel"`
	ActualBalSrc      float64 `json:"actual_bal_src"`
	ActualBalDes      float64 `json:"actual_bal_des"`
	// Cif               string  `json:"cif"`
	CrDrSign          string  `json:"cr_dr_sign"`
	TfrSeqNo          int     `json:"tfr_seq_no"`
	ActualBal         float64 `json:"actual_bal"`
	PreviousActualBal float64 `json:"previous_actual_bal"`
	TfrInternalKey    int     `json:"tfr_internal_key"`
}

func SelectDefault() {
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(start))
	}
}
func Example2() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c) // (1)
	}()
	fmt.Println("Blocking on read...")
	select {
	case <-c: // (2)
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

// multiple channel read
func MulChannel() {
	ch1 := make(chan interface{})
	close(ch1)
	ch2 := make(chan interface{})
	close(ch2)
	var ch1Count, ch2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-ch1:
			ch1Count++
		case <-ch2:
			ch2Count++
		}
	}
	fmt.Printf("ch1Count: %d\n ch2Count: %d\n", ch1Count, ch2Count)
}
