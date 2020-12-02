create or replace function crear_alerta()  returns trigger as $$
	declare
		increment_value  interval := interval '1 day';
		cantRechazo int;
	begin
		insert into alerta values(default, new.nroTarjeta, new.fecha , new.nroRechazo, 0, new.motivo);
		select count(*) into cantRechazo from rechazo where nroTarjeta = new.nroTarjeta and fecha - new.fecha < increment_value and new.motivo = 'supera límite de tarjeta';
		if (cantRechazo > 1) then
			insert into alerta values(default, new.nroTarjeta, new.fecha , new.nroRechazo, 32,'Tarjeta suspendida por varios excesos de limite');
			update tarjeta set estado = '{"s","u","s","p","e","n","d","i","d","a"}' where nroTarjeta = new.nroTarjeta;
		end if;
		return new;
	end; 
$$ language plpgsql;

create trigger alerta_automatica_trg
after insert on rechazo
for each row
execute procedure crear_alerta();
