package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"storage/app/domain/model"
	"storage/app/domain/service/source"

	so "storage/app/domain/service/source"
)

// NewID -
func (s *service) NewID() string {
	return s.node.Generate().String()
}

// UploadFile -
func (s *service) UploadFile(ctx context.Context, transfer *model.Transfer) (*model.File, error) {
	var err error

	// set object
	object := setObject(transfer.ID, transfer)

	if err = s.googleRepo.Upload(ctx, object, transfer.Data); err != nil {
		s.log(ctx).WithField("in :", object).Error()
		return nil, err
	}

	if err = s.googleRepo.ACL(ctx, object); err != nil {
		s.log(ctx).WithField("in :", object).Error()
		return nil, err

	}

	var f *model.File
	if f, err = s.googleRepo.Attrs(ctx, object); err != nil {
		s.log(ctx).WithField("in :", object).Error()
		return nil, err
	}

	s.SetFileInfo(f, transfer)

	return f, nil
}

// uploadImageInbackground -
func (s *service) UploadImageInbackground(ctx context.Context, in *model.File) {
	s.log(ctx).Info("upload image in background", in)

	data, err := s.DownloadFile(in.URL)
	if err != nil {
		return
	}

	uploadFileRequest := &model.UploadFileRequest{
		Data:   data,
		Prefix: in.Prefix,
	}

	// set name
	name, ok := in.Metadata["name"].(string)
	if ok {
		uploadFileRequest.Name = name
	} else {
		s.log(ctx).WithField("in", in).Error("interface conversion failed: conversion name to string failed")
		return
	}

	// set width
	width, ok := in.Metadata["width"].(int)
	if ok {
		uploadFileRequest.Width = width
	} else {
		s.log(ctx).WithField("in", in).Error("interface conversion failed: conversion width to int failed")
		return
	}

	// set isResponsive
	isResponsive, ok := in.Metadata["IsResponsive"].(bool)
	if ok {
		uploadFileRequest.IsResponsive = isResponsive
	} else {
		s.log(ctx).WithField("in", in).Error("interface conversion failed: conversion isResponsive to bool failed")
		return
	}

	if !uploadFileRequest.IsResponsive {
		return
	}

	source, err := s.NewSource(uploadFileRequest)
	if err != nil {
		uploadFileRequest.Data = nil
		s.log(ctx).WithField("NewSource failed in", uploadFileRequest).Error()
	}

	id, ok := in.Metadata["id"].(string)
	if ok {
		source.SetID(id)
	} else {
		s.log(ctx).WithField("in", uploadFileRequest).Error("interface conversion failed: conversion id to string failed")
		return
	}

	devices := source.CompressDevice()

	for k, width := range devices {
		if k == so.ORIGIN {
			continue
		}

		transfer, err := source.Compress(95, width, k)
		if err != nil {
			return
		}

		// UploadFile
		f, err := s.UploadFile(ctx, transfer)
		if err != nil {
			s.log(ctx).WithField("UploadFile (google), in: ", f).Error()
			return
		}

		if err := s.mongoRepo.Create(ctx, f); err != nil {
			s.log(ctx).WithField("CreateMany (mongo), in: ", f).Error()
			return
		}
	}
}

func (s *service) DownloadFile(url string) ([]byte, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Write the body to file
	return ioutil.ReadAll(resp.Body)
}

// DeleteFile -
func (s *service) DeleteFile(ctx context.Context, in []*model.File) error {
	if in == nil {
		return nil
	}

	for _, v := range in {
		s.googleRepo.Delete(ctx, v.Name)
	}

	return nil
}

// SetFileInfo -
func (s *service) SetFileInfo(f *model.File, t *model.Transfer) {
	// ID
	f.ID = t.ID
	f.Prefix = t.Prefix
	// set metadata
	f.Metadata = map[string]interface{}{
		"hash":   t.Hash,
		"device": t.Device,
		"name":   t.Name,
	}

	// custom file name and set url
	u := fmt.Sprintf("%s/%s/%s", s.googleRepo.GetDomain(), f.Bucket, f.Name)
	publicURL, _ := url.Parse(u)

	f.URL = publicURL.String()
	if t.Hash != "" {
		f.URL += "?hash=" + t.Hash
	}
}

// FindFile -
func (s *service) FindFile(ctx context.Context, fs []*model.File, target string) *model.File {
	for _, v := range fs {
		if val, ok := v.Metadata["device"]; ok && val == target {
			return v
		}
	}

	return nil
}

// SetObject -
func setObject(id string, in *model.Transfer) (object string) {
	var (
		target = "/"
		r      = []rune(in.Prefix)
		newR   = make([]rune, 0)
		start  = 0
		end    = len(r) - 1
	)

	for i, v := range r {
		switch i {
		case start, end:
			if string(v) != target {
				newR = append(newR, v)
			}
		default:
			newR = append(newR, v)
		}
	}

	object = string(newR)

	switch in.Device {
	case source.ORIGIN:
	default:
		object += "/" + in.Device
	}

	object += "/" + id + "." + in.Extension

	return object

}
