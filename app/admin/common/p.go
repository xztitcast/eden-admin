package common

type P struct {
	Total    int64       `json:"total"`
	PageList interface{} `json:"pageList"`
}

type Index struct {
	PageNum  int `form:"pageNum" json:"pageNum" binding:"required,gte=1"`
	PageSize int `form:"pageSize" json:"pageSize" binding:"required,gte=10"`
}

func NewP(total int64, pageList interface{}) *P {
	return &P{
		Total:    total,
		PageList: pageList,
	}
}
