package serializers

type CreateTemplateReq struct {
	Type              string `json:"type" binding:"required"`
	Description       string `json:"description"`
	Message           string `json:"message"`
	EmailSubject      string `json:"email_subject"`
	EmailMarkup       string `json:"markup"`
	EmailRenderedHTML string `json:"email_rendered_html"`
	EventId           uint   `json:"event_id" binding:"required"`
	Active            bool   `json:"active"`
	Title             string `json:"title"`
	Link              string `json:"link"`
	File              string `json:"file"`
	Language          string `json:"language"`
}

type UpdateTemplateReq struct {
	Type              string `json:"type" binding:"required"`
	EventId           uint   `json:"event_id" binding:"required"`
	Description       string `json:"description"`
	Message           string `json:"message"`
	EmailSubject      string `json:"email_subject"`
	EmailMarkup       string `json:"markup"`
	EmailRenderedHTML string `json:"email_rendered_html"`
	Active            bool   `json:"active"`
	Title             string `json:"title"`
	Link              string `json:"link"`
	File              string `json:"file"`
	Language          string `json:"language" gorm:"default:'en'"`
}
