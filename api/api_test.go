package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var systemStatusJSON = `{
	"message":"",
	"valid_version":true,
	"system_running":true,
	"skip_phrase":true,
	"timestamp":1371327425000
}`

func TestSystemStatusIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1", "", []byte(systemStatusJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL}

	if _, err := client.SystemStatus(); err != nil {
		t.Fatal(err)
	}
}

var loginJSON = `{
	"country":"SE",
	"expires_in":300,
	"session_key":"441ff696b7bd75fbe50add3e2e728eb761596f1b",
	"environment":"test",
	"private_feed":{"port":443,"hostname":"priv.api.test.nordnet.se","encrypted":true},
	"public_feed":{"port":443,"hostname":"pub.api.test.nordnet.se","encrypted":true}
}`

func TestLoginIntegration(t *testing.T) {
	ts := setupTestServer(t, "POST", "/v1/login?auth=SECRET&service=TEST", "", []byte(loginJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, Credentials: "SECRET", Service: "TEST", SessionKey: ""}
	if _, err := client.Login(); err != nil {
		t.Fatal(err)
	}
}

var logoutJSON = `{"logged_in":false}`

func TestLogoutIntegrationt(t *testing.T) {
	ts := setupTestServer(t, "DELETE", "/v1/login/SESSIONKEY", "SESSIONKEY", []byte(logoutJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}
	if _, err := client.Logout(); err != nil {
		t.Fatal(err)
	}
}

var touchJSON = `{"logged_in":true}`

func TestTouchIntegration(t *testing.T) {
	ts := setupTestServer(t, "PUT", "/v1/login/SESSIONKEY", "SESSIONKEY", []byte(touchJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}
	if _, err := client.Touch(); err != nil {
		t.Fatal(err)
	}
}

var realtimeAccessJSON = `[
	{"marketID":"44","level":2},
	{"marketID":"11","level":2},
	{"marketID":"34","level":1},
	{"marketID":"12","level":2}
]`

func TestReatimeAccessIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/realtime_access", "SESSIONKEY", []byte(realtimeAccessJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}
	if _, err := client.RealtimeAccess(); err != nil {
		t.Fatal(err)
	}
}

var newsSourcesJSON = `[
	{"name":"Dow Jones News","imageurl":"/now/images/loggaDJN.gif","code":"djn","sourceid":3,"level":"REALTIME"},
	{"name":"OMX","imageurl":"/now/images/loggaOmxnews.gif","code":"omxnews","sourceid":7,"level":"REALTIME"},
	{"name":"Thomson Reuters","imageurl":"/now/images/loggaHugin.gif","code":"hugin","sourceid":9,"level":"REALTIME"}
]`

func TestNewsSourcesIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/news_sources", "SESSIONKEY", []byte(newsSourcesJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}
	if _, err := client.NewsSources(); err != nil {
		t.Fatal(err)
	}
}

// FIXME: these are taken from the docs, not real values
var newsItemsJSON = `[
	{"datetime":"2010-03-01 10:40:19 UTC","headline":"LONDON MARKETS: BP Falls","itemid":159619003,"sourceid":3,"type":"NEWS"},
	{"datetime":"2010-03-01 10:40:19 UTC","headline":"LONDON MARKETS: BP Falls","itemid":159619003,"sourceid":3,"type":"NEWS"}
]`

func TestNewsItemsIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/news_items", "SESSIONKEY", []byte(newsItemsJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}
	if _, err := client.NewsItems(nil); err != nil {
		t.Fatal(err)
	}
}

// FIXME: these are taken from the docs, not real values
var newsItemJSON = `{
	"datetime":"2010-03-01 10:40:19 UTC",
	"headline":"Danske Equities",
	"body":"test",
	"itemid":4711,
	"lang":"da",
	"preamble":"test",
	"sourceid":6,
	"type":"NEWS"
}`

func TestNewsItemIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/news_items/4711", "SESSIONKEY", []byte(newsItemJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}
	if _, err := client.NewsItem(4711); err != nil {
		t.Fatal(err)
	}
}

var accountsJSON = `[
	{"alias":null,"default":"true","id":"1000000"}
]`

func TestAccountsIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/accounts", "SESSIONKEY", []byte(accountsJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}
	if _, err := client.Accounts(); err != nil {
		t.Fatal(err)
	}
}

var accountJSON = `{
	"ownCapitalMorning":"1000000.0",
	"accountCurrency":"SEK",
	"ownCapital":"1000000.0",
	"futureSum":"0.0",
	"forwardSum":"0.0",
	"collateral":"0.0",
	"tradingPower":"948000.0",
	"interest":"0.0",
	"pawnValue":"0.0",
	"accountSum":"1000000.0",
	"loanLimit":"1000000.0",
	"fullMarketvalue":"0.0"
}`

func TestAccountIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/accounts/1000000", "SESSIONKEY", []byte(accountJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}
	if _, err := client.Account("1000000"); err != nil {
		t.Fatal(err)
	}
}

var accountLedgersJSON = `[
	{
		"accountSumAcc":"1000000.0",
		"accIntCred":"0.0",
		"currency":"SEK",
		"accIntDeb":"0.0",
		"accountSum":"1000000.0"
	}
]`

func TestAccountLedgersIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/accounts/1000000/ledgers", "SESSIONKEY", []byte(accountLedgersJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.AccountLedgers("1000000"); err != nil {
		t.Fatal(err)
	}
}

// FIXME these are taken from the docs, not real values
var accountPositionsJSON = `[
	{
		"acqPrice":"700.1524",
		"acqPriceAcc":"700.1524",
		"pawnPercent":"85",
		"qty":"9.0",
		"marketValue":"642.6",
		"marketValueAcc":"642.6",
		"instrumentID":{
			"mainMarketId":"11",
			"identifier":"101",
			"type":"A",
			"currecy": "SEK",
			"mainMarketPrice":"55"
		}
	}
]`

func TestAccountPositionsIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/accounts/1000000/positions", "SESSIONKEY", []byte(accountPositionsJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.AccountPositions("1000000"); err != nil {
		t.Fatal(err)
	}
}

var accountOrdersJSON = `[
	{
		"priceCondition":"LIMIT",
		"validity":{"validUntil":1370876700000,"type":"DAY"},
		"price":{"value":65.0,"curr":"SEK"},
		"side":"BUY",
		"orderID":683772,
		"volumeCondition":"NORMAL",
		"tradedVolume":0.0,
		"instrumentID":{"marketID":11,"identifier":"101"},
		"orderState":"LOCAL",
		"accno":9210370,
		"openVolume":0.0,
		"volume":100.0,
		"actionState":"INS_PEND",
		"activationCondition":{"type":"NONE"},
		"modDate":1370797680194
	}
]`

func TestAccountOrdersIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/accounts/1000000/orders", "SESSIONKEY", []byte(accountOrdersJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.AccountOrders("1000000"); err != nil {
		t.Fatal(err)
	}
}

var accountTradesJSON = `[
	{
		"securityTrade":{
			"tradeID":"B8118-20130603",
			"price":{"value":"146","curr":"SEK"},
			"volume":"2",
			"tradetime":"12:06:06",
			"instrumentID":{"marketID":"11","identifier":"3966"},
			"accno":"9210329",
			"counterparty":"MCF",
			"orderID":"683168",
			"side":"BUY"
		}
	}
]`

func TestAccountTradesIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/accounts/1000000/trades", "SESSIONKEY", []byte(accountTradesJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.AccountTrades("1000000"); err != nil {
		t.Fatal(err)
	}
}

var instrumentsJSON = `[
	{
		"type":"A",
		"longname":"Ericsson A",
		"marketID":"11",
		"country":"SE",
		"shortname":"ERIC A",
		"marketname":"OMX Stockholm",
		"isinCode":"SE0000108649",
		"identifier":"100",
		"currency":"SEK"
	}
]`

func TestInstrumentsIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/instruments?country=SE&query=ERI&type=A", "SESSIONKEY", []byte(instrumentsJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.Instruments(&Params{"query": "ERI", "type": "A", "country": "SE"}); err != nil {
		t.Fatal(err)
	}
}

