package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	bolt "github.com/coreos/bbolt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	menu()
}

func menu() {

	//Limpia la terminal
	fmt.Print("\033[H\033[2J")

	fmt.Println(`Introduzca la opcion elegida :
				1. Para crear la base de datos
				2. Para cargar datos
				3. Para agregar las Pk y Fk
				4. Para cargar los stored procedures y triggers
				5. Testear base usando consumo
				6. Borrar Pks  y Fks
				7. Cargar base Bolt.db
				q. Salir`)
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()

	if err != nil {
		fmt.Println(err)
	}

	switch char {
	case '1':
		fmt.Println("Creando")
		crearBase()
		time.Sleep(2 * time.Second)
		break
	case '2':
		fmt.Println("cargando la base")
		cargarDatos()
		time.Sleep(2 * time.Second)
		break
	case '3':
		fmt.Println("cargando PKs y FKs")
		cargarPkYFK()
		time.Sleep(2 * time.Second)
		break
	case '4':
		fmt.Println("cargando stored procedures")
		cargarProceduresYTriggers()
		time.Sleep(2 * time.Second)
		break
	case '5':
		fmt.Println("testando base con consumo")
		testearBaseConConsumo()
		time.Sleep(2 * time.Second)
		break
	case '6':
		fmt.Println("borrando PKs y FKs")
		borrarKeys()
		time.Sleep(2 * time.Second)
		break
	case '7':
		fmt.Println("creando BoltDB")
		crearBoltDB()
		time.Sleep(2 * time.Second)
		break

	case 'q':
		fmt.Println("Chau")
		return
		break
	default:
		fmt.Println("La opcion elegida no es valida")
		time.Sleep(2 * time.Second)
	}
	menu()
}

type Cliente struct {
	NroCliente                  int
	Nombre, Apellido, Domicilio string
	Telefono                    [12]rune
}

type Comercio struct {
	NroComercio                  int
	Nombre, Domicilio, CodPostal string
	Telefono                     [12]rune
}

type Tarjeta struct {
	NroCliente               int
	NroTarjeta               [16]rune
	ValidaDesde, ValidaHasta [6]rune
	CodigoSeguridad          [4]rune
	Estado                   [10]rune
	LimiteCompra             float64
}

type Compra struct {
	NroOperacion, NroComercio int
	NroTarjeta                [16]rune
	Fecha                     time.Time
	Monto                     float64
	Pagado                    bool
}

func crearBase() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error al abrir la base de datos")
	}
	defer db.Close()

	_, err = db.Exec(`drop database if exists transacciones;`)
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

func cargarDatos() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
	if err != nil {
		fmt.Println("Error al abrir la base de datos ya creada")
		log.Fatal(err)
	}
	defer db.Close()

	cargarComandosAPostgres(db, "codigo/crearTablas.sql")

	cargarComandosAPostgres(db, "codigo/datosClientes.sql")

	cargarComandosAPostgres(db, "codigo/datosComercios.sql")

	cargarComandosAPostgres(db, "codigo/datosTarjetas.sql")

	cargarCierre(db, 2020)

	cargarComandosAPostgres(db, "codigo/datosConsumos.sql")
}

