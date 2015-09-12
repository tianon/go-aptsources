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
	switch {
	case suite == "experimental", suite == "rc-buggy":
		source.Suites = append([]string{"sid"}, source.Suites...)
		fallthrough
	case suite == "sid", suite == "unstable":
		return New(source)
	case strings.HasSuffix(suite, "-backports"):
		suite = suite[:len(suite)-len("-backports")]
		source.Suites = append([]string{suite}, source.Suites...)
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
		sources = sources.Append(Source{
			Types:      source.Types,
			URIs:       source.URIs,
			Suites:     []string{suite + "-lts"},
			Components: source.Components,
		})
	}
	return sources
}
