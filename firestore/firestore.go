// firestore/firestore.go
package firestore

import (
	"context"
	"employee-management-api-407905.com/employee-management-system/sharedpackage"
	"errors"
	"fmt"
	"log"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var FirestoreClient *firestore.Client

func InitializeFirestore() error {
	ctx := context.Background()
	
	projectID := "employee-management-api-407905"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	FirestoreClient = client
	fmt.Println("Firestore client successfully initialized")
	return nil
}

// CloseFirestore closes the Firestore client
func CloseFirestore() {
	if FirestoreClient != nil {
		FirestoreClient.Close()
		fmt.Println("Firestore client closed")
	}
}

// EmployeeExists checks if an employee with the given ID already exists in the database
func EmployeeExists(employeeID int) (bool, error) {
	ctx := context.Background()

	iter := FirestoreClient.Collection("employees").Where("id", "==", employeeID).Documents(ctx)
	defer iter.Stop()

	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("ERROR: Error iterating through documents: %v", err)
			return false, err
		}

		return true, nil
	}

	return false, nil
}

// AddEmployee adds a new employee to Firestore
func AddEmployee(employee sharedpackage.Employee) error {
	// Check if the employee ID already exists in the database
	if exists, err := EmployeeExists(employee.ID); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("Employee with ID %d already exists", employee.ID)
	}

	ctx := context.Background()
	_, _, err := FirestoreClient.Collection("employees").Add(ctx, employee)
	if err != nil {
		log.Printf("ERROR: Error adding employee to Firestore: %v", err)
		return err
	}

	log.Printf("INFO: Employee added to Firestore: %+v", employee)
	return nil
}

// GetAllEmployees retrieves all employees from the Firestore database.
func GetAllEmployees() ([]sharedpackage.Employee, error) {
	ctx := context.Background()

	iter := FirestoreClient.Collection("employees").Documents(ctx)

	var employees []sharedpackage.Employee

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("ERROR: Error iterating over documents: %v", err)
			return nil, err
		}

		var employeeData sharedpackage.Employee
		if err := doc.DataTo(&employeeData); err != nil {
			log.Printf("ERROR: Error converting document data: %v", err)
			return nil, err
		}

		employees = append(employees, employeeData)
		log.Printf("INFO: Processed employee %d", employeeData.ID)
	}

	return employees, nil
}

// getEmployeeByID gets an employee by ID from Firestore.
func GetEmployeeByID(employeeID int) (*sharedpackage.Employee, error) {
	ctx := context.Background()

	iter := FirestoreClient.Collection("employees").Where("id", "==", employeeID).Documents(ctx)

	defer iter.Stop()

	log.Printf("INFO: Searching for employee with ID: %d", employeeID)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("ERROR: Error iterating through documents: %v", err)
			return nil, err
		}

		var employee sharedpackage.Employee
		if err := doc.DataTo(&employee); err != nil {
			log.Printf("ERROR: Error converting document data: %v", err)
			return nil, err
		}

		log.Printf("INFO: Found employee with ID: %d", employeeID)
		return &employee, nil
	}

	log.Printf("WARN: Employee not found with ID: %d", employeeID)
	return nil, fmt.Errorf("ERROR: Document not found for ID: %d", employeeID)
}

// DeleteEmployee deletes an employee with a given ID from Firestore.
func DeleteEmployee(employeeID int) error {
	ctx := context.Background()

	collectionName := "employees"
	fieldName := "id"

	log.Printf("INFO: DeleteEmployee - Querying documents in collection '%s' where '%s' == %d", collectionName, fieldName, employeeID)

	query := FirestoreClient.Collection(collectionName).Where(fieldName, "==", employeeID)

	iter := query.Documents(ctx)
	defer iter.Stop()

	numDocuments := 0

	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("ERROR: DeleteEmployee - Error iterating through documents: %v", err)
			return err
		}
		numDocuments++
	}

	log.Printf("INFO: DeleteEmployee - Number of documents found: %d", numDocuments)

	if numDocuments == 0 {
		
		err := errors.New("No document found with given ID")
		log.Printf("INFO: DeleteEmployee - No documents found with ID = %d", employeeID)
		return err
	}

	iter.Stop()
	iter = query.Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("ERROR: DeleteEmployee - Error iterating through documents: %v", err)
			return err
		}

		log.Printf("INFO: DeleteEmployee - Deleting document with ID: %s", doc.Ref.ID)

		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			log.Printf("ERROR: DeleteEmployee - Error deleting document: %v", err)
			return err
		}
	}

	log.Printf("INFO: DeleteEmployee - Employee with ID %d deleted successfully", employeeID)
	return nil
}

// UpdateEmployee updates specific fields of an employee with a given ID in Firestore.
func UpdateEmployee(employeeID int, updatedFields map[string]interface{}) error {
	ctx := context.Background()

	collectionName := "employees"
	fieldName := "id"

	log.Printf("INFO: UpdateEmployee - Querying documents in collection '%s' where '%s' == %d", collectionName, fieldName, employeeID)

	// Query to find the document with the specified id
	query := FirestoreClient.Collection(collectionName).Where(fieldName, "==", employeeID)

	// Get the iterator for the query results
	iter := query.Documents(ctx)

	// Checking if the document with the specified id exists
	doc, err := iter.Next()
	if err == iterator.Done {
		log.Printf("WARN: UpdateEmployee - No documents found with ID = %d", employeeID)
		return nil
	}
	if err != nil {
		log.Printf("ERROR: UpdateEmployee - Error iterating through documents: %v", err)
		return err
	}

	log.Printf("INFO: UpdateEmployee - Updating document with ID: %s", doc.Ref.ID)

	// Creating a map to store the fields to update
	updatedMap := make(map[string]interface{})

	log.Printf("INFO: UpdateEmployee - Updating specific fields in map")

	for key, value := range updatedFields {
		updatedMap[key] = value
	}

	// Perform the update using the provided data
	_, err = doc.Ref.Set(ctx, updatedMap, firestore.MergeAll)
	if err != nil {
		log.Printf("ERROR: UpdateEmployee - Error updating employee: %v", err)
		return err
	}

	log.Printf("INFO: UpdateEmployee - Employee updated successfully with ID: %d", employeeID)
	return nil
}

// SearchEmployee searches for employees in Firestore based on specified query parameters.
func SearchEmployee(searchParams map[string]string) ([]*sharedpackage.Employee, error) {
	ctx := context.Background()

	
	collectionName := "employees"

	log.Printf("INFO: Received search parameters: %+v", searchParams)

	query := FirestoreClient.Collection(collectionName)
	var employees []*sharedpackage.Employee

	// Applying query conditions based on the provided search parameters
	for key, value := range searchParams {
		
		stringValue := fmt.Sprintf("%v", value)

		log.Printf("INFO: Applying query condition: %s == %s", key, stringValue)

		// Applying the query condition
		query1 := query.Where(key, "==", stringValue)

		// Get the iterator for the query results
		iter := query1.Documents(ctx)

		// Iterate over the documents and convert them to Employee objects
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Printf("ERROR: Error iterating through documents: %v", err)
				return nil, err
			}

			var employeeData sharedpackage.Employee
			if err := doc.DataTo(&employeeData); err != nil {
				log.Printf("ERROR: Error converting document data: %v", err)
				return nil, err
			}

			// Log the converted employee data at INFO level
			log.Printf("INFO: Converted employee data: %+v", employeeData)

			employees = append(employees, &employeeData)
		}
	}
	return employees, nil
}
