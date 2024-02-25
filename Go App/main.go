package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func createUser(name string) error {
	db, err := sql.Open("mysql", "joseph:1192948@tcp(db:3306)/paralelo_db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (name) VALUES (?)", name)
	if err != nil {
		return err
	}

	return nil
}

// Función para actualizar un usuario
func updateUser(id int, name string) error {
	db, err := sql.Open("mysql", "joseph:1192948@tcp(db:3306)/paralelo_db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return err
	}

	return nil
}

// Función para eliminar un usuario
func deleteUser(id int) error {
	db, err := sql.Open("mysql", "joseph:1192948@tcp(db:3306)/paralelo_db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func getUsers() []*User {
	// Abra nuestra conexión de base de datos.
	db, err := sql.Open("mysql", "joseph:1192948@tcp(db:3306)/paralelo_db")

	// si hay un error al abrir la conexión, manejarlo
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	// Ejecutar la consulta
	results, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error()) // manejo adecuado de errores en lugar de pánico en su aplicación
	}

	var users []*User
	for results.Next() {
		var u User
		// para cada fila, escanee el resultado en nuestra etiqueta de objeto compuesto
		err = results.Scan(&u.ID, &u.Name)
		if err != nil {
			panic(err.Error()) // manejo adecuado de errores en lugar de pánico en su aplicación
		}

		users = append(users, &u)
	}

	return users
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// Comprobamos si el método de la solicitud es POST
	if r.Method == http.MethodPost {
		// Obtenemos los valores del formulario de inicio de sesión
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Verificamos las credenciales del usuario (aquí podrías hacer la verificación con la base de datos)
		if username == "joseph" && password == "1192948" {
			// Si las credenciales son válidas, redirigimos al usuario a la página de usuarios
			http.Redirect(w, r, "/users", http.StatusFound)
			return
		} else {
			// Si las credenciales son inválidas, mostramos un mensaje de error
			w.Write([]byte("<p style='color:red'>Credenciales incorrectas. Inténtalo de nuevo.</p>"))
		}
	}

	// Renderizamos el formulario de inicio de sesión
	w.Write([]byte(`
        <html>
        <head>
            <title>Página de inicio de sesión</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    margin: 0;
                    padding: 0;
                    background-color: #f4f4f4;
                }
                .container {
                    max-width: 400px;
                    margin: 100px auto;
                    padding: 20px;
                    background-color: #fff;
                    border-radius: 5px;
                    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
                }
                input[type="text"], input[type="password"] {
                    width: 100%;
                    padding: 10px;
                    margin: 5px 0;
                    border: 1px solid #ccc;
                    border-radius: 4px;
                    box-sizing: border-box;
                }
                input[type="submit"] {
                    width: 100%;
                    background-color: #4caf50;
                    color: white;
                    padding: 10px;
                    margin: 10px 0;
                    border: none;
                    border-radius: 4px;
                    cursor: pointer;
                }
                input[type="submit"]:hover {
                    background-color: #45a049;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h2>Iniciar sesión</h2>
                <form method="post">
                    <label for="username">Nombre de usuario:</label>
                    <input type="text" id="username" name="username" required>
                    <label for="password">Contraseña:</label>
                    <input type="password" id="password" name="password" required>
                    <!-- Agregar un campo de entrada para el nombre del usuario -->
                    <label for="name">Nombre:</label>
                    <input type="text" id="name" name="name" required>
                    <input type="submit" value="Iniciar sesión">
                </form>
            </div>
        </body>
        </html>
    `))
}

func userPage(w http.ResponseWriter, r *http.Request) {
	// Obtener usuarios de la base de datos
	users := getUsers()

	// Construir una tabla HTML para mostrar los usuarios
	html := "<!DOCTYPE html><html><head><title>Lista de usuarios</title></head><body>"
	html += "<h2>Lista de usuarios</h2>"
	html += "<table border='1'><tr><th>ID</th><th>Nombre</th></tr>"

	// Iterar sobre los usuarios y agregar filas a la tabla
	for _, user := range users {
		html += fmt.Sprintf("<tr><td>%d</td><td>%s</td></tr>", user.ID, user.Name)
	}

	html += "</table></body></html>"

	// Escribir la respuesta HTTP con la tabla de usuarios
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el nombre del usuario del cuerpo de la solicitud
	name := r.FormValue("name")

	// Crear el usuario en la base de datos
	err := createUser(name)
	if err != nil {
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}

	// Redirigir al usuario a la página de usuarios
	http.Redirect(w, r, "/users", http.StatusFound)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID y el nuevo nombre del usuario del cuerpo de la solicitud
	id := r.FormValue("id")
	name := r.FormValue("name")

	// Convertir el ID a un entero
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID de usuario no válido", http.StatusBadRequest)
		return
	}

	// Actualizar el usuario en la base de datos
	err = updateUser(userID, name)
	if err != nil {
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	// Redirigir al usuario a la página de usuarios
	http.Redirect(w, r, "/users", http.StatusFound)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario a eliminar del cuerpo de la solicitud
	id := r.FormValue("id")

	// Convertir el ID a un entero
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID de usuario no válido", http.StatusBadRequest)
		return
	}

	// Eliminar el usuario de la base de datos
	err = deleteUser(userID)
	if err != nil {
		http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
		return
	}

	// Redirigir al usuario a la página de usuarios
	http.Redirect(w, r, "/users", http.StatusFound)
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", userPage)
	http.HandleFunc("/create", createUserHandler) // Handler para crear usuario
	http.HandleFunc("/update", updateUserHandler) // Handler para actualizar usuario
	http.HandleFunc("/delete", deleteUserHandler) // Handler para eliminar usuario
	log.Fatal(http.ListenAndServe(":5001", nil))
}
