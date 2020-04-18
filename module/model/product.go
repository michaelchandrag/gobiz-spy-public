package model

import (
	"fmt"
	"time"

	db "github.com/michaelchandrag/go-my-skeleton/database"
)

type (
	Product struct {
		ID 					int 			`db:"id" json:"id"`
		ProductID 			string 			`db:"product_id" json:"product_id"`
		Link				string 			`db:"link" json:"link"`
		DetailLatest 		Detail 			`json:"detail_latest,omitempty"`
		DetailHistory 		[]Detail 			`json:"detail_history,omitempty"`
		CreatedAt			string 			`db:"created_at" json:"created_at"`
		UpdatedAt			string 			`db:"updated_at" json:"updated_at"`
		DeletedAt			string  		`db:"deleted_at" json:"deleted_at"`
	}
)

func (this *Product) Finds() (results []Product, err error) {
	query := `
		SELECT
			p.id,
			p.product_id,
			p.link,
			p.created_at,
			p.updated_at,
			COALESCE(p.deleted_at,"") as deleted_at
		FROM
			product as p
		`
	err = db.Engine.Select(&results, query)
	if err != nil {
		return nil, err
	}
	return results, nil
}


func (this *Product) FindByProductID(productId string) error {
	query := `
		SELECT
			id,
			product_id,
			link,
			created_at,
			updated_at,
			COALESCE(deleted_at,"") as deleted_at
		FROM
			product
		WHERE
			product_id = ?`
	if err := db.Engine.Get(this, query, productId); err != nil {
		return err
	}

	return nil
}

func (this *Product) Create(data Product) (result Product, err error) {
	currentTime := time.Now()

	data.CreatedAt = currentTime.Format("2006-01-02 15:04:05")
	data.UpdatedAt = currentTime.Format("2006-01-02 15:04:05")
	query := fmt.Sprintf(`
			INSERT INTO product (
				product_id, link, 
				created_at, updated_at, deleted_at
			) VALUES (
				?, ?, 
				?, ?, null
			)
		`)
	resp, err := db.Engine.Exec(query,
			data.ProductID, data.Link,
			data.CreatedAt, data.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	lastID, _ := resp.LastInsertId()
	result.ID = int(lastID)
	result.ProductID = data.ProductID
	result.Link = data.Link
	result.CreatedAt = data.CreatedAt
	result.UpdatedAt = data.UpdatedAt
	return result, nil
}