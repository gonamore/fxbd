package reporters

// Interface for creating  downloaded account statistics
type Reporter interface {
	Assemble()
}
