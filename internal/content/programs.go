package content

// StandardProgram is a named, read-only built-in exercise set. To offer another
// standard set, add an entry to StandardPrograms with its own Build function.
type StandardProgram struct {
	Key   string          // stable id, stored as the user's active program
	Name  string          // display name in the program selector
	Build func() Resolved // produces the program (regenerated on demand)
}

// StandardPrograms are the selectable built-in exercise sets. The first entry is
// the default (used when nothing else is selected).
var StandardPrograms = []StandardProgram{
	{Key: "sergey", Name: "Sergey", Build: ResolveBuiltin},
	// Add more standard sets here, e.g.:
	// {Key: "mobility", Name: "Mobility", Build: resolveMobility},
}

// StandardByKey returns the Resolved program for a standard key.
func StandardByKey(key string) (Resolved, bool) {
	for _, p := range StandardPrograms {
		if p.Key == key {
			return p.Build(), true
		}
	}
	return Resolved{}, false
}

// DefaultStandardKey is the key of the default standard program.
func DefaultStandardKey() string { return StandardPrograms[0].Key }
