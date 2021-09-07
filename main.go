package main

import (
	"bytes"
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	"github.com/mailjet/mailjet-apiv3-go/v3"
	"github.com/sirupsen/logrus"
	"github.com/xhit/go-simple-mail/v2"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var banned []string


type Request struct {
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Url   string `db:"url"`
}

type Article struct {
	Id         int    `db:"id"`
	RequestId  int    `db:"request_id"`
	Hid        string `db:"hid"`
	Publisher  string `db:"publisher"`
	Location   string `db:"location"`
	Title      string `db:"title"`
	Href       string `db:"href"`
	Image      string `db:"image"`
	Features   string `db:"features"`
	Price      int    `db:"price"`
	LastUpdate string `db:"last_update"`
}
func init() {
	banned = append(banned, "37884-1")
}
func main() {
	defer holdUnexpectedError()
	cfg := LoadConfig()
	InitDB(cfg.DBConfig.GetDBURL())
	defer DB.Close()

	execute()
}

func execute() {
	requests := FindAllRequests()

	if requests != nil {
		for _, request := range requests {
			find(request)
		}
	}
}

func find(request Request) {
	newArticles := []Article{}
	// "https://www.habitaclia.com/alquiler-sabadell.htm?ordenar=mas_recientes&pmax=700&codzonas=2,5,31,8,10,11,12,13,40,19,20,32"
	res, err := http.Get(request.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error : %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// puede haber una segunda lista no que no sigue los filtros como sugerencia
	firstList := doc.Find(".list-items").First()
	articles := firstList.Find(".js-list-item")

	articles.Each(func(i int, s *goquery.Selection) {
		article := Article{}
		article.Publisher, _ = s.Attr("data-publisherid")
		article.RequestId = request.Id
		if isNotBanned(article.Publisher) {

			article.Hid, _ = s.Attr("data-id")
			article.Href, _ = s.Attr("data-href")
			item := s.Find(".list-item")

			article.Image, _ = item.Find(".list-gallery-image").First().Attr("data-image")
			article.Image = "https://" + strings.TrimLeft(article.Image, "//")
			info := item.Find(".list-item-info .list-item-content ")
			article.Title = info.Find(".list-item-title a").Text()
			article.Location = info.Find(".list-item-location span").Text()
			article.Features = info.Find(".list-item-feature").Text()
			article.LastUpdate = strings.Trim(info.Find(".list-item-date").Text(), " ")

			article.Features = clean(article.Features)

			price := strings.Split(item.Find(".list-item-info .list-item-content-second .list-item-price span").Text(), " ")[0]
			article.Price, _ = strconv.Atoi(price)
			
			// find in database
			oldArticle := FindArticle(article.RequestId, article.Hid)
			if oldArticle == nil || oldArticle.Price != article.Price {
				if oldArticle != nil {
					Update(article)
				} else {
					Insert(article)
				}
				newArticles = append(newArticles, article)
			}
		}

		return
	})
	
	if len(newArticles) > 0 {
		sendMail(request, newArticles)
	}
}

func sendMail(request Request, articles []Article) {

	mailjetClient := mailjet.NewMailjetClient("896a90d1eb4be1f4a6f1c3e0fc9c263e", "4aae95854b83ab4975ad47e3fe7a3aa7")


	htmlBody := `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>`+request.Name+`</title>
	</head>
	<body><div><table style="width: 100%;" role="presentation" border="1" width="100%" cellspacing="0">
						<tbody>`


	attachs := mailjet.InlinedAttachmentsV31{}

	// add inline
	for _, ar := range articles {
		htmlBody += `<tr><td style="padding: 5px; width: 350px;">`
		sp := strings.Split(ar.Image, ".")
		ext := "."+sp[len(sp)-1]
		file := downloadImage(ar, ext)

		attach := mailjet.InlinedAttachmentV31{
			AttachmentV31 :mailjet.AttachmentV31{
				ContentType: http.DetectContentType(file.Data),
				Filename:    file.Name,
				Base64Content: base64.StdEncoding.EncodeToString(file.Data),
			},
			ContentID: ar.Hid,
		}

		attachs = append(attachs, attach)

		htmlBody += `<p><img src="cid:`+ar.Hid+`" alt="image" width="350px"/></p></td>`
		htmlBody += `<td style="padding: 5px; width: 70%;">`
		htmlBody += `<h2 style="font-size: 20px; margin: 5px; font-family: Avenir;">`+ar.Title+`</h2>`
		htmlBody += `<p style="margin: 5px; font-size: 16px; line-height: 24px; font-family: Avenir;">`+ar.Location+`</p>`
		htmlBody += `<p style="margin: 5px; font-size: 16px; line-height: 24px; font-family: Avenir;">`+strconv.Itoa(ar.Price)+" €"+`</p>`
		htmlBody += `<p style="margin: 5px; font-size: 12px; line-height: 24px; font-family: Avenir;">`+ar.LastUpdate+`</p>`
		htmlBody += `<td><p style="margin: 5px; font-size: 14px; line-height: 24px; font-family: Avenir;">`+ar.Features+`</p></td>`
		htmlBody += `<p style="margin: 0; font-size: 16px; line-height: 24px; font-family: Avenir;"><a style="color: #ff7a59; text-decoration: underline;" href="`+ar.Href+`">Ver</a></p>`
		htmlBody += `</td></tr>`
	}

	htmlBody += "</tbody></table></div></body></html>"

	messagesInfo := []mailjet.InfoMessagesV31 {
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "eric.gamo@gmail.com",
				Name: "Eric",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31 {
					Email: request.Email,
					Name: "Recipient",
				},
			},
			Subject: request.Name,
			TextPart: "",
			HTMLPart: htmlBody,
			InlinedAttachments: &attachs,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo }
	_, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		logrus.Error(err)
	} else {
		log.Println("Email Sent")
	}
}

func isNotBanned(publisher string) bool {
	for _, p := range banned {
		if p == publisher {
			return false
		}
	}

	return true
}

func clean(s string) string {
	for i := 0; i < 10; i++ {
		s = strings.ReplaceAll(s, "  ", " ")
		s = strings.ReplaceAll(s, "\t", "")
		s = strings.ReplaceAll(s, "\n", "")
	}
	return s
}

func holdUnexpectedError() {
	if err := recover(); err != nil {
		logrus.Error(err)
	}
}

func downloadImage(article Article, ext string) *mail.File {
	//Get the response bytes from the url

	response, err := http.Get(article.Image)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil
	}

	file := &mail.File{Name: article.Hid+ext, Inline: true}
	//Write the bytes to the fiel
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)

	if err != nil {
		return nil
	}

	file.Data = buf.Bytes()

	return file
}
