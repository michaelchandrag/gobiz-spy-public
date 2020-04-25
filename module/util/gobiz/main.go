package gobiz

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

type (
	GobizResponse struct {
		AccessToken 		string 				`json:"access_token,omitempty"`
		RefreshToken 		string 				`json:"refresh_token,omitempty"`
		ResponseData 		GobizData 			`json:"data,omitempty"`
		ResponseUser 		GobizUser 			`json:"user,omitempty"`
		ResponseErrors 		[]GobizErrors 		`json:"errors,omitempty"`
		Success 			*bool				`json:"success,omitempty"`
		Total 				int 				`json:"total,omitempty"`
		ResponseHits		[]GobizHits 		`json:"hits,omitempty"`
	}

	GobizCustomer struct {
		CustomerID 		string 			`json:"customer_id"`
		CustomerName 	string 			`json:"customer_name"`
		CustomerPhone 	string 			`json:"customer_phone"`
		CustomerEmail 	string 			`json:"customer_email"`
	}

	GobizTransaction struct {
		OrderID 		string 				`json:"order_id"`
		Customer 		GobizCustomer 		`json:"customer"`
		OrderedAt 		string 				`json:"ordered_at"`
	}

	GobizHits struct {
		ResponseTransaction 	GobizTransaction 		`json:"transaction"`
		// OrderedAt 				string 					`json:"ordered_at"`
	}

	GobizData struct {
		OTPToken 		string 				`json:"otp_token,omitempty"`
		OTP 			string 				`json:"otp,omitempty"`
		ResponseOrder 	[]GobizOrder 		`json:"hits,omitempty"`
	}

	GobizOrder struct {
		OrderedAt 			string 					`json:"ordered_at"`
		ProductSpecific 	ProductSpecific 		`json:"product_specific"`
	}

	ProductSpecific struct {
		GoResto 		GoResto 			`json:"goresto"`
	}

	GoResto struct {
		CustomerID 		string 			`json:"customer_id"`
		CustomerName 	string 			`json:"customer_name"`
		CustomerPhone 	string 			`json:"customer_phone"`
		CustomerEmail  	string 			`json:"customer_email"`
	}

	GobizErrors struct {
		Message 		string 			`json:"message"`
	}

	GobizUser struct {
		ID 				int 			`json:"id"`
		Phone 			string 			`json:"phone"`
		MerchantID 		string 			`json:"merchant_id"`
	}
)

func RequestOtp (phoneNumber string) (result GobizResponse, err error) {
	type Request struct {
		ClientId  		string 		`json:"client_id"`
		CountryCode  	string 		`json:"country_code"`
		PhoneNumber  	string 		`json:"phone_number"`
	}

	body := Request{
		ClientId: "YEZympJ5WqYRh7Hs",
		CountryCode: "62",
		PhoneNumber: phoneNumber,
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	url := "https://api.gobiz.co.id/goid/login/request"
	
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, buf)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	// printResponse(*res)
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

func RequestToken (otpToken string, otp string) (result GobizResponse, err error) {
	type Request struct {
		ClientId 	string 		`json:"client_id"`
		GrandType 	string 		`json:"grant_type"`
		Data 		GobizData 	`json:"data"`
	}
	body := Request{
		ClientId: "YEZympJ5WqYRh7Hs",
		GrandType: "otp",
		Data: GobizData{
			OTPToken: otpToken,
			OTP: otp,
		},
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	url := "https://api.gobiz.co.id/goid/token"

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, buf)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	// printResponse(*res)
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

func RequestProfile (accessToken string) (result GobizResponse, err error) {
	url := "https://api.midtrans.com/v1/users/me"
	
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Authentication-Type", "go-id")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&result)
	return result, nil

}

func RequestOrder (accessToken string, orderId string) (result GobizResponse, err error) {
	type Match struct {
		OrderNumber 	string  		`json:"order_number"`
	}

	type Query struct {
		Match		Match 		`json:"match"`
	}

	type Request struct {
		Query 		Query 		`json:"query"`
	}

	body := Request{
		Query: Query{
			Match{
				OrderNumber: orderId,
			},
		},
	}

	url := "https://api.gobiz.co.id/cosmo/v1/orders/search"
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Authentication-Type", "go-id")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&result)
	return result, nil

}

func RequestTransactions (accessToken string, merchantId string, startDate string, endDate string, from int) (result GobizResponse, err error) {
	type Order struct {
		Value 		string 		`json:"order"`
	}

	type OrderKey struct {
		Key 		Order 			`json:"transaction.transaction_time"`
	}

	type Clauses struct {
		Field 		string 			`json:"field,omitempty"`
		Op 			string 			`json:"op"`
		Value 		interface{} 	`json:"value,omitempty"`
		Clauses 	[]Clauses 		`json:"clauses,omitempty"`
	}

	type Query struct {
		Clauses 	[]Clauses 		`json:"clauses"`
		Op 			string 			`json:"op"`
	}

	type Request struct {
		From  		int 				`json:"from"`
		Size  		int 				`json:"size"`
		Query 		[]Query 			`json:"query"`
		Sort 		OrderKey 			`json:"sort"`
		Includes 	[]string 			`json:"includes"`
	}

	body := Request{
		From: from,
		Size: 20,
		Sort: OrderKey{
			Key: Order {
				Value: "desc",
			},
		},
		Includes: []string{"shares"},
		Query: []Query{
			Query{
				Clauses: []Clauses{
					Clauses{
						Field: "transaction.status",
						Op: "in",
						Value: []string{"settlement", "capture", "refund", "partial_refund"},
					},
					Clauses{
						Field: "transaction.payment_type",
						Op: "in",
						Value: []string{"qris", "gopay", "cash", "offline_ovo", "offline_telkomsel_cash", "offline_credit_card", "offline_debit_card", "credit_card"},
					},
					Clauses{
						Op: "not",
						Clauses: []Clauses{
							Clauses{
								Op: "or",
								Clauses: []Clauses{
									Clauses{
										Field: "source",
										Op: "in",
										Value: []string{"gopay_online", "GOSAVE_ONLINE", "GoSave", "GODEALS_ONLINE"},
									},
									Clauses{
										Field: "gopay.source",
										Op: "in",
										Value: []string{"gopay_online", "GOSAVE_ONLINE", "GoSave", "GODEALS_ONLINE"},
									},
								},
							},
						},
					},
					Clauses{
						Field: "transaction.transaction_time",
						Op: "gte",
						Value: startDate,
					},
					Clauses{
						Field: "transaction.transaction_time",
						Op: "lte",
						Value: endDate,
					},
					Clauses{
						Field: "transaction.merchant_id",
						Op: "equal",
						Value: merchantId,
					},
				},
				Op: "and",
			},
		},
	}
	url := "https://api.midtrans.com/v1/payments/search"
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, buf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Authentication-Type", "go-id")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}

func printResponse (res http.Response) {
	fmt.Println(res.Status)
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bodyBytes))
}