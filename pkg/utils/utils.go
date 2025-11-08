package utils

import (
	"context"
	"errors"
	"os"

	"github.com/dhowden/tag"
)

func Ptr[T any](t T) *T {
	return &t
}

func SaveID3Image(ctx context.Context, img *tag.Picture, filename string) error {
	if img == nil {
		return errors.New("no image data in ID3 tag")
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(img.Data)

	return nil
}
