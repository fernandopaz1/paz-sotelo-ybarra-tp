create or replace function crear_alerta()  returns trigger as $$
	declare
		mismo_dia  interval := interval '1 day';
		cant_rechazo int;
	begin
		insert into alerta values(default, new.nro_tarjeta, new.fecha , new.nro_rechazo, 0, new.motivo);
		select count(*) into cant_rechazo from rechazo where nro_tarjeta = new.nro_tarjeta and fecha - new.fecha < mismo_dia and new.motivo = 'supera límite de tarjeta';
		if (cant_rechazo > 1) then
			insert into alerta values(default, new.nro_tarjeta, new.fecha , new.nro_rechazo, 32,'Tarjeta suspendida por varios excesos de limite');
			update tarjeta set estado = '{"s","u","s","p","e","n","d","i","d","a"}' where nro_tarjeta = new.nro_tarjeta;
		end if;
		return new;
	end; 
$$ language plpgsql;
create trigger alerta_automatica_trg
after insert on rechazo
for each row
execute procedure crear_alerta();
