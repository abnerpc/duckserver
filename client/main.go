package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {

	ff := flag.String("f", "", "Folder or File to be compressed")
	fileName := flag.String("o", "zipped.zip", "Output zip name")
	flag.Parse()

	Zipit(*ff, *fileName)

	response, err := sendPost(*fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	text, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(text))

}

func sendPost(path string) (*http.Response, error) {

	file, err := os.Open("./" + path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", "http://localhost:8888/upload", body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", writer.FormDataContentType())
	//r.Header.Set("Authorization", "123")
	client := &http.Client{}
	response, err := client.Do(r)
	return response, err
}
