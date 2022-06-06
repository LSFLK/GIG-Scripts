package data_models

import (
	"GIG-Scripts/extended_models"
	"github.com/lsflk/gig-sdk/enums/ValueType"
	"github.com/lsflk/gig-sdk/models"
)

type Country struct {
	extended_models.Location
}

func (c *Country) SetCountryId(countryId string, source string) *Country {
	c.SetAttribute("country_id", models.Value{
		ValueType:   ValueType.String,
		ValueString: countryId,
		Source:      source,
	})
	return c
}