//conviene usar una funcion en Go por el manejo de ciclos
func cargarCierre(db *sql.DB, anio int) {
	d := 1
	var fechainicio string
	var fechacierre string
	var fechavto string
	for m := 1; m < 13; m++ {
		for t := 0; t < 10; t++ {
			fechainicio = fmt.Sprintf("%v-%v-%v", anio, m, d+t)
			if m < 12 {
				fechacierre = fmt.Sprintf("%v-%v-%v", anio, m+1, d+t+1)
				fechavto = fmt.Sprintf("%v-%v-%v", anio, m+1, d+t+5)
			} else {
				fechacierre = fmt.Sprintf("%v-%v-%v", anio, m-11, d+t+1)
				fechavto = fmt.Sprintf("%v-%v-%v", anio, m-11, d+t+5)
			}
			comandoSQL := fmt.Sprintf("insert into cierre values ('%v','%v','%v','%v','%v','%v');", anio, m, t, fechainicio, fechacierre, fechavto)

			_, err := db.Exec(comandoSQL)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func cargarPkYFK() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
	if err != nil {
		fmt.Println("Error al abrir la base de datos ya creada")
		log.Fatal(err)
	}
	defer db.Close()

	cargarComandosAPostgres(db, "codigo/pks.sql")
	cargarComandosAPostgres(db, "codigo/fks.sql")

}

func cargarComandosAPostgres(db *sql.DB, path string) {
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

func cargarProceduresYTriggers() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
	if err != nil {
		fmt.Println("Error al abrir la base de datos ya creada")
		log.Fatal(err)
	}
	defer db.Close()

	cargarComandosAPostgres(db, "codigo/funcionesAuxiliares.sql")

	cargarComandosAPostgres(db, "codigo/autorizacionDeCompra.sql")

	cargarComandosAPostgres(db, "codigo/generacionDeResumen.sql")

	cargarComandosAPostgres(db, "codigo/triggerRechazo.sql")

	cargarComandosAPostgres(db, "codigo/triggerCompra.sql")

}

func testearBaseConConsumo() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
	if err != nil {
		fmt.Println("Error al abrir la base de datos ya creada")
		log.Fatal(err)
	}
	defer db.Close()

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
}

func borrarKeys() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
	if err != nil {
		fmt.Println("Error al abrir la base de datos ya creada")
		log.Fatal(err)
	}
	defer db.Close()

	cargarComandosAPostgres(db, "codigo/removeKeys.sql")
}

