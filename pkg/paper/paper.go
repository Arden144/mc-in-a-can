package paper

import (
	"errors"
	"fmt"

	"github.com/Arden144/paperupdate/pkg/api"
	"github.com/Arden144/paperupdate/pkg/req"
)

type Project string
type Download struct {
	Url, Sha256 string
	Build       int32
}

func (p Project) Versions() ([]string, error) {
	var pInfo api.ProjectResponse
	err := req.GetJson(fmt.Sprintf(api.ProjectUrl, p), &pInfo)
	if err != nil {
		return nil, err
	}
	if len(pInfo.VersionGroups) == 0 {
		return nil, errors.New("no version groups found")
	}
	return pInfo.VersionGroups, nil
}

func (p Project) DownloadInfo(vg string) (*Download, error) {
	var vInfo api.VersionGroupBuildsResponse
	err := req.GetJson(fmt.Sprintf(api.GroupBuildsUrl, p, vg), &vInfo)
	if err != nil {
		return nil, err
	}
	g := vInfo.Builds[len(vInfo.Builds)-1]
	v := vInfo.Versions[len(vInfo.Versions)-1]
	b := g.Build
	d := g.Downloads["application"]
	return &Download{
		Url:    fmt.Sprintf(api.DownloadUrl, p, v, b, d.Name),
		Sha256: d.Sha256,
		Build:  b,
	}, nil
}
