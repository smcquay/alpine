package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const postBody = `
<!DOCTYPE html>
<html>
	<body>
		<form enctype="multipart/form-data" action="/upload/" method="POST">
			<input name="file" type="file" /><br />
			<input type="submit" value="Upload File" />
		</form>
	</body>
</html>
`

var port = flag.Int("port", 8000, "port from which to serve http")
var tlsport = flag.Int("tlsport", 8443, "port from which to serve https")
var hidden = flag.Bool("hidden", false, "allow serving hidden dirs")
var canUpload = flag.Bool("upload", false, "enable upload interface")

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		log.Printf("%s: %s\n", r.RemoteAddr, r.URL)
		if !*hidden && strings.Contains(r.URL.Path, "/.") {
			http.Error(w, "hidden files and directories are not allowed", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func upload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Printf("upload get")
		fmt.Fprintf(w, postBody)
	case "POST":
		file, header, err := r.FormFile("file")
		if err != nil {
			msg := fmt.Sprintf("problem picking off file from request: %v", err)
			http.Error(w, msg, http.StatusBadRequest)
			log.Printf(msg)
			return
		}
		log.Printf("upload for: %v", header.Filename)
		defer file.Close()
		f, err := os.Create(header.Filename)
		if err != nil {
			msg := fmt.Sprintf("problem creating upload file: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			log.Printf(msg)
			return
		}
		if _, err := io.Copy(f, file); err != nil {
			msg := fmt.Sprintf("problem copying file: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			log.Printf(msg)
			return
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func main() {
	flag.Parse()
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("problem getting hostname:", err)
	}
	addr := fmt.Sprintf(":%d", *port)
	tlsaddr := fmt.Sprintf(":%d", *tlsport)

	http.Handle("/", logger(http.FileServer(http.Dir("./"))))
	if *canUpload {
		log.Printf("WARNING: uploading enabled")
		http.HandleFunc("/upload/", upload)
	}

	key := os.Getenv("TLS_KEY")
	cert := os.Getenv("TLS_CERT")
	url := fmt.Sprintf("http://%s:%d/", hostname, *tlsport)
	if key != "" && cert != "" {
		tlsUrl := fmt.Sprintf("https://%s:%d/", hostname, *tlsport)
		go func() {
			log.Printf("serving redirect on: %s", url)
			sm := http.NewServeMux()
			sm.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
				http.Redirect(w, req, tlsUrl, http.StatusPermanentRedirect)
			})
			if err := http.ListenAndServe(addr, sm); err != nil {
				log.Fatal(err)
			}
		}()
		log.Printf("serving on: %s", tlsUrl)
		if err := http.ListenAndServeTLS(tlsaddr, cert, key, nil); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("serving on: %s", url)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal(err)
		}
	}
}
