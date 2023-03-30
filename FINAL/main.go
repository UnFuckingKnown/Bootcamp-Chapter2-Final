package main

import (
	"belajar-golang/connection"
	"belajar-golang/middleware"
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	IsLogin  bool
	Username string
}

var Data = MetaData{}

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Project struct {
	Id           int
	Name         string
	Start_date   time.Time
	End_date     time.Time
	Description  string
	Technologies []string
	Image        string
	NewPostdate  int
	Author       string
	Dif          string
	IsLogin      bool
}

var Projects = []Project{}

func main() {

	router := mux.NewRouter()

	connection.DatabaseConnect()

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	router.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/project", project).Methods("GET")
	router.HandleFunc("/mainblog/{id}", mainblog).Methods("GET")
	router.HandleFunc("/new-blog", middleware.UploadFile(newblog)).Methods("POST")
	router.HandleFunc("/delete/{id}", delete).Methods("GET")
	router.HandleFunc("/contact", contact).Methods("GET")
	router.HandleFunc("/login", formlogin).Methods("GET")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", formregister).Methods("GET")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/update/{id}", formupdate).Methods("GET")
	router.HandleFunc("/update", middleware.UploadFile(update)).Methods("POST")

	fmt.Println("server running on port 5000")
	http.ListenAndServe("localhost:5000", router)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html ; Charset=utf-8")
	w.WriteHeader(http.StatusOK)

	templ, err := template.ParseFiles("html/index.html")

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data.IsLogin = false
	} else {
		Data.IsLogin = session.Values["IsLogin"].(bool)
		Data.Username = session.Values["Name"].(string)
	}

	rows, err := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects;")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var result []Project
	for rows.Next() {
		var each = Project{}

		err := rows.Scan(&each.Id, &each.Name, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image)

		each.Dif = countDuration(each.Start_date, each.End_date)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)

	}
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("message" + err.Error()))
		return
	}

	resp := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
	}

	templ.Execute(w, resp)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html ; Charset=utf-8")
	w.WriteHeader(http.StatusOK)
	templ, err := template.ParseFiles("html/contact.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	templ.Execute(w, nil)

}

func formlogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html ; Charset=utf-8")
	w.WriteHeader(http.StatusOK)
	templ, err := template.ParseFiles("html/login.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	templ.Execute(w, nil)

}

func formregister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html ; Charset=utf-8")
	w.WriteHeader(http.StatusOK)
	templ, err := template.ParseFiles("html/register.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	templ.Execute(w, nil)

}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user(name, email, password) VALUES ($1, $2, $3);", name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, email, password FROM tb_user WHERE email=$1", email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["ID"] = user.Id
	session.Options.MaxAge = 10800

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html ; Charset=utf-8")

	templ, err := template.ParseFiles("html/blog.html")

	if err != nil {

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("message  " + err.Error()))
		return
	}

	templ.Execute(w, nil)
}

func mainblog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html ; Charset=utf-8")
	w.WriteHeader(http.StatusOK)

	templ, err := template.ParseFiles("html/mainblog.html")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {

		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("message  " + err.Error()))
		return
	}

	CatchBlog := Project{}
	for i, data := range Projects {
		if i == id {
			CatchBlog = Project{
				Name:         data.Name,
				Start_date:   data.Start_date,
				End_date:     data.End_date,
				Description:  data.Description,
				Technologies: data.Technologies,
				Image:        data.Image,
			}
		}
	}

	var resp = map[string]interface{}{
		"Data":  nil,
		"Blogs": CatchBlog,
	}

	templ.Execute(w, resp)
}

func newblog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
		return
	}

	Projectname := r.PostForm.Get("projectname")
	StartDate := r.PostForm.Get("startDate")
	EndDate := r.PostForm.Get("endDate")
	Description := r.PostForm.Get("description")
	technologi1 := r.PostForm.Get("technologi1")
	technologi2 := r.PostForm.Get("technologi2")
	technologi3 := r.PostForm.Get("technologi3")
	technologi4 := r.PostForm.Get("technologi4")
	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	technologis := []string{}

	if technologi1 != "" {
		technologis = append(technologis, technologi1)
	}

	if technologi2 != "" {
		technologis = append(technologis, technologi2)
	}

	if technologi3 != "" {
		technologis = append(technologis, technologi3)
	}
	if technologi4 != "" {
		technologis = append(technologis, technologi4)
	}

	start, _ := time.Parse("2006-01-02", StartDate)
	end, _ := time.Parse("2006-01-02", EndDate)

	var refilData = Project{
		Name:         Projectname,
		Start_date:   start,
		End_date:     end,
		Description:  Description,
		Technologies: technologis,
		Image:        image,
	}

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects ( name, start_date, end_date, description, technologies, image) VALUES ( $1, $2, $3, $4, $5, $6); ", refilData.Name, refilData.Start_date, refilData.End_date, refilData.Description, refilData.Technologies, refilData.Image)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "Charset=utf-8")
	w.Header().Set("Cache-Control", "no-chace, no-store, must-revalidate")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id=$1", id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func countDuration(start time.Time, end time.Time) string {
	timeDifferent := float64(end.Sub(start).Milliseconds())

	monthDistance := int(math.Floor(timeDifferent / (30 * 24 * 60 * 60 * 1000)))
	weekDistance := int(math.Floor(timeDifferent / (7 * 24 * 60 * 60 * 1000)))
	dayDistance := int(math.Floor(timeDifferent / (24 * 60 * 60 * 1000)))

	if monthDistance > 0 {
		str := strconv.Itoa(monthDistance) + " month"
		if monthDistance > 1 {
			return str + "s"
		}
		return str
	}
	if weekDistance > 0 {
		str := strconv.Itoa(weekDistance) + " week"
		if weekDistance > 1 {
			return str + "s"
		}
		return str
	}
	if dayDistance > 0 {
		str := strconv.Itoa(dayDistance) + " day"
		if dayDistance > 1 {
			return str + "s"
		}
		return str
	}
	return "cannot get duration"
}

func formupdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "Charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	templ, err := template.ParseFiles("html/update.html")

	Update := Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1 ;", id).Scan(&Update.Id, &Update.Name, &Update.Start_date, &Update.End_date, &Update.Description, &Update.Technologies, &Update.Image)

	str := Update.Start_date.Format("2006-01-02")
	strend := Update.End_date.Format("2006-01-02")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var resp = map[string]interface{}{
		"Data":    Data,
		"Update":  Update,
		"Datestr": str,
		"Date":    strend,
		// "Gambar" : gambar,
	}
	templ.Execute(w, resp)

}

func update(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		return
	}
	Projectname := r.PostForm.Get("projectname")
	StartDate := r.PostForm.Get("startDate")
	EndDate := r.PostForm.Get("endDate")
	Description := r.PostForm.Get("description")
	technologi1 := r.PostForm.Get("technologi1")
	technologi2 := r.PostForm.Get("technologi2")
	technologi3 := r.PostForm.Get("technologi3")
	technologi4 := r.PostForm.Get("technologi4")
	id := r.PostForm.Get("id")
	dataContext := r.Context().Value("dataFile")
	Image := dataContext.(string)

	// start, _ := time.Parse("2006-01-02", StartDate)
	// end, _ := time.Parse("2006-01-02", EndDate)

	Technologis := []string{
		technologi1,
		technologi2,
		technologi3,
		technologi4,
	}

	_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET  name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6 WHERE id=$7;", Projectname, StartDate, EndDate, Description, Technologis, Image, id)

	fmt.Println(id)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}
