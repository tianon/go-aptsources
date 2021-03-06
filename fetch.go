package aptsources // import "go.tianon.xyz/aptsources"

import (
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/xi2/xz"
	"pault.ag/go/resolver"
)

// TODO rewrite all this

var compressions = []string{".xz", ".bz2", ".gz", ""}

func fetchCandidates(can *resolver.Candidates, url string) error {
	for _, comp := range compressions {
		resp, err := http.Get(url + comp)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}
		defer resp.Body.Close()
		switch comp {
		case ".xz":
			reader, err := xz.NewReader(resp.Body, 0)
			if err != nil {
				return err
			}
			return can.AppendBinaryIndexReader(reader)
		case ".bz2":
			return can.AppendBinaryIndexReader(bzip2.NewReader(resp.Body))
		case ".gz":
			reader, err := gzip.NewReader(resp.Body)
			if err != nil {
				return err
			}
			defer reader.Close()
			return can.AppendBinaryIndexReader(reader)
		}
		return can.AppendBinaryIndexReader(resp.Body)
	}
	return fmt.Errorf("unable to find %s (tried all of %q)", url, compressions)
}

func (source Source) fetch(can *resolver.Candidates, arches ...string) error {
	if len(source.Architectures) > 0 {
		// if this source has explict architectures, ignore what's requested
		arches = source.Architectures
	}
	for _, t := range source.Types {
		if t != "deb" {
			continue
		}
		for _, uri := range source.URIs {
			for _, suite := range source.Suites {
				if len(source.Components) == 0 && suite[len(suite)-1] == '/' {
					if err := fetchCandidates(can, uri+"/"+suite+"Packages"); err != nil {
						return err
					}
				} else {
					for _, component := range source.Components {
						for _, arch := range arches {
							if err := fetchCandidates(can, fmt.Sprintf("%s/dists/%s/%s/binary-%s/Packages", uri, suite, component, arch)); err != nil {
								return err
							}
						}
					}
				}
			}
		}
		break
	}
	return nil
}

func (source Source) FetchCandidates(arches ...string) (*resolver.Candidates, error) {
	can := resolver.Candidates{}
	if err := source.fetch(&can, arches...); err != nil {
		return nil, err
	}
	return &can, nil
}

func (sources Sources) FetchCandidates(arches ...string) (*resolver.Candidates, error) {
	can := resolver.Candidates{}
	for _, source := range sources {
		if err := source.fetch(&can, arches...); err != nil {
			return nil, err
		}
	}
	return &can, nil
}
