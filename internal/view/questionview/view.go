package questionview

import (
	"lms-backend/internal/model"
)

type View struct {
	ID          uint    `json:"id,omitempty"`
	Description string  `json:"description"`
	Answer      string  `json:"answer"`
	Cost        float64 `json:"cost"`
	WorksheetID uint    `json:"worksheet_id"`
}

func ToView(question *model.Question) *View {
	return &View{
		ID:          question.ID,
		Description: question.Description,
		Answer:      question.Answer,
		Cost:        question.Cost,
		WorksheetID: question.WorksheetID,
	}
}
