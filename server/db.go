package db

import (
    "database/sql"
    "log"
    "github.com/jung-kurt/gofpdf"
    _ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

// InitDB initialise la base de données SQLite
func InitDB() {
    var err error
    database, err = sql.Open("sqlite3", "data/database.db")
    if err != nil {
        log.Fatal(err)
    }
    createTableQuery := `
    CREATE TABLE IF NOT EXISTS courses (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        content TEXT
    )`
    _, err = database.Exec(createTableQuery)
    if err != nil {
        log.Fatal(err)
    }
}

// AddCourse ajoute un cours à la base de données
func AddCourse(title, content string) {
    _, err := database.Exec("INSERT INTO courses (title, content) VALUES (?, ?)", title, content)
    if err != nil {
        log.Println("Erreur lors de l'ajout du cours :", err)
    }
}

// GetCourses retourne la liste des cours
func GetCourses() []map[string]string {
    rows, err := database.Query("SELECT title, content FROM courses")
    if err != nil {
        log.Println("Erreur lors de la récupération des cours :", err)
        return nil
    }
    defer rows.Close()

    var courses []map[string]string
    for rows.Next() {
        var title, content string
        rows.Scan(&title, &content)
        courses = append(courses, map[string]string{"Title": title, "Content": content})
    }
    return courses
}

// GeneratePDF génère un PDF contenant les cours
func GeneratePDF() []byte {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "B", 16)

    courses := GetCourses()
    for _, course := range courses {
        pdf.Cell(40, 10, course["Title"])
        pdf.Ln(12)
        pdf.MultiCell(0, 10, course["Content"], "", "L", false)
    }

    var buf []byte
    pdf.OutputBuffer(&buf)
    return buf
}
