package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type student struct {
	Name   string `json:"name"`
	Rollno string `json:"rollno"`
	City   string `json:"city"`
}

type allStudents []student

var students = allStudents{
	{
		Name: "sai teja", Rollno: "1", City: "Nlg",
	},
	{
		Name: "sai sri", Rollno: "2", City: "Nlg",
	},
}

func createNewStudent(w http.ResponseWriter, r *http.Request) {
	var newStudent student
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the student name and city only in order to update")
	}

	json.Unmarshal(reqBody, &newStudent)
	students = append(students, newStudent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newStudent)
}

func getAllStudents(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: All Students Endpoint")
	json.NewEncoder(w).Encode(students)
}

func testPostStudents(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test POST endpoint worked")

}

func getOneStudent(w http.ResponseWriter, r *http.Request) {
	studentRollno := mux.Vars(r)["rollno"]

	for _, singleStudent := range students {
		if singleStudent.Rollno == studentRollno {
			json.NewEncoder(w).Encode(singleStudent)
		}
	}

}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	studentRollno := mux.Vars(r)["rollno"]
	var updatedStudent student
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the student name and city only in order to update")
	}
	json.Unmarshal(reqBody, &updatedStudent)

	for i, singleStudent := range students {
		if singleStudent.Rollno == studentRollno {
			singleStudent.Name = updatedStudent.Name
			singleStudent.Rollno = updatedStudent.Rollno
			singleStudent.City = updatedStudent.City
			students = append(students[:i], singleStudent)
			json.NewEncoder(w).Encode(singleStudent)
		}
	}
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	studentRollno := mux.Vars(r)["rollno"]
	for i, singleStudent := range students {
		if singleStudent.Rollno == studentRollno {
			students = append(students[:i], students[:i+1]...)
			fmt.Fprintf(w, "the student with Rollno %v has been deleted successfully", studentRollno)
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/students", getAllStudents).Methods("GET")
	myRouter.HandleFunc("/students", testPostStudents).Methods("POST")
	myRouter.HandleFunc("/students/{rollno}", getOneStudent).Methods("GET")
	myRouter.HandleFunc("/students/{rollno}", updateStudent).Methods("PATCH")
	myRouter.HandleFunc("/students", createNewStudent).Methods("POST")
	myRouter.HandleFunc("/students/{rollno}", deleteStudent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3001", myRouter))
}

func main() {
	handleRequests()
}
