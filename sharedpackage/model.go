// employee/model.go
package sharedpackage

// import "golang.org/x/text/number"

// "fmt"

type Employee struct {
	ID        int  `firestore:"id" json:"id"`
	FirstName string  `firestore:"firstName" json:"firstName"`
	LastName  string  `firestore:"lastName" josn:"lastName"`
	Email     string  `firestore:"email" json:"email"`
	Password  string  `firestore:"password" json:"password"`
	PhoneNo   int  `firestore:"phoneNo" json:"phoneNo"`
	Role      string  `firestore:"role" json:"role"`
	Salary    float64 `firestore:"salary" json:"salary"`
}

