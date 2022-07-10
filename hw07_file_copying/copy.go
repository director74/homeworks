package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"math"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var buffSize int64
	var cutLimit int64
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()
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

	if limit > 0 && (offset+limit) < size {
		size = limit
	}

	buffSize = 512
	steps := int64(math.Ceil(float64(size-offset) / float64(buffSize)))
	bar := pb.Start64(steps)
	defer bar.Finish()

	buff := make([]byte, buffSize)
	for offset < size {
		read, readErr := fromFile.ReadAt(buff, offset)
		offset += int64(read)
		if read > 0 {
			cutLimit = int64(read)
			if limit > 0 {
				if limit <= cutLimit {
					cutLimit = limit
				} else {
					limit -= cutLimit
				}
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
