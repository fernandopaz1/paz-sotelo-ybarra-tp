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


func insertarClientes(db *sql.DB, nroCliente int ,nombre string, apellido string, domicilio string,tel string){
	comandoSQL := fmt.Sprintf(`insert into cliente values ('%v','%v','%v','%v','%v');`,nroCliente,nombre,apellido,domicilio,stringtoSQLArray(tel))
    var err error
    _, err = db.Exec(comandoSQL)
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
    defer db.Close()

    _, err = db.Exec(`create table cliente (nroCliente int ,nombre text, apellido  text, domicilio text, telefono char[])`)
    if err != nil {
        fmt.Println("Error al crear cliente")
        log.Fatal(err)
    }
	
    tel :="111111111111"
    insertarClientes(db,1,"Fernando","Paz","Calle falsa",tel)

    insertarClientes(db,2,"Nacho","Sotelo","Callemasfalsa",tel)

}

func stringtoSQLArray(s string) string{
    nuevo :=`{`
    for i:=0; i<len(s)-1; i++ {
        nuevo += `"`+string(s[i])+`",`
    }
    nuevo+= `"`+string(s[len(s)-1])+`"}`
    return nuevo
}
