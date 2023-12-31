package uploaders

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

type FtpUploader struct {
	Username string
	Password string
	Host     string
	Port     string

	Timeout time.Duration
}

func (u FtpUploader) Upload(ctx context.Context, local, remote string) error {
	c, err := ftp.Dial(
		fmt.Sprintf("%s:%s", u.Host, u.Port),
		ftp.DialWithContext(ctx),
		ftp.DialWithTimeout(u.Timeout),
	)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}

	err = c.Login(u.Username, u.Password)
	if err != nil {
		return fmt.Errorf("login error: %w", err)
	}

	ee, err := c.List("")
	if err != nil {
		return fmt.Errorf("list error: %w", err)
	}

	for _, e := range ee {
		if e.Name == remote {
			if e.Type != ftp.EntryTypeFolder {
				return fmt.Errorf("destination exists but is not a directory")
			}

			err = c.RemoveDirRecur(remote)
			if err != nil {
				return fmt.Errorf("cleaning destination before upload error: %w", err)
			}

			break
		}
	}

	err = c.MakeDir(remote)
	if err != nil {
		return fmt.Errorf("create remote directory error: %w", err)
	}

	err = c.ChangeDir(remote)
	if err != nil {
		return fmt.Errorf("change directory error: %w", err)
	}

	err = filepath.WalkDir(local, func(path string, d fs.DirEntry, err error) error {
		if path == local {
			return nil
		}

		rpath, err := filepath.Rel(local, path)
		if err != nil {
			return fmt.Errorf("find rel path error: %w", err)
		}

		if d.IsDir() {
			err = c.MakeDir(rpath)
			if err != nil {
				return fmt.Errorf("create remote directory error: %w", err)
			}

			return nil
		}

		fh, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("local file read error: %w", err)
		}
		defer func() { _ = fh.Close() }()

		err = c.Stor(rpath, fh)
		if err != nil {
			return fmt.Errorf("file upload error: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk dist directory error: %w", err)
	}

	err = c.Quit()
	if err != nil {
		return fmt.Errorf("close error: %w", err)
	}

	return nil
}
