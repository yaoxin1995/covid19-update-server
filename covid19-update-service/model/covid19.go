package model

type Incidence struct {
	CommonModelFields
	Cases7Per100k float64 `json:"cases7_per_100k"`
}
