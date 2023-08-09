package worksheetview

import (
	"technical-test/internal/model"
	"technical-test/internal/view/questionview"
)

type WorkSheetView struct {
	WorkSheetListView

	Questions []questionview.View `json:"questions"`
}

func ToView(workSheet *model.Worksheet) *WorkSheetView {
	questionViews := make([]questionview.View, len(workSheet.Questions))
	for i, question := range workSheet.Questions {
		//nolint:gosec // loop does not modify struct
		questionViews[i] = *questionview.ToView(&question)
	}

	return &WorkSheetView{
		WorkSheetListView: *ToListView(workSheet),
		Questions:         questionViews,
	}
}
