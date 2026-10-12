package model

type Paging struct {
	Limit  int
	Offset int
	Sort   string
	Order  string
}

type Option func(p *Paging)

func WithSort(sort string) Option {
	return func(p *Paging) {
		p.Sort = sort
	}
}

func WithOrder(order string) Option {
	return func(p *Paging) {
		p.Order = order
	}
}

func NewPaging(limit int, offset int, opts ...Option) (p Paging) {
	page := &Paging{
		Limit:  limit,
		Offset: offset,
		Sort:   "id",
		Order:  "desc",
	}

	for _, o := range opts {
		o(page)
	}
	return *page
}
