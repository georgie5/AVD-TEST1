package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Students struct {
	Student_id   int
	Student_name string
	Email        string
	Class        string
	Major        string
	Age          int
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	ts, err := template.ParseFiles("./home_page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

}

func (app *application) createStudentForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./students_form_page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

func (app *application) createStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}

	student_name := r.PostForm.Get("student_name")
	email := r.PostForm.Get("email")
	class := r.PostForm.Get("class")
	major := r.PostForm.Get("major")
	age := r.PostForm.Get("age")

	errors := make(map[string]string)
	//check each field to validate
	if strings.TrimSpace(student_name) == " " {
		errors["student_name"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(student_name) > 50 {
		errors["student_name"] = "This field is too long(maximum is 50 characters)"
	}

	if strings.TrimSpace(email) == "" {
		errors["email"] = "This field cannot be left blank"
	} else if !isValidEmail(email) {
		errors["email"] = "Invalid email format"
	}

	if strings.TrimSpace(class) == "" {
		errors["class"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(class) > 50 {
		errors["class"] = "This field is too long (maximum is 50 characters)"
	}

	if strings.TrimSpace(major) == "" {
		errors["major"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(major) > 50 {
		errors["major"] = "This field is too long (maximum is 50 characters)"
	}

	if age == "" {
		errors["age"] = "This field cannot be left blank"
	} else {

		ageInt, err := strconv.Atoi(age)
		if err != nil {
			errors["age"] = "Invalid age format"
		} else if ageInt < 0 {
			errors["age"] = "Age cannot be negative"
		}
	}

	// check if there are any errors in the map
	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}
	s := `
			INSERT INTO students (student_name, email, class, major, age)
			VALUES ($1, $2, $3, $4, $5)
		`
	_, err = app.db.Exec(s, student_name, email, class, major, age)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	// Redirect to the home page after successful form submission
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) displayStudents(w http.ResponseWriter, r *http.Request) {

	readStudents := `
		SELECT * 
		FROM students
	`
	rows, err := app.db.Query(readStudents)
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
	}

	defer rows.Close()

	var students []Students
	for rows.Next() {
		var s Students
		err = rows.Scan(&s.Student_id, &s.Student_name, &s.Email, &s.Class, &s.Major, &s.Age)

		if err != nil {
			log.Println(err.Error())
			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		students = append(students, s)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	// 	// Dsplay the studemts using a template
	ts, err := template.ParseFiles("./show_page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, students)

	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

}

func isValidEmail(email string) bool {

	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
