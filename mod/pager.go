package mod

type Pager struct {
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
	Data  []interface{} `json:"data"`
}

func NewPager() *Pager {
	return &Pager{Page: 1, Limit: 20}
}

func (pager *Pager) GetStart() int {
	return (pager.GetPage() - 1) * pager.GetLimit()
}

func (pager *Pager) GetLimit() int {
	if pager.Limit < 1 || pager.Limit > 10000 {
		pager.Limit = 20
	}
	return pager.Limit
}

func (pager *Pager) GetPage() int {
	if pager.Page < 1 || pager.Page > 10000 {
		pager.Page = 1
	}
	return pager.Page
}
