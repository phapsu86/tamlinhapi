package controllers

import (
	"bytes"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	//	"fmt"
	"io/ioutil"
	"net/http"

	//"strconv"

	//"github.com/gorilla/mux"
	"github.com/phapsu86/tamlinhapi/api/auth"
	"github.com/phapsu86/tamlinhapi/api/models"
	"github.com/phapsu86/tamlinhapi/api/responses"
	"github.com/phapsu86/tamlinhapi/api/utils/formaterror"
	"github.com/phapsu86/tamlinhapi/api/utils/formatresult"
)

type ParamTopupPoint struct {
	PaymentType string `json:"payment_type"`
	Amount      uint64 `json:"status"`
	Page        int    `json:"page"`
}

type ParamMomo struct {
	OrderID     uint64 `json:"order_id"`
	Data        string `json:"data"`
	PhoneNumber string `json:"phonenumber"`
}

type PayApp struct {
	PartnerCode    string  `json:"partnerCode"`
	PartnerRefID   string  `json:"partnerRefId"`
	CustomerNumber string  `json:"customerNumber"`
	AppData        string  `json:"appData"`
	Hash           string  `json:"hash"`
	Version        float64 `json:"version"`
	PayType        int     `json:"payType"`
	Description    string  `json:"description"`
	ExtraData      string  `json:"extraData"`
}

type JsonData struct {
	PartnerCode    string `json:"partnerCode"`
	PartnerRefID   string `json:"partnerRefId"`
	PartnerTransID string `json:"partnerTransId"`
	Amount         uint64 `json:"amount"`
}

var pubKeyData = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkpa+qMXS6O11x7jBGo9W
3yxeHEsAdyDE40UoXhoQf9K6attSIclTZMEGfq6gmJm2BogVJtPkjvri5/j9mBnt
A8qKMzzanSQaBEbr8FyByHnf226dsLt1RbJSMLjCd3UC1n0Yq8KKvfHhvmvVbGcW
fpgfo7iQTVmL0r1eQxzgnSq31EL1yYNMuaZjpHmQuT24Hmxl9W9enRtJyVTUhwKh
tjOSOsR03sMnsckpFT9pn1/V9BE2Kf3rFGqc6JukXkqK6ZW9mtmGLSq3K+JRRq2w
8PVmcbcvTr/adW4EL2yc1qk9Ec4HtiDhtSYd6/ov8xLVkKAQjLVt7Ex3/agRPfPr
NwIDAQAB
-----END PUBLIC KEY-----
`)

func (server *Server) CreatePaymentTopUpPoint(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item := ParamTopupPoint{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil || uid == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//Create Order
	mdOd := models.Order{}
	mdOd.Amount = item.Amount
	mdOd.UserID = uid
	mdOd.Status = 0
	mdOd.OrderType = 0
	mdOd.Prepare()
	orderData, err := mdOd.SaveOrder(server.DB)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return

	}

	var rs interface{}
	rs = orderData

	result := formatresult.ReturnGlobalArray(rs)
	responses.JSON(w, http.StatusOK, result)

}

func (server *Server) CreatPaymentMoMo(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item := ParamMomo{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil || uid == 0 {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Lấy thông tin Order  ====

	//Create transaction
	mdTrans := models.TransactionHistory{}
	mdTrans.Prepare()
	mdTrans.OrderID = item.OrderID
	mdTrans.Amount = 1000
	mdTrans.TranCode = "TOPUP_POINT"
	mdTrans.Status = 0
	transData, err := mdTrans.SaveTransactionHistory(server.DB)
	if err != nil {
		err := formaterror.ReturnErr(err)
		responses.JSON(w, http.StatusOK, err)
		return

	}

	//=============================

	var orderId = "15"
	//var requestId = "16"
	var endpoint = "https://test-payment.momo.vn/pay/app"
	var partnerCode = "MOMOIQA420180417"
	//var accessKey = "SvDmj2cOTYZmQQ3H"
	//var serectkey = "PPuDXq1KowPT1ftR8DvlQTHhC03aul17"
	//var orderInfo = "momo all-in-one"
	//	var returnUrl = "https://developers.momo.vn/"
	//var notifyurl = "https://webhook.site/3c5b6488-a159-4f8d-b038-29eed82fab1e"
	//	var amount = 1000
	//var requestType = "captureMoMoWallet"
	var extraData = "merchantName=;merchantId="
	posHash := JsonData{}
	posHash.Amount = 1000
	posHash.PartnerCode = partnerCode
	posHash.PartnerRefID = orderId
	posHash.PartnerTransID = strconv.FormatUint(transData.ID, 10)

	//1. TODO - Parse RSA Public key string to memory Public key
	block, _ := pem.Decode([]byte(pubKeyData))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}
	pkixPub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}
	switch pkixPub := pkixPub.(type) {
	case *rsa.PublicKey:
		fmt.Println("RSA Public Key: {}")
		fmt.Println(pkixPub)
		fmt.Println()
	case *dsa.PublicKey:
		fmt.Println("DSA Public Key: {}")
		fmt.Println(pkixPub)
	case *ecdsa.PublicKey:
		fmt.Println("ECDSA Public Key: {}")
		fmt.Println(pkixPub)
	default:
		panic("Unknow Public Key")
	}
	var publicKey *rsa.PublicKey
	//TODO - Result public key in memory
	publicKey = pkixPub.(*rsa.PublicKey)

	var rawJashJson []byte
	rawJashJson, err = json.Marshal(posHash)
	if err != nil {
		log.Println(err)
	}
	//END define a Json before hash name `hashJson`

	//2. TODO - Encrypt PKCS1V15 RSA encryption using `publicKey` in memory
	randomReader := rand.Reader
	ciphertext, err := rsa.EncryptPKCS1v15(randomReader, publicKey, rawJashJson)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return
	}
	var hash = b64.StdEncoding.EncodeToString(ciphertext)

	postToMoMo := PayApp{}
	postToMoMo.AppData = item.Data
	postToMoMo.PartnerCode = partnerCode
	postToMoMo.PartnerRefID = orderId
	postToMoMo.CustomerNumber = item.PhoneNumber
	postToMoMo.AppData = item.Data
	postToMoMo.Hash = hash
	postToMoMo.Version = 2.0
	postToMoMo.PayType = 3
	postToMoMo.Description = "TOPUP_POINT"
	postToMoMo.ExtraData = extraData

	var posPayappJson []byte
	posPayappJson, err = json.Marshal(postToMoMo)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Payload with hash : {}")
	fmt.Println(string(posPayappJson))

	//send HTTP to momo endpoint
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(posPayappJson))
	if err != nil {
		log.Fatalln(err)
	}

	//result
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println("Response from Momo: ", result)

}
