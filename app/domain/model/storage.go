package model

// File -
type File struct {
	// ID -
	ID string `bson:"id"`

	// Bucket -
	Bucket string `bson:"bucket"`

	// Name -
	Name string `bson:"name"`

	// URL -
	URL string `bson:"url"`

	// Prefix
	Prefix string `bson:"prefix"`

	// ContentType -
	ContentType string `bson:"contentType"`

	// CreatedAt -
	CreatedAt int64 `bson:"createdAt"`

	// Metadata -
	Metadata map[string]interface{} `bson:"metadata"`

	// Size -
	Size int64 `bson:"size"`
}

// Transfer -
type Transfer struct {
	// ID -
	ID string

	// Name -
	Name string

	// Prefix -
	Prefix string

	// Extension -
	Extension string

	// Device -
	Device string

	// Hash
	Hash string

	// Data -
	Data []byte
}

// UploadFileRequest -
type UploadFileRequest struct {
	// Data -
	Data []byte

	// Prefix -
	Prefix string

	// Name -
	Name string

	// Extensions - supported extensions
	Extensions []string

	// Protect -
	Protect bool

	// Width -
	Width int

	// Height -
	Height int

	// IsResponsive -
	IsResponsive bool
}

// DeleteFileRequest -
type DeleteFileRequest struct {
	// ID -
	ID string

	// Prefix -
	Prefix string

	// URL -
	URL string
}
