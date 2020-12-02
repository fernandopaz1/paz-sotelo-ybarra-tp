alter table cliente add constraint cliente_pk primary key (nroCliente);
alter table comercio add constraint comercio_pk primary key (nroComercio);
alter table tarjeta add constraint tarjeta_pk primary key (nroTarjeta);
alter table compra add constraint compra_pk primary key (nroOperacion);
alter table rechazo add constraint rechazo_pk primary key (nroRechazo);
alter table cierre add constraint cierre_pk primary key (año,mes,terminacion);
alter table cabecera add constraint cabecera_pk primary key (nroResumen);
alter table detalle add constraint detalle_pk primary key (nroResumen,nroLinea);
alter table alerta add constraint alerta_pk primary key (nroAlerta)