package model

import (
	"fmt"
	"time"

	db "github.com/michaelchandrag/go-my-skeleton/database"
)

type (
	Detail struct {
		ID 					int 			`db:"id" json:"id"`
		ProductID 			string 			`db:"product_id" json:"product_id"`
		Title				string 			`db:"title" json:"title"`
		Description 		string 			`db:"description" json:"description"`
		Price 				int 			`db:"price" json:"price"`
		Images 				string 			`db:"images" json:"images"`
		CreatedAt			string 			`db:"created_at" json:"created_at"`
		UpdatedAt			string 			`db:"updated_at" json:"updated_at"`
		DeletedAt			string  		`db:"deleted_at" json:"deleted_at"`
	}
)

func (this *Detail) FindByProductID(productId string) (results []Detail, err error) {
	query := `
		SELECT
			d.id,
			d.product_id,
			d.title,
			d.description,
			d.price,
			d.images,
			d.created_at,
			d.updated_at,
			COALESCE(d.deleted_at,"") as deleted_at
		FROM
			detail as d
		WHERE
			product_id = ?
		ORDER BY
			d.created_at ASC
		`
	err = db.Engine.Select(&results, query, productId)
	if err != nil {
		return nil, err
	}
	return results, nil
}


func (this *Detail) FindLatestDetailByProductID(productId string) error {
	query := `
		SELECT
			d.id,
			d.product_id,
			d.title,
			d.description,
			d.price,
			d.images,
			d.created_at,
			d.updated_at,
			COALESCE(d.deleted_at,"") as deleted_at
		FROM
			detail as d
		WHERE
			product_id = ?
		ORDER BY
			d.created_at DESC
		LIMIT 1`
	if err := db.Engine.Get(this, query, productId); err != nil {
		return err
	}

	return nil
}

func (this *Detail) Create(data Detail) (result Detail, err error) {
	currentTime := time.Now()

	data.CreatedAt = currentTime.Format("2006-01-02 15:04:05")
	data.UpdatedAt = currentTime.Format("2006-01-02 15:04:05")
	query := fmt.Sprintf(`
			INSERT INTO detail (
				product_id, title, description, price, images, 
				created_at, updated_at, deleted_at
			) VALUES (
				?, ?, ?, ?, ?,
				?, ?, null
			)
		`)
	resp, err := db.Engine.Exec(query,
			data.ProductID, data.Title, data.Description, data.Price, data.Images,
			data.CreatedAt, data.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		return result, err
	}

	lastID, _ := resp.LastInsertId()
	result.ID = int(lastID)
	result.ProductID = data.ProductID
	result.Title = data.Title
	result.Description = data.Description
	result.Price = data.Price
	result.Images = data.Images
	result.CreatedAt = data.CreatedAt
	result.UpdatedAt = data.UpdatedAt
	return result, nil
}