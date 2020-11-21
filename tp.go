package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
)

type alumne struct {
    legajo           int
    nombre, apellido string
}

func createDatabase() {

    db,err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

	_, err = db.Exec(`drop database if exists guarani`)
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec(`create database guarani`)
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    createDatabase()

    db, err := sql.Open("postgres", "user=postgres host=localhost dbname=guarani sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(`create table alumne (legajo int, nombre text, apellido text)`)
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec(`insert into alumne values (1, 'Cristina', 'Kirchner');
                      insert into alumne values (2, 'Juan Domingo', 'Perón');`)

    if err != nil {
        log.Fatal(err)
	}
	
	 rows, err := db.Query(`select * from alumne where nombre = 'Cristina'`)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    var a alumne
    for rows.Next() {
        if err := rows.Scan(&a.legajo, &a.nombre, &a.apellido); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%v %v %v\n", a.legajo, a.nombre, a.apellido)
    }
    if err = rows.Err(); err != nil {
        log.Fatal(err)
	}

}

	