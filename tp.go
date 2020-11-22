package main

import (
    "database/sql"
    "fmt"
    "github.com/lib/pq"
    "log"
)

func createDatabase() {
    db,err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _ , err = db.Exec(`drop database if exist transacciones`)
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec(`create database transacciones;`)
    if err != nil {
        log.Fatal(err)
    }
}

type cliente struct {
    nroCliente int
    nombre, apellido, domicilio string
    telefono [12] rune
}

func insertarClientes(db sql.DB ,nroCliente int ,nombre string, apellido string, domicilio string, telefono []rune) {
    _, err = db.Exec(`insert into cliente values (`+nroCliente+`, `+nombre+`,`+apellido+`,`+domicilio+`,`+telefono+`);`)
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    createDatabase()
    db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec(`create table cliente (nroCliente int ,nombre text, apellido  text, domicilio textg, telefono char[])`)
    if err != nil {
        log.Fatal(err)
    }

    insertarClientes(db,1 , "Fenando", "Paz", "domicilio string", {"1","1","1","1","1","1","1","1","1","1","1","1"})  
    defer db.Close()
}
