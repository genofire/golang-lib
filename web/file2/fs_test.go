package file_test

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"

	"dev.sum7.eu/genofire/golang-lib/web/file2"
)

var fs file.FS

func ExampleFS() {
	// generate the UUID for the new file
	id := uuid.New()

	// store a file
	{
		f, _ := os.Open("glenda.png")
		fs.Store(id, "glenda.png", "image/png", f)
		f.Close()
	}

	// copy back to a local file
	{
		r, _ := fs.OpenUUID(id)
		w, _ := os.Create("glenda.png")
		io.Copy(w, r)
		r.Close()
		w.Close()
	}
}

func ExampleStoreFromHTTP() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if err := file.StoreFromHTTP(fs, r, "file"); err != nil {
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"message":"%v"}`, err)))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
}
