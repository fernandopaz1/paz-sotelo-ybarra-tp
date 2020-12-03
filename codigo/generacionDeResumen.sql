create or replace function generacion_de_resumen(nro_client int ,anio int, m int) returns void as $$
    declare
        client record;
        tarj  record;
        term_tarj int;
        cierre_actual record;
        total float = 0;
	    v record;
	    num_resumen int;
        
    begin
    
        select * into client from cliente cl where cl.nro_cliente = nro_client;
        
        select * into tarj from tarjeta t where t.nro_cliente = nro_client;
        
        term_tarj = tarj.nro_tarjeta[16]::int;
        
        select * into cierre_actual from cierre c where anio = c.año and m = c.mes and term_tarj = c.terminacion;
        
        select coalesce(sum(monto),0) into total from   compra where tarj.nro_tarjeta = nro_tarjeta and
                                                fecha::date > cierre_actual.fecha_inicio and
                                                fecha::date < cierre_actual.fecha_cierre;


        insert into cabecera values ( default, client.nombre, client.apellido, client.domicilio,
        tarj.nro_tarjeta, cierre_actual.fecha_inicio, cierre_actual.fecha_cierre, cierre_actual.fecha_vto,total);
        
        select nro_resumen into num_resumen from cabecera where tarj.nro_tarjeta = nro_tarjeta and cierre_actual.fecha_inicio = desde;
        
        for v in select * from compra com,comercio comer where comer.nro_comercio = com.nro_comercio and
                    com.nro_tarjeta = tarj.nro_tarjeta and
                                            fecha::date > cierre_actual.fecha_inicio and
                                            fecha::date < cierre_actual.fecha_cierre loop
        
            insert into detalle values (num_resumen, default, v.fecha::date, v.nombre, v.monto);

        end loop;
        
        update compra set pagado = true where nro_tarjeta = tarj.nro_tarjeta and
                                            fecha::date > cierre_actual.fecha_inicio and
                                            fecha::date < cierre_actual.fecha_cierre;	
            
        ALTER SEQUENCE detalle_nro_linea_seq RESTART WITH 1;
    end; 
$$ language plpgsql;
