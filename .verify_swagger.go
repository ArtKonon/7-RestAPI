package main

import (
    "fmt"
    "net/http/httptest"
    "strings"

    "libratrack/internal/app"
)

func main() {
    h := app.NewRouter()
    req := httptest.NewRequest("GET", "/swagger/index.html", nil)
    w := httptest.NewRecorder()
    h.ServeHTTP(w, req)
    fmt.Printf("swagger_status=%d body_has_libratrack=%v\n", w.Code, strings.Contains(w.Body.String(), "Libratrack API"))
}
