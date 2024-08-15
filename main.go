package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type CodeRecord struct {
	Code    string
	Expires time.Time
}

var codeRecords = make(map[string]CodeRecord)
var mutex = &sync.Mutex{}

func saveCodeRecord(code string) {
	mutex.Lock()
	defer mutex.Unlock()
	codeRecords[code] = CodeRecord{
		Code:    code,
		Expires: time.Now().Add(265 * 24 * time.Hour),
	}
}

func validateCode(code string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	record, ok := codeRecords[code]
	if !ok || record.Expires.Before(time.Now()) {
		return false
	}
	return true
}

func main() {
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/getcode", getCodeHandler)
	http.ListenAndServe(":8080", nil)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	cookie_code, err := r.Cookie("code")
	if err != nil {
		fmt.Printf("ERROR:\t[%v]\n", err)
		w.WriteHeader(403)
		fmt.Fprintf(w, "FAILED:\t[COOKIE NOT FOUND]")
	} else {
		code := cookie_code.Value
		if validateCode(code) {
			fmt.Printf("SUCCESS:[%s]\n", code)
			w.WriteHeader(200)
			fmt.Fprintf(w, "SUCCESS:[%s]", code)
		} else {
			fmt.Printf("FAILED:\t[%s]\n", code)
			w.WriteHeader(403)
			fmt.Fprintf(w, "FAILED:\t[%s]", code)
		}
	}
}

func getCodeHandler(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")

	uuidV4 := uuid.New()
	saveCodeRecord(uuidV4.String())
	http.SetCookie(w, &http.Cookie{
		Name:     "code",
		Value:    uuidV4.String(),
		Domain:   os.Getenv("DOMAIN"),
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})
	fmt.Printf("ADD:  \t[%s]\n", uuidV4.String())
	fmt.Printf("302:  \t[%s]\n", redirect)
	http.Redirect(w, r, redirect, http.StatusFound)
}
