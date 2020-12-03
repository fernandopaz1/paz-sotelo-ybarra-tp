create or replace function crear_alerta_compras()  returns trigger as $$
	declare
		tiempo_compra  interval := interval '1' minute;
		tiempo_compraCinco  interval := interval '5' minute;
		cant_compras_un_min int;
		cant_compras_cinco_min int;
		cp text;
		
	begin
		
		select cod_postal into cp from comercio where nro_comercio = new.nro_comercio;
		
		
		select count (cod_postal) into cant_compras_un_min from comercio where nro_comercio in 
			(select distinct nro_comercio from compra c where c.nro_tarjeta = new.nro_tarjeta and 
															c.fecha - new.fecha < tiempo_compra );

		
		select count (distinct cod_postal) into cant_compras_cinco_min from comercio where nro_comercio in 
			(select distinct nro_comercio from compra c where c.nro_tarjeta = new.nro_tarjeta and
															c.fecha - new.fecha < tiempo_compraCinco);

		
		
		if (cant_compras_cinco_min > 1) then
			insert into alerta values(default, new.nro_tarjeta, new.fecha , null , 5,'Se detectaron mas de 1 compra en 5 minutos');			
			return new;
		end if;  	
		
		if cant_compras_un_min > 1 then
			insert into alerta values(default, new.nro_tarjeta, new.fecha , null , 1,'Se detectaron mas de 1 compra en 1 minuto');
			return new;
		
		end if ;	
		
		
		return new;
		
	end; 
$$ language plpgsql;
create trigger alerta_automatica_compras_trg
after insert on compra
for each row
execute procedure crear_alerta_compras();
