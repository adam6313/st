package source

// SetID -
func (i *img) SetID(id string) {
	i.id = id
}

// SetID -
func (f *file) SetID(id string) {
	f.id = id
}

// SetVersionPrefix -
func (i *img) SetVersionPrefix(prefix string) {
	if prefix != "" {
		i.prefix = prefix + "/" + i.prefix
	}
}

// SetVersionPrefix -
func (f *file) SetVersionPrefix(prefix string) {
	if prefix != "" {
		f.prefix = prefix + "/" + f.prefix
	}
}
