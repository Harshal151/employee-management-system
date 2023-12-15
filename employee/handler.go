// employee/handler.go
package employee

import (
	"employee-management-api-407905.com/employee-management-system/firestore"
	"employee-management-api-407905.com/employee-management-system/sharedpackage"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

func GetAllEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve all employees
	employees, err := firestore.GetAllEmployees()
	if err != nil {
		log.Printf("ERROR: Error getting employees: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("INFO: Retrieved %d employees", len(employees))

	// Convert employees to JSON format
	employeesJSON, err := json.Marshal(employees)
	if err != nil {
		log.Printf("ERROR: Error marshaling employees to JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("INFO: Successfully marshaled employees to JSON")

	w.Header().Set("Content-Type", "application/json")

	w.Write(employeesJSON)
	log.Printf("INFO: Sent employees JSON response")
}

// GetEmployeeByIDHandler is an HTTP handler function to retrieve and respond with an employee by ID.
func GetEmployeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Printf("WARN: Invalid URL in GetEmployeeByIDHandler")
		return
	}

	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		log.Printf("WARN: Invalid employee ID in GetEmployeeByIDHandler: %v", err)
		return
	}

	log.Printf("INFO: Request received for employee with ID: %d", employeeID)

	// Get the employee by ID
	employee, err := firestore.GetEmployeeByID(employeeID)
	if err != nil {
		log.Printf("ERROR: Error getting employee by ID: %v", err)
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	log.Printf("INFO: Employee found with ID: %d", employeeID)

	// Convert employee to JSON format
	employeeJSON, err := json.Marshal(employee)
	if err != nil {
		log.Printf("ERROR: Error marshaling employee to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set content type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON formatted employee to the response
	w.Write(employeeJSON)
	log.Printf("INFO: Employee response sent for ID: %d", employeeID)
}

// AddEmployeeHandler is an HTTP handler function to add a new employee to Firestore.
func AddEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var newEmployee sharedpackage.Employee
	err := json.NewDecoder(r.Body).Decode(&newEmployee)
	if err != nil {
		log.Printf("ERROR: Error parsing request body: %v", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: New employee data received: %+v", newEmployee)

	// Check if the employee ID already exists in the database
	if exists, err := firestore.EmployeeExists(newEmployee.ID); err != nil {
		log.Printf("ERROR: Error checking if employee ID exists: %v", err)
		http.Error(w, "Error checking employee ID", http.StatusInternalServerError)
		return
	} else if exists {
		log.Printf("ERROR: Employee with ID %d already exists", newEmployee.ID)
		http.Error(w, "Employee with the given ID already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newEmployee.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}

	newEmployee.Password = string(hashedPassword)

	// Call AddEmployee function to store the new employee in Firestore
	err = firestore.AddEmployee(newEmployee)
	if err != nil {
		log.Printf("ERROR: Error adding employee to Firestore: %v", err)
		http.Error(w, "Error adding employee to Firestore", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Employee added successfully to Firestore")

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Employee added successfully to Firestore"))
}

// UpdateEmployeeHandler is an HTTP handler function to update an employee's information.
func UpdateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the employee ID from the URL using Gorilla Mux
	vars := mux.Vars(r)
	employeeIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Println("WARN: Invalid URL in UpdateEmployeeHandler")
		return
	}

	// Convert employee ID from string to integer
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		log.Printf("WARN: Invalid employee ID in UpdateEmployeeHandler: %v", err)
		return
	}

	log.Printf("INFO: UpdateEmployeeHandler - Request received for updating employee with ID: %d", employeeID)

	// Parse the request body to get the updated fields
	var updatedFields map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&updatedFields)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("ERROR: Error decoding request body: %v", err)
		return
	}

	defer r.Body.Close()

	log.Printf("INFO: UpdateEmployeeHandler - Decoded request body with updated fields: %+v", updatedFields)

	// Perform the update in Firestore
	err = firestore.UpdateEmployee(employeeID, updatedFields)
	if err != nil {
		log.Printf("ERROR: UpdateEmployeeHandler - Error updating employee by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Employee updated successfully"))
	log.Printf("INFO: UpdateEmployeeHandler - Employee with ID %d updated successfully", employeeID)
}

// DeleteEmployeeHandler is an HTTP handler function to delete an employee by ID.
func DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Println("WARN: Invalid URL in DeleteEmployeeHandler")
		return
	}

	// Convert employeeID from string to integer
	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		log.Printf("WARN: Invalid employee ID in DeleteEmployeeHandler: %v", err)
		return
	}

	log.Printf("INFO: DeleteEmployeeHandler - Request received for deleting employee with ID: %d", employeeID)

	// Delete the employee by ID
	err = firestore.DeleteEmployee(employeeID)
	if err != nil {
		log.Printf("ERROR: DeleteEmployeeHandler - Error deleting employee by ID: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Employee deleted successfully"))
	log.Printf("INFO: DeleteEmployeeHandler - Employee with ID %d deleted successfully", employeeID)
}

func SearchEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters from the request URL
	queryParams := r.URL.Query()

	// Convert query parameters to map[string]interface{}
	searchParams := make(map[string]string)
	for key, values := range queryParams {
		if len(values) > 0 {
			searchParams[key] = values[0]
		}
	}

	// Log the received search parameters at INFO level
	log.Printf("INFO: Received search parameters: %+v", searchParams)

	// Call the SearchEmployee function with the correct parameter type
	employees, err := firestore.SearchEmployee(searchParams)
	if err != nil {
		log.Printf("ERROR: Error searching employees: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the number of employees found at INFO level
	log.Printf("INFO: Number of employees found: %d", len(employees))

	// Convert employees to JSON format
	employeesJSON, err := json.Marshal(employees)
	if err != nil {
		log.Printf("ERROR: Error marshaling employees to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set content type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON formatted employees to the response
	w.Write(employeesJSON)

	// Log the successful response at INFO level
	log.Printf("INFO: SearchEmployeeHandler executed successfully")
}
