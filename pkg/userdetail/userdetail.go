package userdetail

type UserDetail struct {
	Name string `json:"name" jsonschema:"title=the name,description="he name of the user"`
	Age  int    `json:"age" jsonschema:"title=the age, description="the age of the user"`
}
