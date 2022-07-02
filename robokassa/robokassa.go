package robokassa

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"investohubBot/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var (
	login  = "investohubhelper"
	outsum = "20127"
	pass1  = "o5EA8XA6BOfU3epVZqE1"
	pass2  = "j2u03ZLl6KsutGrS2vUR"
	url    = "https://auth.robokassa.kz/Merchant/Index.aspx?"

	opstateUrl = "https://auth.robokassa.ru/Merchant/WebService/Service.asmx/OpState?"

	test1 = "ylwVUV1g6Jle4XovR46m"
	test2 = "Hm4vjlJCHqzq87Qx28gj"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetURL(chatId int64) string {
	chat_id := strconv.Itoa(int(chatId))
	return_url := url + "MerchantLogin=" + login + "&InvId=" + chat_id + "&Culture=ru&Encoding=utf-8&OutSum=" + outsum + "&SignatureValue="

	signValue := GetMD5Hash("" + login + ":" + outsum + ":" + chat_id + ":" + pass1 + "")
	return_url += signValue

	fmt.Println(return_url)
	return return_url
}

func CheckPayment(chatId int64) int64 {
	chat_id := strconv.Itoa(int(chatId))
	xmltxt := opstateUrl + "MerchantLogin=" + login + "&InvoiceID=" + chat_id + "&Signature="

	signValue := GetMD5Hash("" + login + ":" + chat_id + ":" + pass2 + "")

	xmltxt += signValue

	fmt.Println(xmltxt)

	resp, err := http.Get(xmltxt)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var osr models.OperationStateResponse

	err = xml.Unmarshal(body, &osr)
	if err != nil {
		return -1
	}

	/*fmt.Println(osr.Result.Code)
	fmt.Println(osr.Result.Description)*/

	if osr.Result.Code != 0 {
		return -1
	} else {
		return osr.State.Code
	}

	/*sb := string(body)
	fmt.Println(sb)*/
}
