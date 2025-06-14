package flags

import (
	"bytes"
	"flag"
	"strings"
)

func Usage(txt string, flags *flag.FlagSet) string {
	u := &Usager{
		Usage: txt,
		Flags: flags,
	}

	return u.String()
}

type Usager struct {
	Usage string
	Flags *flag.FlagSet
}

func (u *Usager) String() string {
	out := new(bytes.Buffer)
	out.WriteString(strings.TrimSpace(u.Usage))
	out.WriteString("\n")
	out.WriteString("\n")

	return strings.TrimRight(out.String(), "\n")
}
