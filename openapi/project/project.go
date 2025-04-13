package project

type Project struct {
	endPoint string
}

func NewDefaultProject() *Project {
	return &Project{
		endPoint: "https://api.dooray.com",
	}
}
func NewProject(endPoint string) *Project {
	return &Project{
		endPoint: endPoint,
	}
}
