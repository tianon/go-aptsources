package sources

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
	URIs       []string `delim:" "` // "http://httpredir.debian.org/debian"
	Suites     []string `delim:" "` // "jessie"
	Components []string `delim:" "` // "main", "contrib", "non-free"

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
				// TODO Architectures
				ret = append(ret, fmt.Sprintf("%s %s %s %s", t, uri, suite, strings.Join(source.Components, " ")))
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

func (sources Sources) ListString() string {
	ret := []string{}
	for _, s := range sources {
		if str := s.ListString(); str != "" {
			ret = append(ret, str)
		}
	}
	return strings.Join(ret, "\n\n")
}
