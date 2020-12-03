create table cliente (nro_cliente int ,nombre text, apellido  text, domicilio text, telefono char[]);
create table comercio (nro_comercio int ,nombre text, domicilio text, cod_postal text, telefono char[]);
create table tarjeta (nro_tarjeta char[] ,nro_cliente int, val_desde char[], val_hasta char[], codigo_seguridad char[], limite_compra float, estado char[]);
create table compra (nro_operacion serial ,nro_tarjeta char[], nro_comercio int, fecha timestamp, monto float, pagado boolean);
create table rechazo (nro_rechazo serial, nro_tarjeta char[], nro_comercio int, fecha timestamp, monto float, motivo text);
create table cierre (año int, mes int, terminacion int, fecha_inicio date, fecha_cierre date, fecha_vto date);
create table cabecera (nro_resumen serial, nombre text, apellido text, domicilio text, nro_tarjeta char[], desde date, hasta date, vence date, total float);
create table detalle (nro_resumen int, nro_linea serial, fecha date, nombre_comercio text, monto float);
create table alerta (nro_alerta serial, nro_tarjeta char[], fecha timestamp, nro_rechazo int, cod_alerta int, descripcion text);
create table consumo (nro_tarjeta char[], codigo_seguridad char[], nro_comercio int, monto float)
