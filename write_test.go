package cancelablefw

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type dummyData struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

func TestWriteFile(t *testing.T) {

	dummy := dummyData{
		Name:   "John Doe",
		Number: 4,
	}
	dummyData, err := json.Marshal(dummy)
	require.NoError(t, err)

	t.Run("Simple test with default buffer", func(t *testing.T) {
		tmpFile, err := ioutil.TempFile("", "cancelable_file_writer")
		require.NoError(t, err)
		defer func() {
			err := os.Remove(tmpFile.Name())
			require.NoError(t, err)
		}()

		nWrittenBytes, err := FileWriteWithContext(context.TODO(), dummyData, tmpFile)
		require.NoError(t, err)
		require.Equal(t, len(dummyData), nWrittenBytes)

		tmpFileContent, err := ioutil.ReadFile(tmpFile.Name())
		require.NoError(t, err)
		require.Equal(t, dummyData, tmpFileContent)
	})

	t.Run("Simple test with custom buffer", func(t *testing.T) {
		tmpFile, err := ioutil.TempFile("", "cancelable_file_writer")
		require.NoError(t, err)
		defer func() {
			err := os.Remove(tmpFile.Name())
			require.NoError(t, err)
		}()

		nWrittenBytes, err := FileWriteWithContextSize(context.TODO(), dummyData, tmpFile, 2)
		require.NoError(t, err)
		require.Equal(t, len(dummyData), nWrittenBytes)

		tmpFileContent, err := ioutil.ReadFile(tmpFile.Name())
		require.NoError(t, err)
		require.Equal(t, dummyData, tmpFileContent)
	})

	t.Run("Simple test with cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		tmpFile, err := ioutil.TempFile("", "cancelable_file_writer")
		require.NoError(t, err)
		defer func() {
			err := os.Remove(tmpFile.Name())
			require.NoError(t, err)
		}()

		cancel()
		nWrittenBytes, err := FileWriteWithContext(ctx, dummyData, tmpFile)
		require.Error(t, err)
		require.Equal(t, "all data was not written to the file due to context cancellation", err.Error())
		require.Equal(t, 0, nWrittenBytes)

		tmpFileContent, err := ioutil.ReadFile(tmpFile.Name())
		require.NoError(t, err)
		require.NotEqual(t, dummyData, tmpFileContent)
	})

	t.Run("Simple test with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		tmpFile, err := ioutil.TempFile("", "cancelable_file_writer")
		require.NoError(t, err)
		defer func() {
			err := os.Remove(tmpFile.Name())
			require.NoError(t, err)
		}()

		dummyData := make([]byte, 100000000)
		nWrittenBytes, err := FileWriteWithContextSize(ctx, dummyData, tmpFile, 1)
		require.Error(t, err)
		require.Equal(t, "all data was not written to the file due to context cancellation", err.Error())
		require.NotEqual(t, 0, nWrittenBytes)
		require.NotEqual(t, len(dummyData), nWrittenBytes)

		tmpFileContent, err := ioutil.ReadFile(tmpFile.Name())
		require.NoError(t, err)
		require.NotEqual(t, dummyData, tmpFileContent)
	})
}
