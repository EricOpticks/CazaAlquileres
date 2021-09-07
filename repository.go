package main

import (
	"CazaAlquileres/database"
	"github.com/sirupsen/logrus"
)

const findallrequestsquery = `SELECT 
       	id, name, email, url
		FROM requests r`

func FindAllRequests() []Request{
	var requests []Request
	err := database.DB.Select(&requests, findallrequestsquery)

	if err != nil {
		logrus.Panic("Error retrieving requests: %s", err)
	}

	return requests
}

const findarticlequery = `SELECT 
       	id,hid,request_id,publisher,location,title,href,image,features,price,last_update
		FROM articles a 
		WHERE request_id = ? and hid = ?`

func FindArticle(requestId int, hid string) *Article {
	var articles []Article
	err := database.DB.Select(&articles, findarticlequery, requestId, hid)

	if err != nil {
		logrus.Panic("Error retrieving article: %s", err)
	}

	if articles != nil && len(articles) > 0 {
		return &articles[0]
	}

	return nil
}

const insertquery =  `
		INSERT INTO articles (hid,request_id,publisher,location,title,href,image,features,price,last_update)
		VALUES (:hid,:request_id, :publisher, :location, :title, :href, :image, :features, 
		        :price, :last_update)`
func Insert(article Article) {

	_, err := database.DB.NamedExec(insertquery, &article)
	if err != nil {
		logrus.Errorf("Error inserting article ", err)
	}
}

const updatequery =  `
		UPDATE articles (hid,request_id,publisher,location,title,href,image,features,price,last_update)
		VALUES (:hid,:request_id, :publisher, :location, :title, :href, :image, :features, 
		        :price, :last_update) where id = :id`
func Update(article Article) {

	_, err := database.DB.NamedExec(insertquery, &article)
	if err != nil {
		logrus.Errorf("Error inserting article ", err)
	}
}