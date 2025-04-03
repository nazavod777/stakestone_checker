package accountsParser

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"main/pkg/global"
	"main/pkg/types"
	"main/pkg/util"
	"strings"
)

func parseBalance(accountData types.AccountData) float64 {
	var err error

	for {
		client := util.GetClient()

		req := fasthttp.AcquireRequest()

		req.SetRequestURI(fmt.Sprintf("https://airdrop.stakestone.io/api/credentials?walletAddress=%s&batchId=0",
			strings.ToLower(accountData.AccountAddress.String())))
		req.Header.Set("accept", "*/*")
		req.Header.Set("accept-language", "ru,en;q=0.9,vi;q=0.8,es;q=0.7,cy;q=0.6")
		req.Header.Set("origin", "https://app.kiloex.io")
		req.Header.SetReferer("https://app.kiloex.io/")
		req.Header.SetMethod("GET")
		req.Header.SetContentType("application/json")

		resp := fasthttp.AcquireResponse()

		if err = client.Do(req, resp); err != nil {
			log.Printf("[%d/%d] | %s | Error When Doing Request When Parsing Balance: %s",
				global.CurrentProgress, global.TargetProgress, accountData.AccountLogData, err)

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		dropAmount := gjson.Get(string(resp.Body()), "tokenQualified")
		if !dropAmount.Exists() {
			log.Printf("[%d/%d] | %s | Wrong Response When Parsing Balance: %s",
				global.CurrentProgress, global.TargetProgress, accountData.AccountLogData, string(resp.Body()))

			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(resp)
			continue
		}

		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)

		return dropAmount.Float()
	}
}

func ParseAccount(accountData types.AccountData) {
	accountBalance := parseBalance(accountData)

	log.Printf("[%d/%d] | %s | %g $STO",
		global.CurrentProgress, global.TargetProgress, accountData.AccountLogData, accountBalance)

	if accountBalance > 0 {
		util.AppendFile("with_balances.txt",
			fmt.Sprintf("%s | %g $STO\n", accountData.AccountLogData, accountBalance))
	}
}
