package usecase

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"storage/app/domain/model"
	"storage/app/domain/repository"
	"storage/app/domain/service"
	so "storage/app/domain/service/source"

	"github.com/sirupsen/logrus"
)

type storageUsecase struct {
	service    service.Service
	googleRepo repository.StorageRepository
	mongoRepo  repository.MongoRepository
	log        func(context.Context) *logrus.Entry
}

// NewStorageID - generate new storage id
func (s *storageUsecase) NewStorageID(ctx context.Context) string {
	return s.service.NewID()
}

// UploadFile -
func (s *storageUsecase) UploadFile(ctx context.Context, in *model.UploadFileRequest) (*model.File, error) {
	// get source
	source, err := s.service.NewSource(in)
	if err != nil {
		in.Data = nil
		s.log(ctx).WithField("NewSource failed in", in).Error()
		return &model.File{}, err
	}

	// set id
	source.SetID(s.NewStorageID(ctx))

	// set prefix
	source.SetVersionPrefix("v160")

	// protect file
	if err := source.Protect(in.Protect); err != nil {
		in.Data = nil
		s.log(ctx).WithField("Protect failed, in: ", in).Error()
		return &model.File{}, err
	}

	// devices
	devices := source.CompressDevice()

	width, _ := devices[so.ORIGIN]

	transfer, err := source.Compress(95, width, so.ORIGIN)
	if err != nil {
		return &model.File{}, err
	}

	// UploadFile
	f, err := s.service.UploadFile(ctx, transfer)
	if err != nil {
		//s.service.DeleteFile(ctx, fs)
		s.log(ctx).WithField("UploadFile (google), in: ", f).Error()
		return &model.File{}, err
	}

	if err := s.mongoRepo.Create(ctx, f); err != nil {
		s.log(ctx).WithField("CreateMany (mongo), in: ", f).Error()
		return &model.File{}, err
	}

	f.Metadata["IsResponsive"] = in.IsResponsive
	f.Metadata["width"] = in.Width
	f.Metadata["id"] = f.ID

	// process in background
	go s.uploadImageInbackground(context.Background(), f)

	return f, nil
}

// uploadImageInbackground -
func (s *storageUsecase) uploadImageInbackground(ctx context.Context, in *model.File) {
	s.log(ctx).Info("upload image in background", in)

	data, err := downloadFile(in.URL)
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

	source, err := s.service.NewSource(uploadFileRequest)
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
			e := fmt.Sprintf("source.Compress failed, device: %s,width %d in", k, width)
			s.log(ctx).WithField(e, uploadFileRequest).Error()
			return
		}

		// UploadFile
		f, err := s.service.UploadFile(ctx, transfer)
		if err != nil {
			e := fmt.Sprintf("UploadFile (google) failed, device: %s,width %d in", k, width)
			s.log(ctx).WithField(e, uploadFileRequest).Error()
			return
		}

		if err := s.mongoRepo.Create(ctx, f); err != nil {
			e := fmt.Sprintf("CreateMany (mongo) failed, device: %s,width %d in", k, width)
			s.log(ctx).WithField(e, uploadFileRequest).Error()
			return
		}
	}
}

func downloadFile(url string) ([]byte, error) {
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
func (s *storageUsecase) DeleteFile(ctx context.Context, in *model.DeleteFileRequest) error {
	if in.ID != "" {
		fs, err := s.mongoRepo.FindByID(ctx, in.ID)
		if err != nil {
			return err
		}

		for _, f := range fs {
			s.mongoRepo.DeleteByURL(ctx, f.URL)
			s.googleRepo.Delete(ctx, f.Name)
		}

	}

	return nil
}
