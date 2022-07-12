package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"math"
	"os"
)

var (
	ErrUnsupportedFile             = errors.New("unsupported file")
	ErrOffsetExceedsFileSize       = errors.New("offset exceeds file size")
	copyBuffSize             int64 = 512
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var cutLimit int64
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	fileInfo, statErr := fromFile.Stat()
	if statErr != nil {
		return statErr
	}
	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	size := fileInfo.Size()
	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	if limit > 0 && (offset+limit) < size {
		size = limit
	}

	steps := int64(math.Ceil(float64(size-offset) / float64(copyBuffSize)))
	bar := pb.Start64(steps)
	defer bar.Finish()

	buff := make([]byte, copyBuffSize)
	for offset < size {
		read, readErr := fromFile.ReadAt(buff, offset)
		cutLimit = int64(read)
		offset += cutLimit
		if cutLimit > 0 {
			if limit > 0 && limit <= cutLimit {
				cutLimit = limit
			} else if limit > 0 && limit > cutLimit {
				limit -= cutLimit
			}

			_, writeErr := toFile.Write(buff[:cutLimit])
			if writeErr != nil {
				return writeErr
			}
			bar.Increment()
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	return nil
}
