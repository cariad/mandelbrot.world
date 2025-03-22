package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/cariad/gandelbrot"
)

const port = 8080

func handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "immutable, max-age=3600, no-transform, public")
		http.FileServer(http.Dir("frontend")).ServeHTTP(w, r)
	}
}

func main() {
	http.Handle("/", handleRoot())
	http.HandleFunc("/tiles/{z}/{x}/{y}", tilesHandler)

	log.Printf("Listening on port %d…", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func readIntParam(r *http.Request, p string) int {
	s := r.PathValue(p)
	i, err := strconv.Atoi(s)

	if err != nil {
		log.Fatalf("failed to convert %s=%s to int\n", p, s)
	}

	return i
}

func tilesHandler(w http.ResponseWriter, r *http.Request) {
	maxIterations := 800

	if i := r.URL.Query().Get("max_iterations"); i != "" {
		if mi, err := strconv.Atoi(i); err == nil {
			maxIterations = mi
		} else {
			log.Printf("failed to convert max_iterations to int (%s)\n", i)
		}
	}

	threadCount := 4

	if tc := os.Getenv("MW_THREAD_COUNT"); tc != "" {
		if mi, err := strconv.Atoi(tc); err == nil {
			threadCount = mi
		} else {
			log.Printf("failed to convert MW_THREAD_COUNT to int (%s)\n", tc)
		}
	}

	w.Header().Set("Cache-Control", "immutable, max-age=86400, no-transform, public")
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)

	x := float64(readIntParam(r, "x"))
	y := float64(readIntParam(r, "y"))
	z := float64(readIntParam(r, "z"))

	complexWidth := 4.0 / math.Pow(2, z)

	gandelbrot.Render(&gandelbrot.RenderArgs{
		ComplexWidth:  complexWidth,
		Imaginary:     (complexWidth * y) - 2.0,
		MaxIterations: maxIterations,
		Real:          (complexWidth * x) - 2.0,
		RenderWidth:   400,
		ThreadCount:   threadCount,
		Writer:        w,
	})
}
