package node

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

func StartServer(port string) {

	http.HandleFunc("/calc", calculate)
	http.ListenAndServe(port, nil)

}

func calculate(w http.ResponseWriter, r *http.Request) {

	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil {
		http.Error(w, "invalid input n", http.StatusBadRequest)
		return
	}
	var sum float64 = 0
	var loop_end float64 = float64(n) * 1e6
	for i := 1.0; i <= loop_end; i++ {
		sum += ((math.Sqrt(i) * math.Sin(i)) / math.Log(i+1))
	}

	fmt.Fprintf(w, "calculated value = %f", sum)
}
