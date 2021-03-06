package host

import (
	"bufio"
	"io"
	"net"
	"regexp"
	"strings"

	"github.com/guumaster/hostctl/pkg/host/errors"
	"github.com/guumaster/hostctl/pkg/host/render"
)

var (
	profileNameRe = regexp.MustCompile(`# profile(?:.(on|off))?\s+([a-z0-9-_.\s]+)`)
	profileEnd    = regexp.MustCompile(`(?i)# end\s*`)
	spaceRemover  = regexp.MustCompile(`\s+`)
	tabReplacer   = regexp.MustCompile(`\t+`)
)

type Parser interface {
	Parse(reader io.Reader) Content
}

// Parse reads content from reader into Data struct.
func Parse(r io.Reader) (*Content, error) {
	data := &Content{
		ProfileNames: []string{},
		Profiles:     map[string]*Profile{},
	}

	currProfile := ""

	s := bufio.NewScanner(r)
	for s.Scan() {
		b := s.Bytes()

		switch {
		case profileNameRe.Match(b):
			p, _ := parseProfileHeader(b)

			currProfile = p.Name
			data.ProfileNames = append(data.ProfileNames, currProfile)
			data.Profiles[currProfile] = &Profile{
				Name:   currProfile,
				Status: p.Status,
			}

		case profileEnd.Match(b):
			currProfile = ""

		case currProfile != "":
			p := data.Profiles[currProfile]
			appendLine(p, string(b))
			data.Profiles[currProfile] = p

		default:
			row := parseToDefault(b, currProfile)
			data.DefaultProfile = append(data.DefaultProfile, row)
		}

		if err := s.Err(); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func appendLine(p *Profile, line string) {
	if line == "" {
		return
	}

	route, ok := ParseLine(line)
	if !ok {
		return
	}

	ip := route.IP.String()
	p.AddRoutes(ip, route.HostNames)
}

func parseToDefault(b []byte, currProfile string) *render.Row {
	var row *render.Row

	if len(b) == 0 {
		if currProfile == "" {
			row = &render.Row{Comment: ""}
		}

		return row
	}

	line, ok := ParseLine(string(b))
	if !ok {
		row = &render.Row{
			Comment: string(b),
		}
	} else {
		status := Enabled
		if off, _ := regexp.Match("^#", b); off {
			status = Disabled
		}
		row = &render.Row{
			Profile: Default,
			Status:  string(status),
			IP:      line.IP.String(),
			Host:    line.HostNames[0],
		}
	}

	return row
}

func parseProfileHeader(b []byte) (*Profile, error) {
	rs := profileNameRe.FindSubmatch(b)
	if len(rs) != 3 || string(rs[2]) == "" {
		return nil, errors.ErrInvalidProfileHeader
	}

	status := Enabled
	if string(rs[1]) == string(Disabled) {
		status = Disabled
	}

	return &Profile{
		Name:   strings.TrimSpace(string(rs[2])),
		Status: status,
	}, nil
}

// ParseLine checks if a line is a host line or a comment line.
func ParseLine(str string) (*Route, bool) {
	clean := spaceRemover.ReplaceAllString(str, " ")
	clean = tabReplacer.ReplaceAllString(clean, " ")
	clean = strings.TrimSpace(clean)

	p := strings.Split(clean, " ")

	i := 0
	if p[0] == "#" && len(p) > 1 {
		i = 1
	}

	ip := net.ParseIP(p[i])

	if ip == nil {
		return nil, false
	}

	return &Route{IP: ip, HostNames: p[i+1:]}, true
}
