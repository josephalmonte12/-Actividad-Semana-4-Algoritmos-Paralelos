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

	// Reset autoincrement
	_, err = db.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
	if err != nil {
		return err
	}

	return nil
}

func getUsers() []*User {
	db, err := sql.Open("mysql", "joseph:1192948@tcp(db:3306)/paralelo_db")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}

	var users []*User
	for results.Next() {
		var u User
		err = results.Scan(&u.ID, &u.Name)
		if err != nil {
			panic(err.Error())
		}

		users = append(users, &u)
	}

	return users
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "joseph" && password == "1192948" {
			http.Redirect(w, r, "/users", http.StatusFound)
			return
		} else {
			w.Write([]byte("<p style='color:red'>Credenciales incorrectas. Inténtalo de nuevo.</p>"))
		}
	}

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
                    <input type="submit" value="Iniciar sesión">
               
					</form>
					</div>
				</body>
				</html>
			`))
}

func userPage(w http.ResponseWriter, r *http.Request) {
	users := getUsers()

	html := "<!DOCTYPE html><html><head><title>Lista de usuarios</title></head><body>"
	html += "<h2>Lista de usuarios</h2>"
	html += "<table border='1'><tr><th>ID</th><th>Nombre</th><th>Acción</th></tr>"

	for _, user := range users {
		html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td><form action='/update' method='post'><input type='hidden' name='id' value='%d'><input type='text' name='name' placeholder='Nuevo nombre' required><input type='submit' value='Actualizar'></form><form action='/delete' method='post'><input type='hidden' name='id' value='%d'><input type='submit' value='Eliminar'></form></td></tr>", user.ID, user.Name, user.ID, user.ID)
	}

	html += "</table>"
	html += "<br>"
	html += "<button onclick=\"window.location.href='/create'\">Crear nuevo usuario</button>"

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func createUserPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")

		err := createUser(name)
		if err != nil {
			http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/users", http.StatusFound)
		return
	}

	w.Write([]byte(`
					<html>
					<head>
						<title>Crear nuevo usuario</title> <!-- Aquí cerré correctamente la etiqueta <title> -->
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
							input[type="text"] {
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
							<h2>Crear nuevo usuario</h2>
							<form method="post">
								<label for="name">Nombre:</label>
								<input type="text" id="name" name="name" required>
								<input type="submit" value="Crear usuario">
							</form>
						</div>
					</body>
					</html>
				`))
}

func deleteUserPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de usuario no válido", http.StatusBadRequest)
			return
		}

		err = deleteUser(id)
		if err != nil {
			http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/users", http.StatusFound)
		return
	}

	w.Write([]byte(`
					<html>
					<head>
						<title>Eliminar usuario</title>
					</head>
					<body>
						<h2>Eliminar usuario</h2>
						<form action="/delete" method="post">
							<label for="id">ID del usuario a eliminar:</label>
							<input type="number" id="id" name="id" required>
							<input type="submit" value="Eliminar">
						</form>
					</body>
					</html>
				`))
}

func updateUserName(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de usuario no válido", http.StatusBadRequest)
			return
		}

		newName := r.FormValue("name")

		err = updateUser(id, newName)
		if err != nil {
			http.Error(w, "Error al actualizar el nombre del usuario", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/users", http.StatusFound)
		return
	}

	http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", userPage)
	http.HandleFunc("/create", createUserPage)
	http.HandleFunc("/delete", deleteUserPage)
	http.HandleFunc("/update", updateUserName) // Agregado handler para actualizar nombre de usuario
	log.Fatal(http.ListenAndServe(":5001", nil))
}
