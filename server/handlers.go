package server

import (
	"cargo_service/internal/deals"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"strconv"
)

func dealCreate(w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get("user")
	slb, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var deal deals.Deal
	err = json.Unmarshal(slb, &deal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = deal.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(deal, user)

	err = deals.Create(Db, deal, user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func dealsList(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("user")
	page := r.URL.Query().Get("page")
	perPage := r.URL.Query().Get("perpage")

	if page == "" || perPage == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	offset := perPageInt * pageInt

	list, err := deals.List(Db, user, perPageInt, offset, pageInt)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func dealInfo(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("user")
	dealId := chi.URLParam(r, "ID")

	deal, err := deals.Info(Db, user, dealId)

	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(deal)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
