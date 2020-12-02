create table cliente (nroCliente int ,nombre text, apellido  text, domicilio text, telefono char[]);
create table comercio (nroComercio int ,nombre text, domicilio text, codPostal text, telefono char[]);
create table tarjeta (nroTarjeta char[] ,nroCliente int, valDesde char[], valHasta char[], codigoSeguridad char[], limiteCompra float, estado char[]);
create table compra (nroOperacion int ,nroTarjeta char[], nroComercio int, fecha timestamp, monto float, pagado boolean);
create table rechazo (nroRechazo serial, nroTarjeta char[], nroComercio int, fecha timestamp, monto float, motivo text);
create table cierre (año int, mes int, terminacion int, fechaInicio date, fechaCierre date, fechaVto date);
create table cabecera (nroResumen serial, nombre text, apellido text, domicilio text, nroTarjeta char[], desde date, hasta date, vence date, total float);
create table detalle (nroResumen int, nroLinea serial, fecha date, nombreComercio text, monto float);
create table alerta (nroAlerta serial, nroTarjeta char[], fecha timestamp, nroRechazo int, codAlerta int, descripcion text);
create table consumo (nroTarjeta char[], codigoSeguridad char[], nroComercio int, monto float)
