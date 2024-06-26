package vbd_test

import (
	"bytes"
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildbuddy-io/buildbuddy/enterprise/server/remote_execution/vbd"
	"github.com/buildbuddy-io/buildbuddy/server/testutil/testfs"
	"github.com/buildbuddy-io/buildbuddy/server/testutil/testmount"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	testmount.RunWithLimitedMountPermissions(m)
}

func TestVBD(t *testing.T) {
	ctx := context.Background()
	root := testfs.MakeTempDir(t)

	// Set up backing file
	f, err := os.CreateTemp(root, "vbd-*")
	require.NoError(t, err)
	t.Cleanup(func() { f.Close() })
	const fSize = 1024 * 512
	b := make([]byte, fSize)
	copy(b, []byte("Hello world!"))
	_, err = f.Write(b)
	require.NoError(t, err)

	// Mount as VBD
	v, err := vbd.New(&FileBlockDevice{f})
	require.NoError(t, err)
	dir, err := os.MkdirTemp(root, "mount-*")
	require.NoError(t, err)
	err = v.Mount(ctx, dir)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := v.Unmount(ctx)
		require.NoError(t, err)
	})

	// Try stat() on the virtual file
	s, err := os.Stat(filepath.Join(dir, vbd.FileName))
	require.NoError(t, err)
	require.Equal(t, int64(fSize), s.Size())

	// Try cleaning up VBD mounts while the VBD is mounted; this should succeed,
	// and should not affect our currently mounted VBD.
	err = vbd.CleanStaleMounts()
	require.NoError(t, err)

	// Try random reads and writes to the virtual file
	{
		f, err := os.OpenFile(filepath.Join(dir, vbd.FileName), os.O_RDWR, 0)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := f.Close()
			require.NoError(t, err)
		})

		for i := 1; i <= 100; i++ {
			offset, length := randSubslice(len(b))
			p := make([]byte, length)
			if shouldRead := rand.Intn(2) == 0; shouldRead {
				// Read, randomly choosing between using Seek+Read / ReadAt.
				if shouldSeek := rand.Intn(2) == 0; shouldSeek {
					_, err := f.Seek(int64(offset), 0)
					require.NoError(t, err)
					_, err = f.Read(p)
					require.NoError(t, err)
					require.True(t, bytes.Equal(p, b[offset:offset+length]))
				} else {
					_, err := f.ReadAt(p, int64(offset))
					require.NoError(t, err)
					require.True(t, bytes.Equal(p, b[offset:offset+length]))
				}
			} else {
				// Write
				_, err := rand.Read(p)
				require.NoError(t, err)
				_, err = f.WriteAt(p, int64(offset))
				require.NoError(t, err)
				// Update our expected buffer contents
				copy(b[offset:offset+length], p)
			}
		}
		// Try sync() on the virtual file (this is a NOP for now, but should at
		// least not fail)
		err = f.Sync()
		require.NoError(t, err)
	}
}

type FileBlockDevice struct{ *os.File }

func (f *FileBlockDevice) SizeBytes() (int64, error) {
	s, err := f.File.Stat()
	if err != nil {
		return 0, err
	}
	return s.Size(), nil
}

// Picks a uniform random subslice of a slice with a given length.
// Returns the offset and length of the subslice.
func randSubslice(sliceLength int) (offset, length int) {
	length = rand.Intn(sliceLength + 1)
	offset = rand.Intn(sliceLength - length + 1)
	return
}
