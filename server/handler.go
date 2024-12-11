package handler

import (
    "fmt"
    "html/template"
    "net/http"
    "github.com/<votre-utilisateur-github>/go_course_manager/server/db"
)

var templates = template.Must(template.ParseGlob("server/templates/*.html"))

// HomeHandler affiche la page d'accueil
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    courses := db.GetCourses()
    templates.ExecuteTemplate(w, "index.html", courses)
}

// AddCourseHandler gère l'ajout d'un cours
func AddCourseHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        title := r.FormValue("title")
        content := r.FormValue("content")
        db.AddCourse(title, content)
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    templates.ExecuteTemplate(w, "add_course.html", nil)
}

// DownloadCourseHandler gère le téléchargement des cours en PDF
func DownloadCourseHandler(w http.ResponseWriter, r *http.Request) {
    pdfData := db.GeneratePDF()
    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "attachment; filename=courses.pdf")
    w.Write(pdfData)
}
