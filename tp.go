package main

import (
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

    datosClientes := `insert into cliente values ('1','Fernando','Paz','Callao 345','{"1","1","3","4","5","6","8","7","6","5"}');
						insert into cliente values ('2','Manolo','Lettiere','Matheu 3942','{"1","1","4","7","5","4","4","3","6","0"}');
						insert into cliente values ('3','Carlota','Correa','San Martin 975','{"1","1","9","4","4","2","7","7","3","5"}');
						insert into cliente values ('4','Florentina','Sosa','Fonruoge 3870','{"1","1","4","6","0","2","0","6","9","6"}');
						insert into cliente values ('5','Yara','Leiva','Charcas 5128','{"1","1","4","7","7","7","7","6","3","8"}');
						insert into cliente values ('6','Cristiano','Borroni','Av Centenario 837','{"1","1","4","7","4","3","2","2","7","3"}');
						insert into cliente values ('7','Leonor','Ortiz','24 de septiembre 263','{"1","1","4","2","1","6","1","5","1","5"}');
						insert into cliente values ('8','Levina','Dellucci','Thames 550','{"1","1","4","8","5","8","0","7","6","7"}');
						insert into cliente values ('9','Salvino','Castiglione','Moreno 1785','{"1","1","4","8","8","3","5","2","9","1"}');
						insert into cliente values ('10','Franco','Cruz','Nuñez 345','{"1","1","4","5","5","4","3","5","0","0"}');
						insert into cliente values ('11','Virgilio','Angelo','Mitre 424','{"1","1","4","4","2","1","3","0","3","0"}');
						insert into cliente values ('12','Galeno','Romero','Gonzalez 461','{"1","1","4","4","3","0","8","7","9","6"}');
						insert into cliente values ('13','Rosa','Rousse','Av Centenario 743','{"1","1","4","7","4","3","8","4","9","5"}');
						insert into cliente values ('14','Agustin','Arcuri','Fotheringham 282','{"1","1","4","4","6","2","2","8","0","8"}');
						insert into cliente values ('15','Nekate','Longo','Av Besares 1170','{"1","1","4","4","2","7","8","2","6","3"}');
						insert into cliente values ('16','Ventana','Garcia','Yrigoyen 739','{"1","1","4","4","5","6","4","8","3","5"}');
						insert into cliente values ('17','Nevada','Lombardi','Boulogne Sur Mer 372','{"1","1","4","4","5","4","0","4","0","7"}');
						insert into cliente values ('18','Telma','Chavez','Av Cabildo 2370','{"1","1","4","7","8","0","5","4","1","4"}');
						insert into cliente values ('19','Augusto','Bravo','San Luis 2745','{"1","1","4","4","8","3","1","1","3","4"}');
						insert into cliente values ('20','Romano','Cocci','Calle 24 1235','{"1","1","4","4","3","0","9","2","6","1"}')`
    
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
    
    compras := `insert into compra values ('0','{"1","1","3","4","5","6","8","7","6","5","5","6","8","7","6","5"}','1','2020-11-27','150.50','t')`
    _, err = db.Exec(compras)
	if err != nil {
        fmt.Println("Error al cargar la compra")
        log.Fatal(err)
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
    
    funcPrintEnSQL :=
    `create function print() returns void as $$
    begin 
        raise notice 'Se ejecuto el trigger';
    end; 
    $$ language plpgsql;` 

    _, err = db.Exec(funcPrintEnSQL)
    if err != nil {
        fmt.Println("Error al cargar triggers")
        log.Fatal(err)
    }    

    autDeCompraFunc := 
    `create or replace function autorizacion_De_Compra()  returns trigger as $$ 
        begin 
            if new.nroTarjeta != old.nroTarjeta then
                insert into compra values (new.nroOperacion ,new.nroTarjeta, new.nroComercio, new.fecha, new.monto, new.pagado);
                select print();
            end if;
            return new;
        end; 
    $$ language plpgsql;`

    _, err = db.Exec(autDeCompraFunc)
    if err != nil {
        fmt.Println("Error al cargar triggers")
        log.Fatal(err)
    }

    autDeCompraTrigg :=
    `create trigger autorizacionCompra_trg
    before insert or update on compra
    for each row
    execute procedure autorizacion_De_Compra();`

	_, err = db.Exec(autDeCompraTrigg)
	if err != nil {
        fmt.Println("Error al cargar triggers")
        log.Fatal(err)
    }
    
    compras2 := `insert into compra values ('1','{"1","1","3","4","5","6","8","7","6","5","5","6","8","7","6","5"}','1','2020-11-27','150.50','t')`
	_, err = db.Exec(compras2)
	if err != nil {
        fmt.Println("Error al cargar la compra")
        log.Fatal(err)
    }
}  


func menu(){
    fmt.Print("\033[H\033[2J") //Limpia la terminal

    fmt.Println(`Introduzca la opcion elegida :
                1. Para crear la base de datos 
                2. Para cargar la tabla 
                3. Para verificar los stored procedures
                4. Carga los mismos datos en NoSQL
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
        fmt.Println("cargando la base")
        break
        case '3':
        fmt.Println("verificando stored procedures")
        break
        case '4':
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
