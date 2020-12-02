create or replace function generacion_de_resumen(nroClient int ,anio int, m int) returns void as $$
    declare
        client record;
        tarj  record;
        termTarj int;
        cierre_actual record;
        total float = 0;
	    v record;
	    numResumen int;
        
    begin
    
        select * into client from cliente cl where cl.nroCliente = nroClient;
        
        select * into tarj from tarjeta t where t.nroCliente = nroClient;
        
        termTarj = tarj.nroTarjeta[16]::int;
        
        select * into cierre_actual from cierre c where anio = c.año and m = c.mes and termTarj = c.terminacion;
        
        select coalesce(sum(monto),0) into total from   compra where tarj.nroTarjeta = nroTarjeta and
                                                fecha::date > cierre_actual.fechaInicio and
                                                fecha::date < cierre_actual.fechaCierre;


        insert into cabecera values ( default, client.nombre, client.apellido, client.domicilio,
        tarj.nroTarjeta, cierre_actual.fechaInicio, cierre_actual.fechaCierre, cierre_actual.fechaVto,total);
        
        select nroResumen into numResumen from cabecera where tarj.nroTarjeta = nroTarjeta and cierre_actual.fechaInicio = desde;
        
        for v in select * from compra com,comercio comer where comer.nroComercio = com.nroComercio and
                    com.nroTarjeta = tarj.nroTarjeta and
                                            fecha::date > cierre_actual.fechaInicio and
                                            fecha::date < cierre_actual.fechaCierre loop
        
            insert into detalle values (numResumen, default, v.fecha::date, v.nombre, v.monto);

        end loop;
        
        update compra set pagado = true where nroTarjeta = tarj.nroTarjeta and
                                            fecha::date > cierre_actual.fechaInicio and
                                            fecha::date < cierre_actual.fechaCierre;	
            
        ALTER SEQUENCE detalle_nroLinea_seq RESTART WITH 1;
    end; 
$$ language plpgsql;
