package main

import (
        "fmt"
        "io"
        "net/http"
        "os"
)

 func uploadHandler(w http.ResponseWriter, r *http.Request) {
                
        // the FormFile function takes in the POST input id file
        file, _, err := r.FormFile("file")

        if err != nil {
                fmt.Fprintln(w, err)
                return
        }


        defer file.Close()

        out, err := os.Create("./tmp/upload.zip")
        if err != nil {
                fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
                return
        }

        defer out.Close()

        // write the content from POST to the file
        _, err = io.Copy(out, file)
        if err != nil {
                fmt.Fprintln(w, err)
        }

        err = Unzip("./tmp/upload.zip", "./docs")
        if err != nil {
                fmt.Fprintln(w, err)
        }
        
        fmt.Println("Docs uploaded!")
 }

 func main() {
        http.HandleFunc("/upload", uploadHandler)
        fs := http.FileServer(http.Dir("docs"))
        http.Handle("/", fs)
        fmt.Println("Listening...")
        http.ListenAndServe(":8888", nil)
 }
