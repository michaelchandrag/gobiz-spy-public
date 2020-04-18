package first

import (
	"fmt"
	"os"
	"net/http"
	"log"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/PuerkitoBio/goquery"


	model "github.com/michaelchandrag/go-my-skeleton/module/model"
)

type (
	MagentoGallery struct {
		DataGallery 	DataGallery 		`json:"[data-gallery-role=gallery-placeholder]"`
		
	}

	DataGallery struct {
		MageGallery 	MageGallery 		`json:"mage/gallery/gallery"`
	}

	MageGallery struct {
		Mixins 			[]string 			`json:"mixins"`
		Data 			[]Data 				`json:"data"`
	}

	Data struct {
		Img 			string 				`json:"img"`
	}

)

func RenderFirstPage (c *gin.Context) {
	baseUrl := os.Getenv("BASE_URL")
	c.HTML(http.StatusOK, "page1.html", gin.H{
		"title": "Main website",
		"baseUrl": baseUrl,
	})
}

func FetchDataFromFabelio (c *gin.Context) {
	type Body struct {
		Link 	string 	`json:"link"`
	}
	var payload Body
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Body request is not correct.",
		})
		return
	}
	res, err := http.Get(payload.Link)
	if err != nil {
	    log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
	    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	    c.JSON(400, gin.H{
	    	"success": false,
	    	"message": "Load fabelio product failed.",
	    })
	    return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		c.JSON(400, gin.H{
	    	"success": false,
	    	"message": "Load fabelio product failed.",
	    })
	    return
	}

	var title string
	var description string
	var price string
	var find string
	title = doc.Find(".base").Text()
	description = doc.Find("#description").Text()
	productId, _ := doc.Find("#productId").Attr("value")
	newProductId, _ := strconv.Atoi(productId)
	price, _ = doc.Find(fmt.Sprintf("#product-price-%s", productId)).Attr("data-price-amount")
	convPrice, _ := strconv.Atoi(price)
	doc.Find("script").Each( func (i int, s*goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "mage/gallery/gallery") {
			find = text
		}
	})
	if newProductId <= 0 {
		c.JSON(400, gin.H{
	    	"success": false,
	    	"message": "Load fabelio product failed.",
	    })
	    return
	}
	var magentoGallery MagentoGallery
	_ = json.Unmarshal([]byte(find), &magentoGallery)

	var whereProduct model.Product
	err = whereProduct.FindByProductID(productId) // err means product not found
	
	if err != nil {
		var newProduct model.Product
		object := model.Product{
			ProductID: productId,
			Link: payload.Link,
		}
		_, err := newProduct.Create(object)
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"message": "Something error when create product.",
			})
			return
		}
	}

	var newDetail model.Detail
	imagesToJSON, _ := json.Marshal(magentoGallery.DataGallery.MageGallery.Data)
	object := model.Detail{
		ProductID: productId,
		Title: title,
		Description: description,
		Price: convPrice,
		Images: string(imagesToJSON),
	}
	resultDetail, err := newDetail.Create(object)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Something error when create detail",
		})
		return
	}


	
	c.JSON(200, gin.H{
		"success": true,
		"detail": resultDetail,
		"message": "Success fetch fabelio product.",
	})
	return
}