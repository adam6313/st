package source

import (
	exifv3 "github.com/dsoprea/go-exif/v3"
	log "github.com/dsoprea/go-logging"

	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

type exif struct {
	// Make - ex: Canon
	Make string

	// Model - ex: Canon EOS 5D Mark III
	Model string

	// Orientation - ex: 1
	Orientation string

	// DateTime - ex: 2017:12:02 08:18:50
	DateTime string
}

func getEXIF(data []byte) (e *exif) {
	foundAt := -1
	e = new(exif)
	// check exif exist
	for i := 0; i < len(data); i++ {
		if _, err := exifv3.ParseExifHeader(data[i:]); err == nil {
			foundAt = i
			break
		} else if log.Is(err, exifv3.ErrNoExif) == false {
			return nil
		}
	}

	if foundAt == -1 {
		return
	}

	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return
	}

	// new index
	ti := exifv3.NewTagIndex()

	visitor := func(ite *exifv3.IfdTagEntry) (err error) {
		defer func() {
			if state := recover(); state != nil {
				return
			}
		}()

		tagID := ite.TagId()

		valueString, err := ite.FormatFirst()
		if err != nil {
			return err
		}

		set(e, tagID, valueString)

		return nil
	}

	_, _, err = exifv3.Visit(exifcommon.IfdStandardIfdIdentity, im, ti, data[foundAt:], visitor, nil)
	if err != nil {
		return
	}

	return
}

func set(e *exif, tagID uint16, value string) {
	switch tagID {
	case 271:
		e.Make = value
	case 272:
		e.Model = value
	case 274:
		e.Orientation = value
	case 306:
		e.DateTime = value
	}
}
