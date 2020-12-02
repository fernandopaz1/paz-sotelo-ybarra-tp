create or replace function crear_alerta_compras()  returns trigger as $$
	declare
		tiempoCompra  interval := interval '1' minute;
		tiempoCompraCinco  interval := interval '5' minute;
		cantComprasUnMin int;
		cantComprasCincoMin int;
		cp text;

	begin
		select codPostal into cp from comercio where nroComercio in (select nroComercio from compra where nroTarjeta = new.nroTarjeta);
		
		select count(*) into cantComprasUnMin from compra c,comercio com where c.nroTarjeta = new.nroTarjeta 
		and c.fecha - new.fecha < tiempoCompra and cp = com.codPostal and c.nroComercio != new.nroComercio;
		
		select count(*) into cantComprasCincoMin from compra c,comercio com where
		 c.nroTarjeta = new.nroTarjeta and c.fecha - new.fecha < tiempoCompraCinco and cp != com.codPostal;
		
		
		if (cantComprasUnMin > 2) then
			insert into alerta values(default, new.nroTarjeta, new.fecha , null , 1,'Se detectaron mas de 1 compra en 1 minuto');
		return new;	
		elsif cantComprasCincoMin > 2 then
			insert into alerta values(default, new.nroTarjeta, new.fecha , null , 5,'Se detectaron mas de 1 compra en 5 minutos');
		return new;	
		end if;
		return new;
	end; 
$$ language plpgsql;
create trigger alerta_automatica_compras_trg
after insert on compra
for each row
execute procedure crear_alerta_compras();
