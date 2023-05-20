package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	sumsubmodels "github.com/johnyeocx/usual/server2/db/models/sumsub_models"
)

const URL = "https://api.sumsub.com"
const SumsubAppToken = "sbx:6IT1WBqqwRHJJzUSZnA7TE2R.l8CAqpg3prnfaj4W3LX1GuoSJk0EyuPt"
const SumsubSecretKey = "KX7ttypyiiGMYTiQa0yEtqm51DgatRsN"

func main() {
	var levelName = "basic-kyc-level"

	var fixedInfo = sumsubmodels.Info{}
	fixedInfo.Country = "GBR"
	fixedInfo.FirstName = "Test First Name"
	var applicant = sumsubmodels.Applicant{
		FixedInfo: fixedInfo,
		ExternalUserID: "35",
	}

	fmt.Println("Level Name: ", levelName)
	fmt.Println("Applicant Name: ", applicant)
	
	// applicant = CreateApplicant(applicant, levelName)
	// applicant = GetApplicantInfo(applicant)


	// idDoc := AddDocument("646507dd69edf812c35d53d2")
	// fmt.Println(idDoc)

	// // https://developers.sumsub.com/api-reference/#getting-applicant-data

	// // https://developers.sumsub.com/api-reference/#access-tokens-for-sdks
	accessToken := GenerateAccessToken("55547a13-929f-4307-9421-28b4af103ebf", levelName)
	fmt.Println(accessToken.Token)
}

func GenerateAccessToken(externalUserId string, levelName string) sumsubmodels.AccessToken {

	b, err := _makeSumsubRequest("/resources/accessTokens?userId="+externalUserId+"&levelName="+levelName,
		"POST",
		"application/json",
		[]byte(""))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	ioutil.WriteFile("generateAccessToken.json", b, 0777)

	var token sumsubmodels.AccessToken
	err = json.Unmarshal(b, &token)

	return token
}


func CreateApplicant(applicant sumsubmodels.Applicant, levelName string) sumsubmodels.Applicant {
	postBody, _ := json.Marshal(applicant)

	b, err := _makeSumsubRequest(
		"/resources/applicants?levelName="+levelName,
		"POST",
		"application/json",
		postBody)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
	ioutil.WriteFile("createApplicant.json", b, 0777)

	var ac sumsubmodels.Applicant
	err = json.Unmarshal(b, &ac)
	if err != nil {
		log.Fatal(err)
	}

	return ac
}

func AddDocument(applicantId string) sumsubmodels.IdDoc {
	file, err := os.Open("./test_id.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	meta, _ := json.Marshal(sumsubmodels.IdDoc{
		IdDocType: "PASSPORT",
		Country:   "GBR",
	})

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	var fw io.Writer
	if fw, err = w.CreateFormFile("content", file.Name()); err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		log.Fatal(err)
	}

	if fw, err = w.CreateFormField("metadata"); err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(fw, strings.NewReader(string(meta))); err != nil {
		log.Fatal(err)
	}
	w.Close()

	resp, _ := _makeSumsubRequest(
		"/resources/applicants/"+applicantId+"/info/idDoc",
		"POST",
		w.FormDataContentType(),
		b.Bytes(),
	)

	var doc sumsubmodels.IdDoc
	err = json.Unmarshal(resp, &doc)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func GetApplicantInfo(applicant sumsubmodels.Applicant) sumsubmodels.Applicant {
	p := fmt.Sprintf("/resources/applicants/%s/one", applicant.ID)
	b, err := _makeSumsubRequest(
		p,
		"GET",
		"application/json",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("getApplicant.json", b, 0777)

	var r sumsubmodels.Applicant
	err = json.Unmarshal(b, &r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)

	return r
}


//X-App-Token - an App Token that you generate in our dashboard
//X-App-Access-Sig - signature of the request in the hex format (see below)
//X-App-Access-Ts - number of seconds since Unix Epoch in UTC
func _makeSumsubRequest(path, method, contentType string, body []byte) ([]byte, error) {

	request, err := http.NewRequest(method, URL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	ts := fmt.Sprintf("%d", time.Now().Unix())
	fmt.Println("App Token:", SumsubAppToken)
	fmt.Println("Secret Key:", SumsubSecretKey)
	request.Header.Add("X-App-Token", SumsubAppToken)
	request.Header.Add("X-App-Access-Sig", _sign(ts, SumsubSecretKey, method, path, &body))
	request.Header.Add("X-App-Access-Ts", ts)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", contentType)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func _sign(ts string, secret string, method string, path string, body *[]byte) string {
	hash := hmac.New(sha256.New, []byte(secret))
	data := []byte(ts + method + path)

	if body != nil {
		data = append(data, *body...)
	}

	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}