var instrumentJSON = `{
	"type":"A",
	"longname":"Ericsson B",
	"marketID":"11",
	"country":"SE",
	"shortname":"ERIC B",
	"multiplier":"1",
	"marketname":"OMX Stockholm",
	"ticksizeID":"11002",
	"isinCode":"SE0000108656",
	"identifier":"101",
	"currency":"SEK"
}`

func TestInstrumentIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/instruments?identifier=101&marketID=11", "SESSIONKEY", []byte(instrumentJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.Instrument(&Params{"identifier": "101", "marketID": "11"}); err != nil {
		t.Fatal(err)
	}
}

var chartDataJSON = `[
	{"timestamp":"09:38","change":12.18,"volume":1000,"float":82.0}
]`

func TestChartDataIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/chart_data?identifier=101&marketID=11", "SESSIONKEY", []byte(chartDataJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.ChartData(&Params{"identifier": "101", "marketID": "11"}); err != nil {
		t.Fatal(err)
	}
}

var listsJSON = `[
	{"name":"First North SE","country":"SE","id":"6"},
	{"name":"Small Cap Copenhagen","country":"DK","id":"16"}
]`

func TestListsIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/lists", "SESSIONKEY", []byte(listsJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.Lists(); err != nil {
		t.Fatal(err)
	}
}

var listJSON = `[
	{"shortname":"WISE","marketID":"11","identifier":"40017"},
	{"shortname":"WINT","marketID":"11","identifier":"43370"}
]`

func TestListIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/lists/6", "SESSIONKEY", []byte(listJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.List(6); err != nil {
		t.Fatal(err)
	}
}

var marketsJSON = `[
	{
		"name":"Nasdaq",
		"country":"US",
		"marketID":"19",
		"ordertypes":[
			{"text":"Normal order","type":"NORMAL"}
		]
	}
]`

func TestMarketsIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/markets", "SESSIONKEY", []byte(listJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.Markets(); err != nil {
		t.Fatal(err)
	}
}

var marketTradingDaysJSON = `[
	{"date":"2013-06-18","display_date":"2013-06-18"},
	{"date":"2013-06-19","display_date":"2013-06-19"}
]`

func TestMarketTradingDaysIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/markets/11/trading_days", "SESSIONKEY", []byte(marketTradingDaysJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.MarketTradingDays(11); err != nil {
		t.Fatal(err)
	}
}

var indicesJSON = `[
	{
		"type":"INDEX",
		"longname":"OBX",
		"source":"OSE",
		"country":"NO",
		"imageurl":"/now/images/flaggaNoSmall.gif",
		"id":"XOBX"
	},
	{
		"type":"COMMODITY",
		"longname":"Aluminium 3M USD",
		"source":"SIX",
		"id":"B-ALUM-3M"
	}
]`

func TestIndicesIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/indices", "SESSIONKEY", []byte(indicesJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.Indices(); err != nil {
		t.Fatal(err)
	}
}

var ticksizesJSON = `[
	{"tick":0.0001,"above":0.0,"decimals":4},
	{"tick":0.0005,"above":0.5,"decimals":4},
	{"tick":0.001,"above":1.0,"decimals":3}
]`

func TestTicksizesIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/ticksizes/11002", "SESSIONKEY", []byte(ticksizesJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.Ticksizes(11002); err != nil {
		t.Fatal(err)
	}
}

var derivateCountriesJSON = `["SE","FI","NO"]`

func TestDerivateContriesIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/derivatives/A", "SESSIONKEY", []byte(derivateCountriesJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.DerivateCountries("A"); err != nil {
		t.Fatal(err)
	}
}

var derivateUnderlyingsJSON = `[
	{"shortname":"OMXS30","marketID":"11","identifier":"OMXS30"},
	{"shortname":"TLSN","marketID":"11","identifier":"5095"},
	{"shortname":"ERIC B","marketID":"11","identifier":"101"},
	{"shortname":"NOKI SEK","marketID":"11","identifier":"39854"}
]`

func TestDerivateUnderyingsIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/derivatives/O/underlyings/SE", "SESSIONKEY", []byte(derivateUnderlyingsJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.DerivateUnderlyings("O", "SE"); err != nil {
		t.Fatal(err)
	}
}

// FIXME these are taken from the doc, not real values
var derivateJSON = `[
	{
		"shortname":"ERI1N 60SHB",
		"multiplier":"1",
		"strikeprice":"60.000000",
		"expirydate":"2011-02-18 00:00:00",
		"marketID":"11",
		"expirytype":"european",
		"kind":"WNT",
		"identifier":"76987",
		"currency":"SEK",
		"callPut":"Warrant Put"
	}
]`

func TestDerivativesIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/derivatives/WNT/derivatives", "SESSIONKEY", []byte(derivateJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	params := &Params{"identifier": "101", "marketID": "11"}
	if _, err := client.Derivatives("WNT", params); err != nil {
		t.Fatal(err)
	}
}

var relatedMarketsJSON = `[
	{"marketID":11,"identifier":"101"},
	{"marketID":30,"identifier":"1965"}
]`

func TestRelatedMarketsIntegration(t *testing.T) {
	ts := setupTestServer(t, "GET", "/v1/related_markets?identifier=101&marketID=11", "SESSIONKEY", []byte(relatedMarketsJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	params := &Params{"identifier": "101", "marketID": "11"}
	if _, err := client.RelatedMarkets(params); err != nil {
		t.Fatal(err)
	}
}

var createOrderJSON = `{
	"orderID":684870,
	"resultCode":"OK",
	"orderState":"LOCAL",
	"accNo":1000000,
	"actionState":"INS_PEND"
}`

func TestCreateOrderIntegration(t *testing.T) {
	ts := setupTestServer(t, "POST", "/v1/accounts/1000000/orders?currency=SEK&identifier=101&marketID=11&price=65&side=buy&volume=100", "SESSIONKEY", []byte(createOrderJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	params := &Params{"identifier": "101", "marketID": "11", "price": "65", "volume": "100", "side": "buy", "currency": "SEK"}
	if _, err := client.CreateOrder("1000000", params); err != nil {
		t.Fatal(err)
	}
}

var updateOrderJSON = `{
	"orderID":684870,
	"resultCode":"OK",
	"orderState":"ON_MARKET",
	"accNo":1000000,
	"actionState":"MOD_PEND"
}`

func TestUpdateOrderIntegration(t *testing.T) {
	ts := setupTestServer(t, "PUT", "/v1/accounts/1000000/orders/684870?price=68", "SESSIONKEY", []byte(updateOrderJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.UpdateOrder("1000000", 684870, &Params{"price": "68"}); err != nil {
		t.Fatal(err)
	}
}

var deleteOrderJSON = `{
	"orderID":684870,
	"resultCode":"OK",
	"orderState":"ON_MARKET",
	"accNo":9210370,
	"actionState":"DEL_PEND"
}`

func TestDeleteOrderIntegration(t *testing.T) {
	ts := setupTestServer(t, "DELETE", "/v1/accounts/1000000/orders/684870", "SESSIONKEY", []byte(deleteOrderJSON))
	defer ts.Close()

	client := &APIClient{URL: ts.URL, SessionKey: "SESSIONKEY"}

	if _, err := client.DeleteOrder("1000000", 684870); err != nil {
		t.Fatal(err)
	}
}

func setupTestServer(t *testing.T, method, path, session string, stubData []byte) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Fatal(errors.New(fmt.Sprintln("Method was expected to be:", method, "got:", r.Method)))
		} else if r.RequestURI != path {
			t.Fatal(errors.New(fmt.Sprintln("Path was expected to be:", path, "got:", r.RequestURI)))
		} else if auth := r.Header.Get("Authorization"); auth != "" {
			if decoded, err := base64.StdEncoding.DecodeString(auth[6:]); err != nil {
				t.Fatal(err)
			} else if userpass := session + ":" + session; userpass != string(decoded) {
				t.Fatal(errors.New(fmt.Sprintln("Session was expected to be:", userpass, "got:", string(decoded))))
			}
		}

		w.Write(stubData)
	})

	return httptest.NewServer(handler)
}