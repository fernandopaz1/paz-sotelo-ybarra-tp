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
    
    tabla := `create table cliente (nroCliente int ,nombre text, apellido  text, domicilio text, telefono char[]);
			create table comercio (nroComercio int ,nombre text, domicilio text, codPostal text, telefono char[]);
			create table tarjeta (nroTarjeta char[] ,nroCliente int, valDesde char[], valHasta char[], codigoSeguridad char[], limiteCompra float, estado char[]);
			create table compra (nroOperacion int ,nroTarjeta char[], nroComercio int, fecha timestamp, monto float, pagado boolean);
			create table rechazo (nroRechazo int, nroTarjeta char[], nroComercio int, fecha date, monto float, motivo text);
			create table cierre (año int, mes int, terminacion int, fechaInicio date, fechaCierre date, fechaVto date);
			create table cabecera (nroResumen int, nombre text, apellido text, dommicilio text, nroTarjeta char[], desde date, hasta date, vence date, total float);
			create table detalle (nroResumen int, nroLinea int, fecha date, nombreComercio text, monto float);
			create table alerta (nroAlerta int, nroTarjeta char[], fecha date, nroRechazo int, codAlerta int, descripcion text);
			create table consumo (nroTarjeta char[], codigoSeguridad char[], nroComercio int, monto float)`

    _, err = db.Exec(tabla)
    if err != nil {
        fmt.Println("Error al crear las tablas")
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
    
    datosComercios := `insert into comercio values ('1','Adidas','Pte peron 3221','1643','{"1","1","4","9","2","1","1","9","7","1"}');
						insert into comercio values ('2','Nike','Miraflores','1643','{"1","1","4","4","5","1","8","7","6","5"}');
						insert into comercio values ('3','Mc Donals','French 231','1643','{"1","1","4","4","1","1","0","9","6","5"}');
						insert into comercio values ('4','Burger King','Av Almafuerte 436','1643','{"1","1","4","4","4","3","0","8","2","5"}');
						insert into comercio values ('5','Compumundo','Guido Spano 2534','1643','{"1","1","4","4","5","5","3","2","2","3"}');
						insert into comercio values ('6','Garbarino','Parana 3771','1642','{"1","1","4","5","5","6","3","3","6","8"}');
						insert into comercio values ('7','Musimundo','Callao 3245','1642','{"1","1","4","3","0","2","8","7","3","5"}');
						insert into comercio values ('8','Fravega','Falucho 5411','1642','{"1","1","4","1","3","4","8","6","3","0"}');
						insert into comercio values ('9','Rodo','Av Corrientes','1642','{"1","1","4","4","4","2","1","7","1","8"}');
						insert into comercio values ('10','Samsung','Callao','1642','{"1","1","4","8","3","5","1","1","6","5"}');
						insert into comercio values ('11','Freddo','Pte peron 2121','1640','{"1","1","4","1","4","1","8","0","0","1"}');
						insert into comercio values ('12','Mostaza','Ugarte 1212','1640','{"1","1","4","2","6","6","7","7","3","1"}');
						insert into comercio values ('13','Green Eat','Haiti 3367','1640','{"1","1","4","4","5","6","1","7","3","5"}');
						insert into comercio values ('14','Starbucks','Pte peron 1299','1640','{"1","1","4","4","5","2","2","5","6","3"}');
						insert into comercio values ('15','Wendy','Palpa 782','1640','{"1","1","4","3","5","6","7","7","6","9"}');
						insert into comercio values ('16','Bowen','Zelarrayan 485','1638','{"1","1","4","4","5","6","4","2","1","6"}');
						insert into comercio values ('17','Cristobal Colon','Baigorria 1513','1638','{"1","1","4","2","5","6","2","7","1","5"}');
						insert into comercio values ('18','Falabella','Pte peron 1576','1638','{"1","1","4","4","5","2","8","7","6","5"}');
						insert into comercio values ('19','Carrefour','Ugarte 3221','1638','{"1","1","4","4","1","8","6","1","3","5"}');
						insert into comercio values ('20','Etiqueta','Nazca 2356','1638','{"1","1","4","2","5","6","8","1","1","2"}')`
    
    _, err = db.Exec(datosComercios)
    if err != nil {
        fmt.Println("Error al cargar los comercios")
        log.Fatal(err)
    }
    datosTarjetas := `insert into tarjeta values ('{"5","1","5","4","5","6","8","7","6","5","5","6","8","7","6","5"}','1','{"2","0","1","1","0","6"}','{"2","0","2","7","0","6"}','{"9","6","8","7"}','60000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"5","4","2","2","5","6","8","1","6","2","5","3","8","7","6","5"}','2','{"2","0","1","2","0","2"}','{"2","0","2","8","0","2"}','{"2","4","9","2"}','70000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"5","5","3","4","5","6","4","7","3","3","5","6","8","5","5","1"}','3','{"2","0","1","3","0","1"}','{"2","0","2","9","0","2"}','{"4","4","8","2"}','70000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"5","2","3","4","4","4","8","8","6","8","5","2","2","7","1","1"}','4','{"2","0","1","4","0","2"}','{"2","0","2","2","0","1"}','{"2","6","6","3"}','80000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"5","3","3","2","5","9","8","9","6","3","3","6","1","7","6","2"}','5','{"2","0","1","3","0","4"}','{"2","0","2","1","0","1"}','{"2","1","6","3"}','60000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"5","1","5","3","5","5","8","2","6","2","5","3","8","7","6","3"}','6','{"2","0","1","2","0","6"}','{"2","0","2","4","0","6"}','{"3","1","5","5"}','40000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"5","9","9","4","5","6","3","7","3","5","3","6","2","3","6","5"}','7','{"2","0","1","1","0","5"}','{"2","0","2","2","0","6"}','{"8","2","5","5"}','70000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"5","5","2","4","5","6","8","7","6","3","5","2","8","8","8","3"}','8','{"2","0","1","0","0","2"}','{"2","0","2","4","0","8"}','{"7","2","4","7"}','60000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"5","3","1","4","5","7","7","7","6","5","2","6","8","4","6","4"}','9','{"2","0","1","2","0","3"}','{"2","0","2","2","0","8"}','{"9","6","3","6"}','80000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"5","5","0","4","5","6","8","7","6","2","2","6","2","2","6","5"}','10','{"2","0","1","3","0","6"}','{"2","0","2","2","0","3"}','{"2","6","8","7"}','90000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","4","3","4","5","6","8","7","6","5","5","6","8","7","6","1"}','11','{"2","0","1","4","0","2"}','{"2","0","2","2","0","3"}','{"2","3","2","8"}','90000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","7","3","4","2","6","8","6","6","5","3","6","8","2","2","5"}','12','{"2","0","1","2","0","2"}','{"2","0","2","3","0","2"}','{"4","3","8","8"}','20000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","1","0","4","4","6","2","2","6","5","5","6","8","1","1","1"}','13','{"2","0","1","1","0","1"}','{"2","0","2","3","0","2"}','{"6","4","2","4"}','50000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","9","3","2","2","6","1","7","6","1","5","6","1","7","6","9"}','14','{"2","0","1","0","0","1"}','{"2","0","2","2","0","1"}','{"6","4","3","4"}','60000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","1","4","4","5","6","1","1","6","5","2","2","8","7","6","2"}','15','{"2","0","1","1","0","6"}','{"2","0","2","1","0","1"}','{"5","6","3","4"}','80000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","4","0","4","5","6","8","7","6","5","5","6","4","4","6","5"}','16','{"2","0","1","4","0","8"}','{"2","0","2","2","0","5"}','{"4","5","5","2"}','80000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","1","3","4","5","6","8","7","6","5","5","6","8","7","6","5"}','17','{"2","0","1","3","0","8"}','{"2","0","2","3","0","5"}','{"3","6","5","2"}','90000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","4","9","7","7","6","7","7","6","5","8","8","8","7","1","2"}','18','{"2","0","1","2","0","7"}','{"2","0","2","4","0","9"}','{"1","5","6","4"}','80000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","2","3","4","5","6","7","7","7","5","5","9","8","9","6","3"}','19','{"2","0","1","5","0","7"}','{"2","0","2","5","0","7"}','{"2","6","6","7"}','60000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","5","3","4","5","6","8","5","4","5","5","6","8","1","1","1"}','20','{"2","0","1","2","0","6"}','{"2","0","2","2","0","7"}','{"3","7","8","5"}','50000.00','{"v","i","g","e","n","t","e"}');
						insert into tarjeta values ('{"4","5","0","4","5","6","8","3","3","5","5","6","0","7","6","0"}','19','{"2","0","0","8","0","6"}','{"2","0","1","4","0","6"}','{"4","7","6","6"}','40000.00','{"a","n","u","l","a","d","a"}');
                        insert into tarjeta values ('{"4","0","3","4","1","6","1","7","6","5","2","2","8","0","6","5"}','20','{"2","0","0","9","0","6"}','{"2","0","1","5","0","6"}','{"5","6","8","7"}','40000.00','{"a","n","u","l","a","d","a"}');
                        insert into tarjeta values ('{"4","0","5","4","1","6","1","7","6","5","2","2","8","0","6","5"}','20','{"2","0","0","9","0","6"}','{"2","0","1","5","0","6"}','{"5","6","8","8"}','40000.00','{"s","u","s","p","e","n","d","i","d","a"}')`

	_, err = db.Exec(datosTarjetas)
	if err != nil {
        fmt.Println("Error al cargar las tarjetas")
        log.Fatal(err)
    }
}


