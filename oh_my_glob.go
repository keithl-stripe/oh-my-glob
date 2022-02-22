package oh_my_glob

import (
	"log"
	"strings"
)

type part struct {
	stars int8 // 0, 1, or 2
	lit string
}

var starstar part = part {
	stars: 2,
	lit: "",
}

var star part = part {
	stars: 1,
	lit: "",
}

type Glob struct {
	// this is for pretty-printing the glob
	original string
	// we pre-break the glob at `/` boundaries for less processing
	// later
	parts []part
}

func Compile(glob string) Glob {
	var parts []part
	for _, fragment := range strings.Split(glob, "/") {
		if fragment == "**" {
			parts = append(parts, starstar)
		} else if fragment == "*" {
			parts = append(parts, star)
		} else {
			parts = append(parts, part {
				stars: 0,
				lit: fragment,
			})
		}
	}
	return Glob{
		original: glob,
		parts:    parts,
	}
}

// this code is borrowed directly from Russ Cox's research page,
// albeit with the single-character wildcard removed (since I don't
// believe we need it.) For more details:
//
//   https://research.swtch.com/glob
func match(pattern, name string) bool {
	px := 0
	nx := 0
	nextPx := 0
	nextNx := 0
	for px < len(pattern) || nx < len(name) {
		if px < len(pattern) {
			c := pattern[px]
			switch c {
			default:
				if nx < len(name) && name[nx] == c {
					px++
					nx++
					continue
				}

			case '*':
				// zero-or-more-character wildcard
				// Try to match at nx.
				// If that doesn't work out,
				// restart at nx+1 next.
				nextPx = px
				nextNx = nx + 1
				px++
				continue
			}
		}
		// Mismatch. Maybe restart.
		if 0 < nextNx && nextNx <= len(name) {
			px = nextPx
			nx = nextNx
			continue
		}
		return false
	}
	// Matched all of pattern to all of name. Success.
	return true
}

func (g *Glob) Match(path string) bool {
	chunks := strings.Split(path, "/")
	px := 0
	nx := 0
	nextPx := 0
	nextNx := 0

	for px < len(g.parts) || nx < len(chunks) {
		if px < len(g.parts) {
			c := g.parts[px]
			switch c.stars {
			case 0:
				if nx < len(chunks) && match(c.lit, chunks[nx]) {
					px++
					nx++
					continue
				}
			case 2:
				nextPx = px
				nextNx = nx + 1
				px++
				continue
			case 1:
				if nx < len(chunks) {
					px++
					nx++
					continue
				}
			default:
				// this should never happen and
				// indicates a bug in library code
				log.Fatalf("Unexpected compiled glob value")

			}
		}

		if 0 < nextNx && nextNx <= len(chunks) {
			px = nextPx
			nx = nextNx
			continue
		}
		return false
	}
	return true
}
