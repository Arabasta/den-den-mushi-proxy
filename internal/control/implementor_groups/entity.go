package implementor_groups

type Entity struct {
	ImplementorGroupId string   `json:"implementor_group_id"`
	Members            []string `json:"members"`
}