func crearBoltDB() {

	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// se cargan los cliente

	fernando := Cliente{1, "Fernando", "Paz", "Callao 345", [12]rune{'1', '1', '3', '4', '5', '6', '8', '7', '6', '5', '6', '5'}}
	data, err := json.Marshal(fernando)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "cliente", []byte(strconv.Itoa(fernando.NroCliente)), data)
	resultado1, err := ReadUnique(db, "cliente", []byte(strconv.Itoa(fernando.NroCliente)))
	fmt.Printf("%s\n", resultado1)

	manolo := Cliente{2, "Manolo", "Lettiere", "Matheu 3942", [12]rune{'1', '1', '4', '7', '5', '4', '4', '3', '6', '0'}}
	data, err = json.Marshal(manolo)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "cliente", []byte(strconv.Itoa(manolo.NroCliente)), data)
	resultado2, err := ReadUnique(db, "cliente", []byte(strconv.Itoa(manolo.NroCliente)))
	fmt.Printf("%s\n", resultado2)

	carlota := Cliente{3, "Carlota", "Correa", "San Martin 975", [12]rune{'1', '1', '9', '4', '4', '2', '7', '7', '3', '5'}}
	data, err = json.Marshal(carlota)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "cliente", []byte(strconv.Itoa(carlota.NroCliente)), data)
	resultado3, err := ReadUnique(db, "cliente", []byte(strconv.Itoa(carlota.NroCliente)))
	fmt.Printf("%s\n", resultado3)

	// se cargan los comercio

	adidas := Comercio{1, "Adidas", "Pte peron 3221", "1643", [12]rune{'1', '1', '4', '9', '2', '1', '1', '9', '7', '1'}}
	data, err = json.Marshal(adidas)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "comercio", []byte(strconv.Itoa(adidas.NroComercio)), data)
	resultado4, err := ReadUnique(db, "comercio", []byte(strconv.Itoa(adidas.NroComercio)))
	fmt.Printf("%s\n", resultado4)

	nike := Comercio{2, "Nike", "Miraflores 2121", "1643", [12]rune{'1', '1', '4', '4', '5', '1', '8', '7', '6', '5'}}
	data, err = json.Marshal(nike)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "comercio", []byte(strconv.Itoa(nike.NroComercio)), data)
	resultado5, err := ReadUnique(db, "comercio", []byte(strconv.Itoa(nike.NroComercio)))
	fmt.Printf("%s\n", resultado5)

	mcDonals := Comercio{3, "Mc Donals", "French 231", "1643", [12]rune{'1', '1', '4', '4', '1', '1', '0', '9', '6', '5'}}
	data, err = json.Marshal(mcDonals)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "comercio", []byte(strconv.Itoa(mcDonals.NroComercio)), data)
	resultado6, err := ReadUnique(db, "comercio", []byte(strconv.Itoa(mcDonals.NroComercio)))
	fmt.Printf("%s\n", resultado6)

	// se cargan compras

	compra1 := Compra{1, 1, [16]rune{'5', '1', '5', '4', '5', '6', '8', '7', '6', '5', '5', '6', '8', '7', '6', '5'}, stringATime("2020-11-27"), 150.50, false}
	data, err = json.Marshal(compra1)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "compra", []byte(strconv.Itoa(compra1.NroOperacion)), data)
	resultado7, err := ReadUnique(db, "compra", []byte(strconv.Itoa(compra1.NroOperacion)))
	fmt.Printf("%s\n", resultado7)

	compra2 := Compra{2, 3, [16]rune{'4', '0', '3', '4', '1', '6', '1', '7', '6', '5', '2', '2', '8', '0', '6', '5'}, stringATime("2020-11-27"), 150.50, false}
	data, err = json.Marshal(compra2)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "compra", []byte(strconv.Itoa(compra2.NroOperacion)), data)
	resultado8, err := ReadUnique(db, "compra", []byte(strconv.Itoa(compra2.NroOperacion)))
	fmt.Printf("%s\n", resultado8)

	compra3 := Compra{3, 3, [16]rune{'5', '5', '3', '4', '5', '6', '4', '7', '3', '3', '5', '6', '8', '5', '5', '1'}, stringATime("2020-11-27"), 150000.50, false}
	data, err = json.Marshal(compra3)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "compra", []byte(strconv.Itoa(compra3.NroOperacion)), data)
	resultado9, err := ReadUnique(db, "compra", []byte(strconv.Itoa(compra3.NroOperacion)))
	fmt.Printf("%s\n", resultado9)

	// cargando tarjetas

	tarjeta1 := Tarjeta{2, [16]rune{'5', '4', '2', '2', '5', '6', '8', '1', '6', '2', '5', '3', '8', '7', '6', '5'}, [6]rune{'2', '0', '1', '2', '0', '2'}, [6]rune{'2', '0', '2', '8', '0', '2'}, [4]rune{'2', '4', '9', '2'}, [10]rune{'v', 'i', 'g', 'e', 'n', 't', 'e'}, 70000.00}
	data, err = json.Marshal(tarjeta1)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "tarjeta", []byte(strconv.Itoa(tarjeta1.NroCliente)), data)
	resultado10, err := ReadUnique(db, "tarjeta", []byte(strconv.Itoa(tarjeta1.NroCliente)))
	fmt.Printf("%s\n", resultado10)

	tarjeta2 := Tarjeta{3, [16]rune{'5', '5', '3', '4', '5', '6', '4', '7', '3', '3', '5', '6', '8', '5', '5', '1'}, [6]rune{'2', '0', '1', '3', '0', '1'}, [6]rune{'2', '0', '2', '9', '0', '2'}, [4]rune{'4', '4', '8', '2'}, [10]rune{'v', 'i', 'g', 'e', 'n', 't', 'e'}, 70000.00}
	data, err = json.Marshal(tarjeta2)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "tarjeta", []byte(strconv.Itoa(tarjeta2.NroCliente)), data)
	resultado11, err := ReadUnique(db, "tarjeta", []byte(strconv.Itoa(tarjeta2.NroCliente)))
	fmt.Printf("%s\n", resultado11)

	tarjeta3 := Tarjeta{5, [16]rune{'5', '3', '3', '2', '5', '9', '8', '9', '6', '3', '3', '6', '1', '7', '6', '2'}, [6]rune{'2', '0', '1', '3', '0', '4'}, [6]rune{'2', '0', '2', '1', '0', '1'}, [4]rune{'2', '1', '6', '3'}, [10]rune{'v', 'i', 'g', 'e', 'n', 't', 'e'}, 60000.00}
	data, err = json.Marshal(tarjeta3)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(db, "tarjeta", []byte(strconv.Itoa(tarjeta3.NroCliente)), data)
	resultado12, err := ReadUnique(db, "tarjeta", []byte(strconv.Itoa(tarjeta3.NroCliente)))
	fmt.Printf("%s\n", resultado12)

}

func CreateUpdate(db *bolt.DB, bucketName string, key []byte, val []byte) error {
	// abre transacción de escritura
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))
	err = b.Put(key, val)
	if err != nil {
		return err
	}
	// cierra transacción
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func ReadUnique(db *bolt.DB, bucketName string, key []byte) ([]byte, error) {
	var buf []byte
	// abre una transacción de lectura
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		buf = b.Get(key)
		return nil
	})
	return buf, err
}

func stringATime(str string) (t time.Time) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
