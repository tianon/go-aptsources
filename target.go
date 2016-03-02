package aptsources

import (
	"strings"
)

var (
	DefaultTypes = []string{"deb", "deb-src"}
	DefaultDebianURIs = []string{"http://httpredir.debian.org/debian"}
)

func DebianSources(suite string, components ...string) Sources {
	suite = strings.TrimSuffix(suite, "-security")
	source := Source{
		Types:      DefaultTypes,
		URIs:       DefaultDebianURIs,
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
