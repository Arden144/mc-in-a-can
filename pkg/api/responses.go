package api

type BuildResponse struct {
	ProjectId   string
	ProjectName string
	Version     string
	Build       int32
	Time        string
	Changes     []Change
	Downloads   map[string]Download
}

type Change struct {
	Commit  string
	Summary string
	Message string
}

type Download struct {
	Name   string
	Sha256 string
}

type ProjectResponse struct {
	ProjectId     string   `json:"project_id"`
	ProjectName   string   `json:"project_name"`
	VersionGroups []string `json:"version_groups"`
	Versions      []string `json:"versions"`
}

type ProjectsResponse struct {
	Projects []string
}

type VersionGroupBuild struct {
	Build     int32
	Time      string
	Changes   []Change
	Downloads map[string]Download
}

type VersionGroupBuildsResponse struct {
	ProjectId    string
	ProjectName  string
	VersionGroup string
	Versions     []string
	Builds       []VersionGroupBuild
}

type VersionGroupResponse struct {
	ProjectId    string
	ProjectName  string
	VersionGroup string
	Versions     []string
}

type VersionResponse struct {
	ProjectId   string
	ProjectName string
	Version     string
	Builds      []int32
}
