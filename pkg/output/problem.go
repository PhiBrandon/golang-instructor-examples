package output

type Problem struct {
	Id      int64  `json:"id" jsonschema:"title=the id,description=id of the problem"`
	Problem string `json:"problem" jsonschema:"title=percieved or real problem/obstable that the job poster could or is currently facing."`
}

type Problems struct {
	Problems []Problem `json:"problems" jsonschema:"title=percieved or real problem/obstable that the job poster could or is currently facing."`
}
