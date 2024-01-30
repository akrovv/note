package service

type CreateNote struct {
	Text string `json:"text"`
}

type UpdateNote struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type DeleteNote struct {
	ID string `json:"id"`
}

type GetNote struct {
	ID string `json:"id"`
}

type GetNotes struct {
	OrderBy string
}
