package sources

import (
	"pault.ag/go/resolver"
)

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
				for _, component := range source.Components { // TODO add support for no components (suite/)
					for _, arch := range arches {
						if err := resolver.AppendBinaryIndex(can, uri, suite, component, arch); err != nil {
							return err
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
