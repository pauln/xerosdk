package accounting

import (
	"encoding/json"
	"net/http"

	"github.com/quickaco/xerosdk/helpers"
)

const (
	brandingThemeURL = "https://api.xero.com/api.xro/2.0/BrandingThemes"
)

//BrandingTheme applies structure and visuals to an invoice when printed or sent
type BrandingTheme struct {

	// Xero identifier
	BrandingThemeID string `json:"BrandingThemeID,omitempty" xml:"BrandingThemeID,omitempty"`

	// Name of branding theme
	Name string `json:"Name,omitempty" xml:"Name,omitempty"`

	// Integer – ranked order of branding theme. The default branding theme has a value of 0
	SortOrder float64 `json:"SortOrder,omitempty" xml:"SortOrder,omitempty"`

	// UTC timestamp of creation date of branding theme
	CreatedDateUTC string `json:"CreatedDateUTC,omitempty" xml:"CreatedDateUTC,omitempty"`
}

func unmarshalBrandingTheme(brandingThemeBytes []byte) ([]BrandingTheme, error) {
	response := struct {
		Themes []BrandingTheme `json:"BrandingThemes,omitempty"`
	}{}
	err := json.Unmarshal(brandingThemeBytes, &response)
	if err != nil {
		return nil, err
	}

	for n := len(response.Themes) - 1; n >= 0; n-- {
		response.Themes[n].CreatedDateUTC, err = helpers.DotNetJSONTimeToRFC3339(response.Themes[n].CreatedDateUTC, true)
		if err != nil {
			return nil, err
		}
	}
	return response.Themes, nil
}

// FindBrandingThemes will get all BrandingThemes.
func FindBrandingThemes(cl *http.Client) ([]BrandingTheme, error) {
	brandingThemeBytes, err := helpers.Find(cl, brandingThemeURL, nil, nil)
	if err != nil {
		return nil, err
	}

	return unmarshalBrandingTheme(brandingThemeBytes)
}
