package aptsources // import "go.tianon.xyz/aptsources"

import (
	"fmt"
	"strings"

	"pault.ag/go/debian/control"
)

// http://manpages.debian.org/cgi-bin/man.cgi?manpath=experimental&query=sources.list
type Source struct {
	control.Paragraph

	// required
	Types      []string `delim:" "` // "deb", "deb-src"
	URIs       []string `delim:" "` // "http://deb.debian.org/debian"
	Suites     []string `delim:" "` // "jessie", "jsmith-unstable/"
	Components []string `delim:" "` // "main", "contrib", "non-free" (optional if suite ends in "/")

	// optional
	Architectures []string `delim:" "` // "amd64", "i386"
	// Languages
	// Targets
}

func (source Source) ListString() string {
	ret := []string{}
	for _, uri := range source.URIs {
		for _, suite := range source.Suites {
			for _, t := range source.Types {
				options := []string{}
				if len(source.Architectures) > 0 {
					options = append(options, fmt.Sprintf("arch=%s", strings.Join(source.Architectures, ",")))
				}
				parts := []string{t}
				if len(options) > 0 {
					parts = append(parts, fmt.Sprintf("[ %s ]", strings.Join(options, " ")))
				}
				parts = append(parts, uri, suite)
				if len(source.Components) > 0 {
					parts = append(parts, strings.Join(source.Components, " "))
				} // else if !strings.HasSuffix(suite, "/") { ERROR }
				ret = append(ret, strings.Join(parts, " "))
			}
		}
	}
	return strings.Join(ret, "\n")
}

type Sources []Source

func New(source ...Source) Sources {
	return Sources(source)
}

func (sources Sources) Append(source ...Source) Sources {
	return append(sources, source...)
}

func (sources Sources) Prepend(source ...Source) Sources {
	return append(source, sources...)
}

func (sources Sources) ListString() string {
	ret := []string{}
	for _, s := range sources {
		if str := s.ListString(); str != "" {
			ret = append(ret, str)
		}
	}
	return strings.Join(ret, "\n\n")
}
