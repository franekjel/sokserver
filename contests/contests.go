package contests

import (
	"github.com/franekjel/sokserver/fs"
	"github.com/franekjel/sokserver/rounds"
)

//Contest holds rounds in contest and groups allowed to participate
type Contest struct {
	fs     *fs.Fs
	round  map[string]rounds.Round
	groups []string
}
