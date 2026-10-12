package response

type SummaryExposure struct {
	Currency    string `json:"currency"`
	Exposure    string `json:"exposure"`
	ExposureUSD string `json:"exposure_USD"`
}
