package audioreader

import "os"

type fileReader struct {
	*os.File
}

func (f fileReader) Size() (int64, error) {
	fstat, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return fstat.Size(), nil
}
