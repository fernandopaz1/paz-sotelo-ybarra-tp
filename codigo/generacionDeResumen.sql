create or replace function generacion_de_resumen(nroClient int ,anio int, m int) returns void as $$
    declare
        client record;
        tarj  record;
        termTarj int;
        varCierre record;
        total float = 0;
	v record;
	numResumen int;
        
    begin
     
        select * into client from cliente cl where cl.nroCliente = nroClient;
        
        select * into tarj from tarjeta t where t.nroCliente = nroClient and t.estado = '{"v","i","g","e","n","t","e"}';
        
        termTarj = tarj.nroTarjeta[16]::int;
        
        select * into varCierre from cierre c where anio = c.año and m = c.mes and termTarj = c.terminacion;
        
        select sum(monto) into total from   compra where tarj.nroTarjeta = nroTarjeta;
			
			insert into cabecera values ( default, client.nombre, client.apellido, client.domicilio,
        tarj.nroTarjeta, varCierre.fechaInicio, varCierre.fechaCierre, varCierre.fechaVto,total);
        
        select nroResumen into numResumen from cabecera where tarj.nroTarjeta = nroTarjeta and varCierre.fechaInicio = desde;
        
        for v in select * from compra com,comercio comer where comer.nroComercio = com.nroComercio and
								  com.nroTarjeta = tarj.nroTarjeta loop
				
		insert into detalle values (numResumen, default, v.fecha::date, v.nombre, v.monto);
			
		end loop;
			
			
        ALTER SEQUENCE detalle_nroLinea_seq RESTART WITH 1;
     
    end; 
$$ language plpgsql;
