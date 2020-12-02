create or replace function crear_alerta_compras()  returns trigger as $$
	declare
		tiempoCompra  interval := interval '1' minute;
		tiempoCompraCinco  interval := interval '5' minute;
		cantComprasUnMin int;
		cantComprasCincoMin int;
		cp text;
		
	begin
		
		select codPostal into cp from comercio where nroComercio = new.nroComercio;
		
		
		select count (codPostal) into cantComprasUnMin from comercio where nroComercio in 
			(select distinct nroComercio from compra c where c.nroTarjeta = new.nroTarjeta and 
															c.fecha - new.fecha < tiempoCompra );

		
		select count (distinct codPostal) into cantComprasCincoMin from comercio where nroComercio in 
			(select distinct nroComercio from compra c where c.nroTarjeta = new.nroTarjeta and
															c.fecha - new.fecha < tiempoCompraCinco);

		
		
		if (cantComprasCincoMin > 1) then
			insert into alerta values(default, new.nroTarjeta, new.fecha , null , 5,'Se detectaron mas de 1 compra en 5 minutos');
			
			return new;
		end if;  	
		
		if cantComprasUnMin > 1 then
			
			insert into alerta values(default, new.nroTarjeta, new.fecha , null , 1,'Se detectaron mas de 1 compra en 1 minuto');
			return new;
		
		end if ;	
		
		
		return new;
		
	end; 
$$ language plpgsql;
create trigger alerta_automatica_compras_trg
after insert on compra
for each row
execute procedure crear_alerta_compras();
