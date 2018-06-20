package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type dataSet struct {
	data []int
}

func main() {
	port := os.Getenv("PORT")

	if port == ":" {
		log.Fatal("$PORT must be set")
	}
	http.HandleFunc("/", imageHandler)
	http.HandleFunc("/result/", imageHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(content))
}

func processDataSet(rawData string) []int {
	tempData := strings.Split(rawData, " ")
	data := []int{}

	for _, v := range tempData {
		tempVal, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println("Error converting data to int")
		}
		data = append(data, tempVal)
	}

	fmt.Printf("processed Data=%v", data)
	return data
}

func createImage(data []int) image.Image {
	const width, height = 256, 256

	buff := 10
	wbar := (width - buff) / len(data) // width of bar

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{255, 255, 255, 255})
		}
	}

	//honestly don't really understand how this works although I did it by myself
	//first range over data
	//then range over the height (notice how because pixels are 0,0 at the top left corner, y := height, and x := 0)
	//then range over the width
	for i, dp := range data {
		for y := height; y > (height - (dp * height / 100)); y-- {
			for x := wbar*i + buff; x <= wbar*(i+1); x++ {
				img.Set(x, y, color.NRGBA{24, 83, 150, 255})
				//wbar*i+x < wbar*(i+1)
			}
		}
	}

	var m image.Image = img
	return m
}

func imageHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method)
	r.ParseForm()

	if r.Method != "get" && r.FormValue("userData") != "" {

		fmt.Printf("userData=%s", r.FormValue("userData"))
		data := processDataSet(r.FormValue("userData"))
		img := createImage(data)
		writeImageWithTemplate(w, &img)
	} else {
		data := []int{10, 20, 50, 60, 44, 67, 33, 35} //expect this is a percentage
		img := createImage(data)
		writeImageWithTemplate(w, &img)
	}

}

// taken from sanarias.com
// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImageAsPng(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, *img); err != nil {
		log.Println("unable to encode image")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image")
	}
}

// taken from sanarias.com
// Writeimagewithtemplate encodes an image 'img' in jpeg format and writes it into ResponseWriter using a template.
func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {

	//imageTemplate := template.Must(template.ParseFiles("index.html"))

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.ParseFiles("index.html"); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
}
