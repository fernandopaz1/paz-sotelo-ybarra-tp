create or replace function crear_alerta()  returns trigger as $$
	begin
		insert into alerta values(default, new.nroTarjeta, new.fecha , new.nroRechazo, 0, new.motivo);
		return new;
	end; 
$$ language plpgsql;

create trigger alerta_automatica_trg
after insert on rechazo
for each row
execute procedure crear_alerta();
