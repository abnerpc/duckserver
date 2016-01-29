package server

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/abnerpc/duckdoc/common"
)

var config *common.Configuration

func main() {

	// load config
	var err error
	config, err = common.LoadConfiguration()
	if err != nil {
		fmt.Println("Error loading the configuration file")
		return
	}

	upload := http.HandlerFunc(uploadHandler)

	http.Handle("/upload", middlewareProtect(upload))
	fs := http.FileServer(http.Dir("docs"))
	http.Handle("/", fs)
	fmt.Println("Listening...")
	http.ListenAndServe(":8888", nil)
}

func middlewareProtect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = "breakpoint"
		key := r.Header.Get("Authorization")
		if _, ok := config.AdminKeys[key]; !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("User not authorized"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer file.Close()

	out, err := ioutil.TempFile("tmp", "upload")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	err = extract(out.Name(), "docs")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	os.Remove(out.Name())

	fmt.Fprintln(w, "Docs uploaded!")
}

func extract(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
