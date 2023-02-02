package model

type PageReq struct {
	Index int `json:"page" form:"index"`
	Size  int `json:"size" form:"size"`
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func (p *PageReq) Validate() {
	if p.Index < 1 {
		p.Index = 1
	}
	if p.Size < 1 {
		p.Size = 100000
	}
	// if p.PerPage < 1 {
	// 	p.PerPage = MaxInt
	// }
}
