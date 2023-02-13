package file

import (
	"io"
	"os"
)

func Copy(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() {
		_ = in.Close()
	}()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		closeErr := out.Close()
		if err == nil {
			err = closeErr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
