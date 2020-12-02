package main

import (
	"io/ioutil"
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "log"
    "time"
    "bufio"
    "os"
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
type consumo struct {
    nroTarjeta [16] rune
    codigoSeguridad [4] rune
	nroComercio int 
    monto float64
}

func cargarDatos() {
    createDatabase()
    db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
    if err != nil {
        fmt.Println("Error al abrir la base de datos ya creada")
        log.Fatal(err)
    }
    defer db.Close()
    
    cargarComandosAPostgres(db,"codigo/crearTablas.sql")

    cargarComandosAPostgres(db,"codigo/datosClientes.sql")

    cargarComandosAPostgres(db,"codigo/datosComercios.sql")

    cargarComandosAPostgres(db,"codigo/datosTarjetas.sql")
    

    cargarCierre(db,2020)
    
}

//conviene usar una funcion en Go por el manejo de ciclos
func cargarCierre(db *sql.DB,anio int){
	d := 1
	var fechainicio string
	var fechacierre string
	var fechavto string
	for m:= 1; m < 13; m++{
		for t:= 0; t< 10; t++{
			fechainicio = fmt.Sprintf("%v-%v-%v",anio,m,d+t)
			if m<12{
			fechacierre = fmt.Sprintf("%v-%v-%v",anio,m+1,d+t+1)
			fechavto = fmt.Sprintf("%v-%v-%v",anio,m+1,d+t+5)
			}else {
				fechacierre = fmt.Sprintf("%v-%v-%v",anio,m-11,d+t+1)
				fechavto = fmt.Sprintf("%v-%v-%v",anio,m-11,d+t+5)		
			}
            comandoSQL := fmt.Sprintf("insert into cierre values ('%v','%v','%v','%v','%v','%v');",anio, m, t, fechainicio, fechacierre, fechavto)
            
            _, err := db.Exec(comandoSQL)
            if err != nil {
                log.Fatal(err)
            }
		}
	}
}	


func main (){
    //menu()

	cargarDatos()
	
	
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
    if err != nil {
        fmt.Println("Error al abrir la base de datos ya creada")
        log.Fatal(err)
    }
    defer db.Close()
    
    cargarPkYFK(db)        
    cargarComandosAPostgres(db,"codigo/funcionesAuxiliares.sql")

    cargarComandosAPostgres(db,"codigo/autorizacionDeCompra.sql")

    cargarComandosAPostgres(db,"codigo/generacionDeResumen.sql")

    cargarComandosAPostgres(db,"codigo/triggerRechazo.sql")
    
    cargarComandosAPostgres(db,"codigo/triggerCompra.sql")
    
    cargarComandosAPostgres(db,"codigo/datosConsumos.sql")


    _, err = db.Exec(`select cargar_consumos_en_compra()`)
	if err != nil {
        fmt.Println("Error al cargar el consumo")
        log.Fatal(err)
    }
    
    _, err = db.Exec(`select generar_resumenes_del_anio()`)
	if err != nil {
        fmt.Println("Error al cargar el consumo")
        log.Fatal(err)
    }

    cargarComandosAPostgres(db, "codigo/removeKeys.sql")
}  


func menu(){
    fmt.Print("\033[H\033[2J") //Limpia la terminal

    fmt.Println(`Intsroduzca la opcion elegida :
                1. Para crear la base de datos
                2. Para agregar las Pk y Fk 
                3. Para cargar la tabla 
                4. Para verificar los stored procedures
                5. Carga los mismos datos en NoSQL
                q. Salir`) 
    reader := bufio.NewReader(os.Stdin)
    char, _, err := reader.ReadRune()

    if err != nil {
    fmt.Println(err)
    }

    switch char {
        case '1':
        fmt.Println("Creando")
        break
        case '2':
        fmt.Println("verificando stored procedures")
        break
        case '3':
        fmt.Println("cargando la base")
        break
        case '4':
        fmt.Println("verificando stored procedures")
        break
        case '5':
        fmt.Println("Cargando en NoSQL")
        break
        case 'q':
        fmt.Println("Chau")
        break
        default:
        fmt.Println("La opcion elegida no es valida")
        time.Sleep(2 * time.Second)
        menu()
    }
}

func cargarPkYFK(db *sql.DB){
    cargarComandosAPostgres(db, "codigo/pks.sql")
    cargarComandosAPostgres(db, "codigo/fks.sql")
    
}


func cargarComandosAPostgres(db *sql.DB, path string){
    file, err := ioutil.ReadFile(path)

    if err != nil {
        log.Fatal(err)
    }

    request := string(file)

    _, err = db.Exec(request)    
    if err != nil {
        log.Fatal(err)
    }
}