func main (){
    //menu()

	cargarDatos()
	cargarPkYFK();
	
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
    if err != nil {
        fmt.Println("Error al abrir la base de datos ya creada")
        log.Fatal(err)
    }
    defer db.Close()
    
    funcPrintEnSQL :=
    `create function print(pk int) returns void as $$
    begin 
        raise notice 'Se ejecuto el trigger con pk % ', pk;
    end; 
    $$ language plpgsql;` 

    _, err = db.Exec(funcPrintEnSQL)
    if err != nil {
        fmt.Println("Error al cargar triggers")
        log.Fatal(err)
    }    

    arrayDeCharADate := `create or replace function 
        array_de_char_a_date(venc char[])  
        returns date as $$
            declare
                anio int;
                mes  int;
                result record;
            begin
                anio = venc[1]::int * 1000 + venc[2]::int * 100 + venc[3]::int * 10 +venc[4]::int;
                mes = venc[5]::int * 10 + venc[6]::int;
                select into result format('%s-%s-%s', anio, mes, 1)::date;

                raise notice 'Esta es la fecha que le paso % ', result;
                return result;
            end;
    $$ language plpgsql;`

    _, err = db.Exec(arrayDeCharADate)
	if err != nil {
        fmt.Println("Error al cargar la funcion generacion de resumen")
        log.Fatal(err)
    }

    autDeCompraFunc := 
    `create or replace function 
        autorizacion_de_compra(nroOperacion int ,nroTarj char[], nroComercio int, fecha date, monto float, pagado boolean)  
        returns boolean as $$
            declare
                aceptado boolean = true;
                f_validez  char[];
                f_vencimiento date;
            begin
                if not exists(
                    select * from tarjeta t where t.nroTarjeta = nroTarj and t.estado = '{"v","i","g","e","n","t","e"}') then
                    insert into rechazo values (
                        nroOperacion, nroTarj, nroComercio, fecha, monto, 'tarjeta no válida ó no vigente');
                    aceptado = false;
                    return aceptado;
                end if;
					
                
                if not exists(
                    select * from tarjeta t,consumo c
                        where t.nroTarjeta = nroTarj and c.codigoSeguridad = t.codigoSeguridad and nroTarj = c.nroTarjeta) and
                        exists (select * from consumo c2 where c2.nroTarjeta = nroTarj) then
							insert into rechazo values (
							nroOperacion, nroTarj, nroComercio, fecha, monto, 'código de seguridad inválido');
                    aceptado = false;
                    return aceptado;
                end if;
                
                if exists (select * from tarjeta t where t.nroTarjeta = nroTarj and t.limiteCompra < monto ) then
                    insert into rechazo values (
                        nroOperacion, nroTarj, nroComercio, fecha, monto, 'supera límite de tarjeta');
                    aceptado = false;
                    return aceptado;
                end if;

                select valHasta into f_validez from tarjeta t where t.nroTarjeta = nroTarj;
                select into f_vencimiento array_de_char_a_date(f_validez);
                
                if f_vencimiento < fecha then
                    insert into rechazo values (
                    nroOperacion, nroTarj, nroComercio, fecha, monto, 'plazo de vigencia expirado');
                    raise notice 'corrio el if de fecha de vencimient';
                    aceptado = false;
                    return aceptado;
                end if;
                
                if exists (select * from tarjeta t where t.nroTarjeta = nroTarj and t.estado = '{"s","u","s","p","e","n","d","i","d","a"}') then
                    insert into rechazo values (
                        nroOperacion, nroTarj, nroComercio, fecha, monto, 'la tarjeta se encuentra suspendida');
                    aceptado = false;
                    return aceptado;
                end if;

                if aceptado then
                        insert into compra values (nroOperacion ,nroTarj, nroComercio , fecha, monto , pagado);
                end if;
            return aceptado;
        end; 
    $$ language plpgsql;`

    _, err = db.Exec(autDeCompraFunc)
    if err != nil {
        fmt.Println("Error al cargar funcion de autenticacion de compra")
        log.Fatal(err)
    }


    generacionDeResumen := 
    `create or replace function 
        generacion_de_resumen(nroClient int ,anio int, mes int)  
        returns void as $$
            declare
                client record;
                tarj  record;
            begin

                -- create table cliente (nroCliente int ,nombre text, apellido  text, domicilio text, telefono char[]);
                select * into client from cliente c where c.nroCliente = nroClient;
                select * into tarj from tarjeta t where t.nroCliente = nroClient and t.estado = '{"v","i","g","e","n","t","e"}';
                raise notice 'hola %', client.nombre;
                --insert into cierre values (anio , mes , tarj.nroTarjeta[16]::int, '2020-10-29','2020-11-29', periodo[2], periodo[2] + integer '7' );            
            
            end; 
    $$ language plpgsql;`

    _, err = db.Exec(generacionDeResumen)
	if err != nil {
        fmt.Println("Error al cargar la funcion generacion de resumen")
        log.Fatal(err)
    }

	consumos := `insert into consumo values ('{"5","1","5","4","5","6","8","7","6","5","5","6","8","7","6","5"}','{"1","1","1","1"}','1','150.50');`
	_, err = db.Exec(consumos)
	if err != nil {
        fmt.Println("Error al cargar el consumo")
        log.Fatal(err)
    }
    
    compras := `select autorizacion_de_compra ('1','{"5","1","5","4","5","6","8","7","6","5","5","6","8","7","6","5"}','1','2020-11-27','150.50','f');
                select autorizacion_de_compra ('2','{"4","0","3","4","1","6","1","7","6","5","2","2","8","0","6","5"}','3','2020-11-27','150.50','f');
                select autorizacion_de_compra ('3','{"5","5","3","4","5","6","4","7","3","3","5","6","8","5","5","1"}','3','2020-11-27','150000.50','f');
                select autorizacion_de_compra ('5','{"4","0","5","4","1","6","1","7","6","5","2","2","8","0","6","5"}','5','2020-11-27','155.50','f');
                select autorizacion_de_compra ('4','{"4","0","5","4","1","6","1","7","6","5","2","2","8","0","6","5"}','5','2020-11-27','155.50','f');
                select autorizacion_de_compra ('20','{"5","5","0","4","5","6","8","7","6","2","2","6","2","2","6","5"}','5','2050-11-27','155.50','f')`
                _, err = db.Exec(compras)
	if err != nil {
        fmt.Println("Error al cargar la compra")
        log.Fatal(err)
    }
	
    resumen := `select generacion_de_resumen ('1','2020', '11');`
    _, err = db.Exec(resumen)
	if err != nil {
        fmt.Println("Error al cargar el resumen")
        log.Fatal(err)
    }
}  


func menu(){
    fmt.Print("\033[H\033[2J") //Limpia la terminal

    fmt.Println(`Introduzca la opcion elegida :
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
        cargarPkYFK();
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

func cargarPkYFK(){
    db, err := sql.Open("postgres", "user=postgres host=localhost dbname=transacciones sslmode=disable")
    if err != nil {
        fmt.Println("Error al abrir la base de datos ya creada")
        log.Fatal(err)
    }
    defer db.Close()
    
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
}
