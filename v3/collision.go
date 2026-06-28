package swagno3

import (
	"fmt"
	"sort"
	"strings"
)

// NameCollisionError is returned by ToJson (and panicked by MustToJson) when the
// HidePackageName option causes two or more distinct, package-qualified types to
// map to the same stripped name. The generated document would otherwise silently
// drop one of the colliding schemas, so generation fails instead.
type NameCollisionError struct {
	// Collisions maps each colliding stripped name to the sorted list of full
	// (package-qualified) type names that produced it.
	Collisions map[string][]string
}

func (e *NameCollisionError) Error() string {
	shortNames := make([]string, 0, len(e.Collisions))
	for name := range e.Collisions {
		shortNames = append(shortNames, name)
	}
	sort.Strings(shortNames)

	parts := make([]string, 0, len(shortNames))
	for _, name := range shortNames {
		parts = append(parts, fmt.Sprintf("%q maps to multiple types [%s]", name, strings.Join(e.Collisions[name], ", ")))
	}

	return fmt.Sprintf(
		"swagno: HidePackageName name collision: %s; rename one of the types or set HidePackageName to false",
		strings.Join(parts, "; "),
	)
}

// collisionError inspects the recorded stripped-name -> full-type-name sets and
// returns a *NameCollisionError if any stripped name was produced by more than one
// distinct type. It returns nil when there are no collisions.
func collisionError(definitionTypeNames map[string]map[string]struct{}) error {
	var collisions map[string][]string
	for shortName, fullNames := range definitionTypeNames {
		if len(fullNames) < 2 {
			continue
		}
		names := make([]string, 0, len(fullNames))
		for fullName := range fullNames {
			names = append(names, fullName)
		}
		sort.Strings(names)
		if collisions == nil {
			collisions = map[string][]string{}
		}
		collisions[shortName] = names
	}
	if collisions == nil {
		return nil
	}
	return &NameCollisionError{Collisions: collisions}
}
