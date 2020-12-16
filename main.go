package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Score float64 `json:"score"`
}

func main() {
	api := gin.Default()

	api.POST("/face-mesh", func(c *gin.Context) {

		faceMesh, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			errResp(c, err)
			return
		}

		tmpfile, err := ioutil.TempFile("", "face-mesh")
		if err != nil {
			errResp(c, err)
			return
		}
		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write(faceMesh); err != nil {
			errResp(c, err)
			return
		}
		if err := tmpfile.Close(); err != nil {
			errResp(c, err)
			return
		}

		cmdStr := fmt.Sprintf("./verify.py --face-mesh %s", tmpfile.Name())
		cmd := exec.Command("bash", "-c", cmdStr)
		output, err := cmd.Output()
		if err != nil {
			errResp(c, err)
			return
		}

		var response Response
		if err := json.Unmarshal(output, &response); err != nil {
			errResp(c, err)
			return
		}

		c.JSON(http.StatusOK, response)
	})

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.RequestURI, "/api/") {
			req.URL.Path = req.URL.Path[4:]
			api.ServeHTTP(resp, req)
		} else {
			serveStaticAssets(resp, req)
		}
	})

	fmt.Println("Serving on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

var publicDir = "."

func serveStaticAssets(resp http.ResponseWriter, req *http.Request) {
	p := filepath.Join(publicDir, filepath.Clean(req.URL.Path))
	if info, err := os.Stat(p); err != nil {
		http.ServeFile(resp, req, filepath.Join(publicDir, "index.html"))
	} else if info.IsDir() {
		http.ServeFile(resp, req, filepath.Join(publicDir, "index.html"))
	} else {
		http.ServeFile(resp, req, p)
	}
}

func errResp(c *gin.Context, err error) {
	fmt.Printf("Error: %s\n", err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}
