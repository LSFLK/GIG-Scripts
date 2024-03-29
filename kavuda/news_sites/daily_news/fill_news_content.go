package daily_news

import (
	"GIG-Scripts/extended_models"
	"GIG-Scripts/kavuda/models"
	"github.com/lsflk/gig-sdk/libraries/clean_html"
)

func (d DailyNewsDecoder) FillNewsContent(newsItem extended_models.NewsArticle) (extended_models.NewsArticle, string, error) {
	return models.FillNewsContent(newsItem, ".node-main-content .content", "", clean_html.HtmlCleaner{
		Config: clean_html.Config{
			IgnoreElements: []string{"section"},
			IgnoreClasses:  []string{"field-name-field-articletags", "print-edition", "field-name-field-section"},
		}}, d)
}
