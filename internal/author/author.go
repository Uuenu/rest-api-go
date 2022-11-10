package author

type Author struct {
	ID   string `bson:"id"`
	Name string `json:"name"`
}
