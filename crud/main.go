package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

//Hace referencia a la carperta y todos sus archivos
var plantilla = template.Must(template.ParseGlob("templates/*"))

//Datos del empleado
type Empleado struct {
	Id            int
	Nombre, Email string
}

func main() {

	//Rutas
	http.HandleFunc("/", inicio)
	http.HandleFunc("/crear", crear)
	http.HandleFunc("/insertar", insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Println("Servidor iniciado...")
	http.ListenAndServe(":8080", nil)

}

//Sirve para hacer la conexión con la base de datos
func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasena := "root"
	NombreDB := "crud"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasena+"@tcp(127.0.0.1)/"+NombreDB)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

//Página principal, muestra todos los usuarios
func inicio(rw http.ResponseWriter, r *http.Request) {

	empleado := Empleado{}
	var ListaEmpleados = []Empleado{}

	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")
	if err != nil {
		panic(err.Error())
	}

	for registros.Next() {
		var id int
		var nombre, email string
		err := registros.Scan(&id, &nombre, &email)
		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Email = email

		ListaEmpleados = append(ListaEmpleados, empleado)
	}

	plantilla.ExecuteTemplate(rw, "inicio", ListaEmpleados)
}

//Sirve para crear nuevos usuarios
func crear(rw http.ResponseWriter, r *http.Request) {
	plantilla.ExecuteTemplate(rw, "crear", nil)
}

//Recibe datos de la función crear para insertar el usuario en la BD
func insertar(rw http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		insercion, err := conexionEstablecida.Prepare("INSERT INTO empleados (Nombre, Correo) VALUES (?,?)")
		if err != nil {
			panic(err.Error())
		}

		insercion.Exec(nombre, correo)

		http.Redirect(rw, r, "/", http.StatusFound)
	}
}

//Borra un usuario a partir del id
func Borrar(rw http.ResponseWriter, r *http.Request) {

	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionBD()
	eliminacion, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE Id=?")
	if err != nil {
		panic(err.Error())
	}

	eliminacion.Exec(idEmpleado)

	http.Redirect(rw, r, "/", http.StatusFound)
}

//Página para editar un usario existente
func Editar(rw http.ResponseWriter, r *http.Request) {

	empleado := Empleado{}

	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE Id=?", idEmpleado)
	if err != nil {
		panic(err.Error())
	}

	for registro.Next() {
		var id int
		var nombre, email string

		err := registro.Scan(&id, &nombre, &email)
		if err != nil {
			panic(err.Error())
		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Email = email
	}
	plantilla.ExecuteTemplate(rw, "editar", empleado)
}

//Recibe datos de la funcion Editar para así poder actulizar los datos del usuario
func Actualizar(rw http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		email := r.FormValue("correo")

		conexionEstablecida := conexionBD()

		modificacion, err := conexionEstablecida.Prepare("UPDATE empleados SET Nombre=?, Correo=? WHERE Id=?")

		if err != nil {
			panic(err.Error())
		}

		modificacion.Exec(nombre, email, id)
		http.Redirect(rw, r, "/", http.StatusFound)
	}
}
