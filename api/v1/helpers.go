package v1

import (
    "encoding/json"
    "net/http"
    "strconv"
)

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	e := json.NewEncoder(w)
	e.Encode(map[string]string{"error": message})
}

func respondOk(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	e := json.NewEncoder(w)
	e.Encode(body)
}

func getPagerParam(r *http.Request, key string, value int) int {
    p := r.URL.Query().Get(key)

    if n, e := strconv.ParseInt(p, 10, 64); e == nil && n > 0 {
        return int(n)
    }

    return value
}

func getPager(r *http.Request) (int, int) {
    start := getPagerParam(r, "start", 0)
    limit := getPagerParam(r, "limit", 50)

    if limit > 250 {
        limit = 250
    }

    return start, limit
}
