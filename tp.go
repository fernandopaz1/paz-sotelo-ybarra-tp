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
        fmt.Println("Error al abrir la base de datos")
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

func cargarDatos() {
    createDatabase()
    db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
    if err != nil {
        fmt.Println("Error al abrir la base de datos ya creada")
        log.Fatal(err)
    }
    defer db.Close()
    
    tabla := `create table cliente (nroCliente int ,nombre text, apellido  text, domicilio text, telefono char[]);
			create table comercio (nroComercio int ,nombre text, domicilio text, codPostal text, telefono char[]);
			create table tarjeta (nroTarjeta char[] ,nroCliente int, valDesde char[], valHasta char[], codigoSeguridad char[], limiteCompra float, estado char[]);
			create table compra (nroOperacion int ,nroTarjeta char[], nroComercio int, fecha date, monto float, pagado boolean);
			create table rechazo (nroRechazo int, nroTarjeta char[], nroComercio int, fecha date, monto float, motivo text);
			create table cierre (año int, mes int, terminacion int, fechaInicio date, fechaCierre date, fechaVto date);
			create table cabecera (nroResumen int, nombre text, apellido text, dommicilio text, nroTarjeta char[], desde date, hasta date, vence date, total float);
			create table detalle (nroResumen int, nroLinea int, fecha date, nombreComercio text, monto float);
			create table alerta (nroAlerta int, nroTarjeta char[], fecha date, nroRechazo int, codAlerta int, descripcion text)`

    _, err = db.Exec(tabla)
    if err != nil {
        fmt.Println("Error al crear las tablas")
        log.Fatal(err)
    }
    
    pk := `alter table cliente add constraint cliente_pk primary key (nroCliente);
			alter table comercio add constraint comercio_pk primary key (nroComercio);
			alter table tarjeta add constraint tarjeta_pk primary key (nroTarjeta);
			alter table compra add constraint compra_pk primary key (nroOperacion);
			alter table rechazo add constraint rechazo_pk primary key (nroRechazo);
			alter table cierre add constraint cierre_pk primary key (año,mes,terminacion);
			alter table cabecera add constraint cabecera_pk primary key (nroResumen);
			alter table detalle add constraint detalle_pk primary key (nroResumen,nroLinea);
			alter table alerta add constraint alerta_pk primary key (nroAlerta)`
			
	_, err = db.Exec(pk)
    if err != nil {
        fmt.Println("Error al cargar las pk")
        log.Fatal(err)
    }		
    
    fk := `alter table tarjeta add constraint tarjeta_fk foreign key (nroCliente) references cliente (nroCliente);
			alter table compra add constraint compra_nroTarjeta_fk foreign key (nroTarjeta) references tarjeta (nroTarjeta);
			alter table compra add constraint compra_nroComercio_fk foreign key (nroComercio) references comercio (nroComercio);
			alter table compra add constraint rechazo_nroTarjeta_fk foreign key (nroTarjeta) references tarjeta (nroTarjeta);
			alter table compra add constraint rechazo_nroComercio_fk foreign key (nroComercio) references comercio (nroComercio);
			alter table cabecera add constraint cabecera_fk foreign key (nroTarjeta) references tarjeta (nroTarjeta);
			alter table detalle add constraint detalle_fk foreign key (nroResumen) references cabecera (nroResumen);
			alter table alerta add constraint alerta_nroTarjeta_fk foreign key (nroTarjeta) references tarjeta (nroTarjeta);
			alter table alerta add constraint alerta_nroRechazo_fk foreign key (nroRechazo) references rechazo (nroRechazo)`
	
	_, err = db.Exec(fk)
    if err != nil {
        fmt.Println("Error al cargar las fk")
        log.Fatal(err)
    }		
    		
    datosClientes := `insert into cliente values ('1','Fernando','Paz','callao 345','{"1","1","3","4","5","6","8","7","6","5"}');
						insert into cliente values ('2','Fer','Paz','calla 345','{"1","1","3","4","5","6","8","7","6","5"}')`
    
	_, err = db.Exec(datosClientes)
    if err != nil {
        fmt.Println("Error al cargar clientes")
        log.Fatal(err)
    }
    
    datosComercios := `insert into comercio values ('1','Rey del pancho','Pte peron 3221','1613','{"2","2","3","4","5","6","8","7","6","5"}')`
    
    _, err = db.Exec(datosComercios)
    if err != nil {
        fmt.Println("Error al cargar los comercios")
        log.Fatal(err)
    }
    datosTarjetas := `insert into tarjeta values ('{"1","1","3","4","5","6","8","7","6","5","5","6","8","7","6","5"}','2','{"5","6","8","7","6","5"}','{"5","6","8","7","6","4"}','{"9","6","8","7"}','121212.12','{"v","i","g","e","n","t","e"}')`
   
	_, err = db.Exec(datosTarjetas)
	if err != nil {
        fmt.Println("Error al cargar las tarjetas")
        log.Fatal(err)
    }
}
	
func main (){
	cargarDatos()
}  

