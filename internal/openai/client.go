package openai

type Model struct {
	Id         string       `json:"id"`
	Object     string       `json:"object"`
	Created    int          `json:"created"`
	OwnedBy    string       `json:"owned_by"`
	Permission []Permission `json:"permission"`
}

type Permission struct {
	Id                 string   `json:"id"`
	Object             string   `json:"object"`
	Created            int      `json:"created"`
	AllowCreateEngine  bool     `json:"allow_create_engine"`
	AllowSampling      bool     `json:"allow_sampling"`
	AllowLogprobs      bool     `json:"allow_logprobs"`
	AllowSearchIndices bool     `json:"allow_search_indices"`
	AllowView          bool     `json:"allow_view"`
	AllowFineTuning    bool     `json:"allow_fine_tuning"`
	Organization       string   `json:"organization"`
	Group              []string `json:"group"`
	IsBlocking         bool     `json:"is_blocking"`
}

func NewModel(id string, owned_by string) *Model {
	permission := Permission{
		Id:                 "",
		Object:             "permission",
		Created:            1737331200,
		AllowCreateEngine:  false,
		AllowSampling:      true,
		AllowLogprobs:      true,
		AllowSearchIndices: false,
		AllowView:          true,
		AllowFineTuning:    false,
		Organization:       "*",
		Group:              nil,
		IsBlocking:         false,
	}
	return &Model{
		Id:         id,
		Object:     "model",
		Created:    1737331200,
		OwnedBy:    owned_by,
		Permission: []Permission{permission},
	}
}
