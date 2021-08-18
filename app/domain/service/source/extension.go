package source

import (
	"strings"

	"github.com/tyr-tech-team/hawk/status"
)

// Extension -
func (f *file) Extension() string {
	return f.extension
}

func (i *img) Extension() string {
	return i.extension
}

// ExtensionSupported -
func extensionSupported(extension string, extensions []string) error {
	if len(extensions) == 0 {
		return nil
	}

	for _, v := range extensions {
		if strings.Contains(strings.ToLower(v), extension) {
			return nil
		}
	}

	return status.InvalidParameter.SetServiceCode(status.ServiceStorage).WithDetail("檔案格式不符合").Err()
}
