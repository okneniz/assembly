package arm64

// Meta — descriptive instruction metadata.
type Meta struct {
	Name        string
	Aliases     []string
	Group       string
	Tags        []string
	Description string
	DocURL      string
}

// NewMeta — metadata by name and group; optional tags are passed variadically.
func NewMeta(name, group string, tags ...string) Meta {
	return Meta{
		Name:  name,
		Group: group,
		Tags:  tags,
	}
}
