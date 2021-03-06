package file

import (
	"github.com/guumaster/hostctl/pkg/host"
	"github.com/guumaster/hostctl/pkg/host/errors"
)

// RemoveProfiles removes given profiles from the list
func (f *File) RemoveProfiles(profiles []string) error {
	for _, p := range profiles {
		err := f.RemoveProfile(p)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveProfile removes given profile from the list
func (f *File) RemoveProfile(name string) error {
	var names []string

	if name == host.Default {
		return errors.ErrDefaultProfile
	}

	_, ok := f.data.Profiles[name]
	if !ok {
		return errors.ErrUnknownProfile
	}

	delete(f.data.Profiles, name)

	for _, n := range f.data.ProfileNames {
		if n != name {
			names = append(names, n)
		}
	}

	f.data.ProfileNames = names

	return nil
}
