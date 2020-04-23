package second

import (
	"os"
	"net/http"
	"sync"

	gobiz "github.com/michaelchandrag/gobiz-spy/module/util/gobiz"
	"github.com/gin-gonic/gin"
)

func RenderSecondPage (c *gin.Context) {
	baseUrl := os.Getenv("BASE_URL")
	c.HTML(http.StatusOK, "page2.html", gin.H{
		"title": "Main website",
		"baseUrl": baseUrl,
	})
}

func RequestTransactionsGobiz (c *gin.Context) {
	type Body struct {
		MerchantID 		string 		`json:"merchant_id"`
		AccessToken 	string 		`json:"access_token"`
		StartDate 		string 		`json:"start_date"`
		EndDate 		string 		`json:"end_date"`
	}

	var payload Body
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Body request is not correct.",
		})
	}

	var customData []gobiz.GobizTransaction
	from := 0
	res, err := gobiz.RequestTransactions(payload.AccessToken, payload.MerchantID, payload.StartDate, payload.EndDate, from)
	total := res.Total
	ctr := 0
	for ctr < total {
		var wg sync.WaitGroup
		wg.Add(len(res.ResponseHits))
		for i, _ := range res.ResponseHits {
			go func(idx int) {
				ord, _ := gobiz.RequestOrder(payload.AccessToken, res.ResponseHits[idx].ResponseTransaction.OrderID)
				temp := gobiz.GobizCustomer{
					CustomerID: ord.ResponseData.ResponseOrder[0].ProductSpecific.GoResto.CustomerID,
					CustomerName: ord.ResponseData.ResponseOrder[0].ProductSpecific.GoResto.CustomerName,
					CustomerPhone: ord.ResponseData.ResponseOrder[0].ProductSpecific.GoResto.CustomerPhone,
					CustomerEmail: ord.ResponseData.ResponseOrder[0].ProductSpecific.GoResto.CustomerEmail,
				}

				new := gobiz.GobizTransaction{
					OrderID: res.ResponseHits[idx].ResponseTransaction.OrderID,
					Customer:  temp,
					OrderedAt: ord.ResponseData.ResponseOrder[0].OrderedAt,
				}
				// res.ResponseHits[idx].ResponseTransaction.Customer = temp
				// res.ResponseHits[idx].ResponseTransaction.OrderedAt = ord.ResponseData.ResponseOrder[0].OrderedAt
				customData = append(customData, new)
				defer wg.Done()
			}(i)
			ctr+=1
		}
		wg.Wait()
		if ctr < total {
			res, err = gobiz.RequestTransactions(payload.AccessToken, payload.MerchantID, payload.StartDate, payload.EndDate, ctr)
		}
	}


	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to request transactions merchant.",
			"error": err,
			"response_data": customData,
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "Success request transactions merchant.",
		"response_data": customData,
	})
	return
}