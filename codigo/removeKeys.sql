alter table tarjeta drop constraint tarjeta_fk;
alter table compra drop constraint compra_nroTarjeta_fk;
alter table compra drop constraint compra_nroComercio_fk;
alter table compra drop constraint rechazo_nroTarjeta_fk;
alter table compra drop constraint rechazo_nroComercio_fk;
alter table cabecera drop constraint cabecera_fk;
alter table detalle drop constraint detalle_fk;
alter table alerta drop constraint alerta_nroTarjeta_fk;
alter table alerta drop constraint alerta_nroRechazo_fk;
alter table cliente drop constraint cliente_pk;
alter table comercio drop constraint comercio_pk;
alter table tarjeta drop constraint tarjeta_pk;
alter table compra drop constraint compra_pk;
alter table rechazo drop constraint rechazo_pk;
alter table cierre drop constraint cierre_pk;
alter table cabecera drop constraint cabecera_pk;
alter table detalle drop constraint detalle_pk;
alter table alerta drop constraint alerta_pk;