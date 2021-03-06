package models

import (
	"GIG-SDK/libraries"
	"GIG-SDK/libraries/clean_html"
	"GIG-SDK/models"
	"GIG-SDK/request_handlers"
	"golang.org/x/net/html"
	"strings"
)

type IDecoder interface {
	ExtractNewsItems() ([]NewsItem, error)
	FillNewsContent(newsItem NewsItem) (NewsItem, string, error)
	GetSourceTitle() string
	GetDefaultImageUrl() string
}

func FillNewsContent(newsItem NewsItem, contentClass string, htmlCleaner clean_html.HtmlCleaner, decoder IDecoder) (NewsItem, string, error) {
	resp, err := request_handlers.GetRequest(newsItem.Link)
	if err != nil {
		return newsItem, "", err
	}

	newsDoc, err := libraries.HTMLStringToDoc(resp)
	if err != nil {
		return newsItem, "", err
	}

	newsSelection := newsDoc.Find(contentClass).First()
	newsHtml, err := newsSelection.Html()
	if err != nil {
		return newsItem, "", err
	}

	news, err := html.Parse(strings.NewReader(newsHtml))
	if err != nil {
		return newsItem, "", err
	}

	//clean html code by removing unwanted information
	var imageList []models.Upload
	newsItem.Content, _, imageList, newsItem.ImageURL = htmlCleaner.CleanHTML(newsItem.Link, news)

	return UploadImagesToServer(newsItem, imageList, decoder.GetDefaultImageUrl()), newsSelection.Text(), nil
}

func UploadImagesToServer(newsItem NewsItem, imageList []models.Upload, defaultImageUrl string) NewsItem {
	if newsItem.ImageURL == "" {
		var imageUploadClass models.Upload
		newsItem.ImageURL, imageUploadClass = clean_html.GenerateImagePath(defaultImageUrl, defaultImageUrl)
		imageList = append(imageList, imageUploadClass)
	}

	for _, image := range imageList {
		go func(payload models.Upload) {
			request_handlers.UploadImage(payload)
		}(image)
	}

	return newsItem
}
