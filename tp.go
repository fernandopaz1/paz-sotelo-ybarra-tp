package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
    "time"
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
    telefono [12] rune
}

type comercio struct {
    nroComercio int
    nombre, domicilio, codPostal string
    telefono [12] rune
}

type tarjeta struct {
    nroCliente int
    nroTarjeta [16] rune
    validaDesde, validaHasta [6] rune
    codigoSeguridad [4] rune
    estado [10] rune
    limiteCompra float64
}

type compra struct {
    nroOperacion,nroComercio int
    nroTarjeta [16] rune
    fecha time.Time
    monto float64
    pagado bool
}

func insertarCompra(db *sql.DB, nroOperacion int ,nroTarjeta string, nroComercio int, fecha string, monto float64, pagado bool){
	comandoSQL := fmt.Sprintf(`insert into compra values ('%v','%v','%v','%v','%v','%v');`,nroOperacion,stringtoSQLArray(nroTarjeta),nroComercio,fecha,monto,pagado)
   fmt.Println(comandoSQL)
   fmt.Println(fecha)
   var err error
    _, err = db.Exec(comandoSQL)
    if err != nil {
        log.Fatal(err)
    }
}

func insertarTarjeta(db *sql.DB, numeroTarjeta string ,nroCliente int, valDesde string, valHasta string, codigoSeguridad string, limiteCompra float64, estadO string){
	nroTarjeta := stringtoSQLArray(numeroTarjeta)
	validaDesde := stringtoSQLArray(valDesde)
	validaHasta := stringtoSQLArray(valHasta)
	codSeguridad := stringtoSQLArray(codigoSeguridad)
	estado := stringtoSQLArray(estadO)
	
	comandoSQL := fmt.Sprintf(`insert into tarjeta values ('%v','%v','%v','%v','%v','%v','%v');`,
	nroTarjeta,nroCliente,validaDesde,validaHasta,codSeguridad,limiteCompra,estado)
   
   var err error
    _, err = db.Exec(comandoSQL)
    if err != nil {
        log.Fatal(err)
    }
}

func insertarComercio(db *sql.DB, nroComercio int ,nombre string, domicilio string, codPostal string, tel string){
	comandoSQL := fmt.Sprintf(`insert into comercio values ('%v','%v','%v','%v','%v');`,nroComercio , nombre,domicilio,codPostal,stringtoSQLArray(tel))
   var err error
    _, err = db.Exec(comandoSQL)
    if err != nil {
        log.Fatal(err)
    }
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
    
    _, err = db.Exec(`create table comercio (nroComercio int ,nombre text, domicilio text, codPostal text, telefono char[])`)
    if err != nil {
        fmt.Println("Error al crear comercio")
        log.Fatal(err)
    }
    
    _, err = db.Exec(`create table tarjeta (numeroTarjeta char[] ,nroCliente int, valDesde char[], valHasta char[], codigoSeguridad char[], limiteCompra float, estado char[])`)
    if err != nil {
        fmt.Println("Error al crear tarjeta")
        log.Fatal(err)
    }
    _, err = db.Exec(`create table compra (nroOperacion int ,nroTarjeta char[], nroComercio int, fecha date, monto float, pagado boolean)`)
    if err != nil {
        fmt.Println("Error al crear compra")
        log.Fatal(err)
    }

    tel :="111111111111"
    telComercio := "222222222222"
    layout := "2006-01-02"
    updatedAt, _ := time.Parse(layout, "2016-06-10") // lee un string y lo transforma a formato fecha
    fechita :=updatedAt.Format("2006-01-02") // acorta la fecha sin min y zona horaria 
    
    //var limCompra float64 
    //limCompra = 30000.50
    
    insertarClientes(db,1,"Fernando","Paz","Calle falsa",tel)

    insertarClientes(db,2,"Nacho","Sotelo","Callemasfalsa",tel)
    
    insertarClientes(db,3,"Flavio","Ybarra","Calle 123",tel)
    
    insertarClientes(db,4,"Florentina","Sosa","Fonrouge 3870",tel)
    
    insertarClientes(db,5,"Verónica","Roldán","Callemasfalsa",tel)
    
    
    insertarComercio(db,1,"Rey del pancho","Pte Peron 222","1613",telComercio)
    
    insertarTarjeta(db,"1234567891011123",1,"202005","202205","1234",30000.50,"suspendida")
    
    insertarCompra(db,1,"1234567891011123",3156,fechita,122.2,true)
    
    
}

func stringtoSQLArray(s string) string{
    nuevo :=`{`
    for i:=0; i<len(s)-1; i++ {
        nuevo += `"`+string(s[i])+`",`
    }
    nuevo+= `"`+string(s[len(s)-1])+`"}`
    return nuevo
}
