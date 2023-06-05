package cli

// Flags stores data for the options and variables for a command.
type Flags struct {
	mapping map[string]flagMetadata
}

// flagsMetadata stores data for a single options and variable if applicable.
type flagMetadata struct {
	isset    bool
	hasVar   bool
	variable string
}

// Exists returns whether the given label exists in the Flags.
func (flags Flags) Exists(label string) bool {
	_, ok := flags.mapping[label]
	return ok
}

// IsSet returns whether the given label has been set in the command.
func (flags Flags) IsSet(label string) bool {
	meta, ok := flags.mapping[label]
	if !ok {
		return false
	}
	return meta.isset
}

// GetVar returns the variable set for the option with the given label.
// If the option has not been set, doesn't have a variable or doesn't
// exist in Flags, ("", false) will be returned instead.
func (flags Flags) GetVar(label string) (variable string, exists bool) {
	meta, ok := flags.mapping[label]
	if !ok || !meta.isset || !meta.hasVar {
		return "", false
	}
	return meta.variable, true
}
