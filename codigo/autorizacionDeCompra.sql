create or replace function 
    autorizacion_de_compra(nro_tarj char[], nro_comercio int, fecha date, monto float, pagado boolean)  
    returns boolean as $$
    declare
        aceptado boolean = true;
        f_validez  char[];
        f_vencimiento date;
    begin
		
		select val_hasta into f_validez from tarjeta t where t.nro_tarjeta = nro_tarj;
        select into f_vencimiento array_de_char_a_date(f_validez);
        
		if exists (select * from tarjeta t where t.nro_tarjeta = nro_tarj and t.estado = '{"s","u","s","p","e","n","d","i","d","a"}') then
            insert into rechazo values (
                default, nro_tarj, nro_comercio, fecha, monto, 'la tarjeta se encuentra suspendida');
                aceptado = false;
            return aceptado;
        end if;
        
        if not exists(
            select * from tarjeta t where t.nro_tarjeta = nro_tarj and t.estado = '{"v","i","g","e","n","t","e"}') then
            insert into rechazo values (
                default, nro_tarj, nro_comercio, fecha, monto, 'tarjeta no válida ó no vigente');
            aceptado = false;
            return aceptado;
        end if;
				
        
        if not exists(
            select * from tarjeta t,consumo c
                where t.nro_tarjeta = nro_tarj and c.codigo_seguridad = t.codigo_seguridad and nro_tarj = c.nro_tarjeta) and
                    exists (select * from consumo c2 where c2.nro_tarjeta = nro_tarj) then
						insert into rechazo values (
						default, nro_tarj, nro_comercio, fecha, monto, 'código de seguridad inválido');
                        aceptado = false;
            return aceptado;
        end if;
        
        if f_vencimiento < fecha then
            insert into rechazo values (
                default, nro_tarj, nro_comercio, fecha, monto, 'plazo de vigencia expirado');
            aceptado = false;
            return aceptado;
        end if;
        
        if exists (select * from tarjeta t where t.nro_tarjeta = nro_tarj and t.limite_compra < monto ) then
            insert into rechazo values (
                default, nro_tarj, nro_comercio, fecha, monto, 'supera límite de tarjeta');
            aceptado = false;
            return aceptado;
        end if;      

        if aceptado then
            insert into compra values (default ,nro_tarj, nro_comercio , fecha, monto , pagado);
        end if;
    return aceptado;
    end; 
$$ language plpgsql;
