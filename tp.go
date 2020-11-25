package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
)

func createDatabase() {
    db,err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
    if err != nil {
        log.Fatal(err)
        fmt.Println("Error al abir la base de datos")
    }
    defer db.Close()

    _ , err = db.Exec(`drop database if exists transacciones;`)
    if err != nil {
        log.Fatal(err)
        fmt.Println("Error al eliminar la base si ya existia")
    }

    _, err = db.Exec(`create database transacciones;`)
    if err != nil {
        log.Fatal(err)
        fmt.Println("Error al crear la base transacciones")
    }
}

type cliente struct {
    nroCliente int
    nombre, apellido, domicilio string
    //telefono [12] rune
}

func insertarClientes(nroCliente int ,nombre string, apellido string, domicilio string) {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
    if err != nil {
        fmt.Println("Error al abir la base de datos ya creada")
        log.Fatal(err)
    }
    defer db.Close()
	salida := fmt.Sprintf("insert into cliente values (%v,%v,%v,%v);",nroCliente,nombre,apellido,domicilio)
	_, err = db.Exec(salida)
     if err != nil {
         log.Fatal(err)
     }
}

func main() {
    createDatabase()
    db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
    if err != nil {
        fmt.Println("Error al abir la base de datos ya creada")
        log.Fatal(err)
    }

    _, err = db.Exec(`create table cliente (nroCliente int ,nombre text, apellido  text, domicilio text)`)
    if err != nil {
        fmt.Println("Error al crear cliente")
        log.Fatal(err)
    }
	_, err = db.Exec(`insert into cliente values ('1','Fenando', 'Paz', 'Callefalsa');`)
    if err != nil {
        fmt.Println("Error al insertar datos en cliente")
        log.Fatal(err)
    }
	db.Close()
	insertarClientes(2,'Nacho','Sotelo','Callemasfalsa')
}
