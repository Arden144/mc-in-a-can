package update

import (
	"errors"
	"os"

	"github.com/Arden144/paperupdate/pkg/hash"
	"github.com/Arden144/paperupdate/pkg/paper"
	"github.com/Arden144/paperupdate/pkg/req"
)

var ErrInvalidVersion = errors.New("not a valid version group")
var ErrLatest = errors.New("already up to date")

func TryUpdate(version, path string) (build int32, err error) {
	p := paper.Project("paper")

	vs, err := p.Versions()
	if err != nil {
		return
	}

	var v string
	if version == "latest" {
		v = vs[len(vs)-1]
	} else {
		var exists bool
		for _, v := range vs {
			if v == version {
				exists = true
				break
			}
		}
		if !exists {
			return 0, ErrInvalidVersion
		}
		v = version
	}

	d, err := p.DownloadInfo(v)
	if err != nil {
		return
	}

	if cs, err := hash.FromFile(path); err == nil {
		if cs == d.Sha256 {
			return 0, ErrLatest
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return 0, err
	}

	_, err = req.GetFile(d.Url, path)
	if err != nil {
		return
	}

	return d.Build, nil
}
