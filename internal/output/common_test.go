package output

import (
	"github.com/barasher/dep-carto/internal/model"
)

func generateServers() []model.Server {
	s1a := model.Server{
		Hostname: "s1.domain",
		Key:      "a",
		IPs:      []string{"ip1a", "ip1b"},
		Dependencies: []model.Dependency{
			{Resource: "ip2"},
			{Resource: "s.otherdomain"},
		},
	}
	s1b := model.Server{
		Hostname: "s1.domain",
		Key:      "b",
		IPs:      []string{"ip1a", "ip1b"},
		Dependencies: []model.Dependency{
			{Resource: "s3.domain"},
		},
	}
	s2 := model.Server{
		Hostname: "s2.domain",
		IPs:      []string{"ip2"},
	}
	s3 := model.Server{
		Hostname: "s3.domain",
		IPs:      []string{"ip3"},
	}
	return []model.Server{s1a, s1b, s2, s3}
}
