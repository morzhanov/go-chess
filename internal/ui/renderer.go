package ui

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/morzhanov/go-termui-chess/internal/board"
)

type RenderRes struct {
	Move  bool
	Board [][]string
}

const (
	boardReplaceString = "<!--BOARD-->"
)

//go:embed index.html
var indexHTML string

//go:embed board.html
var boardHTML string

type Renderer interface {
	Start() error
	Update(request string)
}

type renderer struct {
	board  board.Board
	router *mux.Router
	events chan struct{}
}

func NewRenderer(b board.Board) (Renderer, error) {
	router := mux.NewRouter()
	r := &renderer{
		board:  b,
		router: router,
		events: make(chan struct{}),
	}
	router.HandleFunc("/", r.handleIndex)
	router.HandleFunc("/sse", r.handleEvents)
	staticPath, err := filepath.Abs("./internal/ui/static/")
	if err != nil {
		return nil, err
	}
	dir := http.Dir(staticPath)
	log.Println(dir)
	router.PathPrefix("/").Handler(http.FileServer(dir))
	return r, nil
}

func (r *renderer) Start() error {
	log.Println("starting chess server on localhost:5000")
	return http.ListenAndServe(":5000", r.router)
}

func (r *renderer) Update(request string) {
	log.Println("handling update: ", request)
	log.Println("board before update: ", r.board)
	if err := r.board.Update(request); err != nil {
		log.Println(err)
		return
	}
	log.Println("board after update: ", r.board)
	r.events <- struct{}{}
}

func (r *renderer) handleIndex(w http.ResponseWriter, _ *http.Request) {
	b, err := r.renderBoard()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := strings.Replace(indexHTML, boardReplaceString, b, -1)
	if _, err = w.Write([]byte(res)); err != nil {
		log.Println("failed to render board: ", err.Error())
	}
}

func (r *renderer) handleEvents(w http.ResponseWriter, _ *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	defer flusher.Flush()
	addSSEHeaders(w)

	initialSemaphore := 2
	for {
		if initialSemaphore == 0 {
			<-r.events
		} else {
			initialSemaphore--
		}
		log.Println("handling event in the for loop")
		b, err := r.renderBoard()
		if err != nil {
			log.Println("failed to render board: ", err)
			continue
		}
		if _, err = fmt.Fprintf(w, "data: %s\n\n", b); err != nil {
			log.Println("failed to send events: ", err)
			r.events <- struct{}{}
			continue
		}
		log.Println("flushing in the for loop")
	}
}

func (r *renderer) renderBoard() (string, error) {
	log.Println("in render board")
	t := template.New("")
	t.Funcs(template.FuncMap{"mod": func(i, j, v int) bool {
		if i%v != 0 {
			return j%v != 0
		}
		return j%v == 0
	}})
	tmpl, err := t.Parse(boardHTML)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err = tmpl.Execute(&tpl, &RenderRes{
		Move:  r.board.GetMove(),
		Board: r.board.GetPieces(),
	}); err != nil {
		return "", err
	}

	res := tpl.String()
	res = strings.Replace(res, "\n", "", -1)
	return res, nil
}

func addSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}
