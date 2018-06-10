package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	const width, height = 256, 256

	data := []int{10, 20, 50, 60, 44, 67, 33, 35} //expect this is a percentage

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

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	port := "8000" //os.Getenv("PORT")

	if port == ":" {
		log.Fatal("$PORT must be set")
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, string(content))
}
