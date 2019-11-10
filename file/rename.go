package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func RenamePictures(root string) (count int) {
	w := sync.WaitGroup{}

	if e := filepath.Walk(root, func(pathX string, infoX os.FileInfo, err error) error {

		// only handles files with the extension name
		if err == nil && !infoX.IsDir() {

			if ext := filepath.Ext(pathX); ext != "" {

				go func(pathX string, infoX os.FileInfo) {

					w.Add(1)

					dir := filepath.Dir(pathX)
					name := infoX.ModTime().Format("2006-01-02_150405")
					newPath := filepath.Join(dir, name+ext)

					if pathX != newPath {
						if er := os.Rename(pathX, newPath); err != nil {
							fmt.Println(er)
						} else {
							//fmt.Println(pathX, " -->", newPath)
							count++
						}
					}

					w.Done()

				}(pathX, infoX)
			}
		} else {
			fmt.Print(err)
		}

		return nil
	}); e != nil {
		fmt.Print(e)
	}

	w.Wait()

	return
}
