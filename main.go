// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"employee-management-api-407905.com/employee-management-system/employee"
	"employee-management-api-407905.com/employee-management-system/firestore"
	"github.com/gorilla/mux"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Your existing code here...

	// Initialize Firestore client
	firestore.InitializeFirestore()

	// Set up HTTP routes
	router := mux.NewRouter()
	router.HandleFunc("/employees", employee.GetAllEmployeesHandler).Methods("GET")
	router.HandleFunc("/employee/{id}", employee.GetEmployeeByIDHandler).Methods("GET")
	router.HandleFunc("/employee", employee.AddEmployeeHandler).Methods("POST")
	router.HandleFunc("/employee/{id}", employee.DeleteEmployeeHandler).Methods("DELETE")
	router.HandleFunc("/employee/{id}", employee.UpdateEmployeeHandler).Methods("PATCH")
	router.HandleFunc("/findemployee/search", employee.SearchEmployeeHandler).Methods("GET")

	// Custom middleware to log requests to undefined routes
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log the requested URL
			log.Printf("INFO: Requested URL: %s", r.URL.Path)

			// Continue with the next handler
			next.ServeHTTP(w, r)
		})
	})

	// Use the router as the handler
	router.ServeHTTP(w, r)
}

func CloudFunction(w http.ResponseWriter, r *http.Request) {
    // Cloud Functions will call the exported Handler function
    Handler(w, r)
}


func main() {
	// Cloud Functions will call the exported Handler function
	http.HandleFunc("/", Handler)

	// Start the server
	port := 8080
	fmt.Printf("Server listening on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
