package cancelablefw

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

const (
	defaultBufSize = 4096
)

// CancelableFileWriter writes data to a file in a buffered manner with
// integrated cancellation. Will return an error in case that the file is
// incompletely written. It is the user's responsibility to clean up a
// potentially incomplete file
func CancelableFileWriter(ctx context.Context, data []byte, file *os.File) (nWrittenBytes int, err error) {
	return CancelableFileWriterSize(ctx, data, file, defaultBufSize)
}

// CancelableFileWriterCustomBufferLen writes data to a file in a buffered
// manner with integrated cancellation with a custom buffer length. Will
// return an error in case that the file is incompletely written. It is the
// user's responsibility to clean up a potentially incomplete file
func CancelableFileWriterSize(ctx context.Context, data []byte, file *os.File,
	bufferLen int) (nWrittenBytes int, err error) {
	bufferedWriter := bufio.NewWriterSize(file, bufferLen)

	for {
		select {
		case <-ctx.Done():
			return nWrittenBytes, fmt.Errorf("all data was not written to the file due to context cancellation")
		default:
		}

		var nn int

		if (nWrittenBytes + bufferLen) < len(data) {
			nn, err = bufferedWriter.Write(data[nWrittenBytes : nWrittenBytes+bufferLen])
		} else {
			nn, err = bufferedWriter.Write(data[nWrittenBytes : nWrittenBytes+(len(data)-nWrittenBytes)])
		}
		if err != nil {
			return nWrittenBytes, err
		}
		nWrittenBytes += nn
		err = bufferedWriter.Flush()
		if err != nil {
			return nWrittenBytes, err
		}

		if nWrittenBytes == len(data) {
			break
		}
	}
	return nWrittenBytes, nil
}
