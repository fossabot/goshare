package api

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/ligenhw/goshare/session"

	"github.com/ligenhw/goshare/blog"
)

// /user
func Get(w http.ResponseWriter, r *http.Request) (err error) {
	var result interface{}
	id, err := strconv.Atoi(path.Base(r.URL.Path))

	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		bd := &blog.BlogDetail{
			Blog: blog.Blog{Id: id},
		}
		bd.QueryByID()
		result = bd
	} else {
		result, err = blog.GetAllBlogs()
		if err != nil {
			return err
		}
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(result)
	return
}

// /user
func Post(w http.ResponseWriter, r *http.Request) (err error) {
	// TODO: close the body
	decoder := json.NewDecoder(r.Body)
	b := blog.Blog{}
	err = decoder.Decode(&b)
	if err != nil {
		return
	}
	err = b.Create()
	return
}

// /user/:id
func Delete(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	b := blog.Blog{Id: id}
	err = b.Delete()
	return
}

// /user
func Put(w http.ResponseWriter, r *http.Request) (err error) {
	decoder := json.NewDecoder(r.Body)
	b := blog.Blog{}
	err = decoder.Decode(&b)
	if err != nil {
		return
	}

	err = b.Update()
	return
}

func BlogHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL, r.Method)

	var err error
	switch r.Method {
	case http.MethodGet:
		err = Get(w, r)
	case http.MethodPost:
		err = Post(w, r)
	case http.MethodDelete:
		err = Delete(w, r)
	case http.MethodPut:
		err = Put(w, r)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WithSession(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := session.Instance.SessionStart(w, r)
		if err != nil {
			log.Println(err)
			log.Println(session)
		}
		handler(w, r)
	}
}

func init() {
	http.HandleFunc("/api/blog/", WithSession(BlogHandler))
}
