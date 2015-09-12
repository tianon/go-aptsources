package sources

import (
	"strings"
)

func SuiteSources(suite string, components ...string) Sources {
	source := Source{
		Types:      []string{"deb", "deb-src"},
		URIs:       []string{"http://httpredir.debian.org/debian"},
		Suites:     []string{suite},
		Components: components,
	}
	switch suite {
	case "experimental", "rc-buggy":
		source.Suites = append([]string{"sid"}, source.Suites...)
		fallthrough
	case "sid", "unstable":
		return New(source)
	}
	origSuite := suite
	for _, suffix := range []string{"backports", "lts"} {
		suffix = "-" + suffix
		if strings.HasSuffix(suite, suffix) {
			suite = suite[:len(suite)-len(suffix)]
			source.Suites = append([]string{suite}, source.Suites...)
		}
	}
	source.Suites = append(source.Suites, suite+"-updates")
	sources := New(source, Source{
		Types:      source.Types,
		URIs:       []string{"http://security.debian.org"},
		Suites:     []string{suite + "/updates"},
		Components: source.Components,
	})
	switch suite {
	case "squeeze":
		if origSuite != suite+"-lts" {
			sources = sources.Append(Source{
				Types:      source.Types,
				URIs:       source.URIs,
				Suites:     []string{suite + "-lts"},
				Components: source.Components,
			})
		}
	}
	return sources
}
