package routem

type (
	group struct {
		creator
		prefix string
	}
)

// =-=-=-=
// Getters
// =-=-=-=

func (g *group) Path() string {
	return g.prefix
}

// =-=-=-=
// Helpers
// =-=-=-=

func newGroup(defs config, prefix string) *group {
	g := &group{
		creator: newCreator(defs),
		prefix:  prefix,
	}

	return g
}
