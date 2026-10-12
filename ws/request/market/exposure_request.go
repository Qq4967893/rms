package market

type GetExposureListRequest struct {
	Params ExposureListRequestParams `json:"params" validate:"required"`
}

type ExposureListRequestParams struct {
	Books []string `json:"books" validate:"required"`
	Sort  string   `json:"sort" validate:"required,oneof=exposure exposure_USD currency"`
	Order string   `json:"order" validate:"required,oneof=asc desc"`
}
