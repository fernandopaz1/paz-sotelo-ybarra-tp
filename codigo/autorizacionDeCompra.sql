create or replace function 
    autorizacion_de_compra(nroTarj char[], nroComercio int, fecha date, monto float, pagado boolean)  
    returns boolean as $$
    declare
        aceptado boolean = true;
        f_validez  char[];
        f_vencimiento date;
    begin
		
		select valHasta into f_validez from tarjeta t where t.nroTarjeta = nroTarj;
        select into f_vencimiento array_de_char_a_date(f_validez);
        
		if exists (select * from tarjeta t where t.nroTarjeta = nroTarj and t.estado = '{"s","u","s","p","e","n","d","i","d","a"}') then
            insert into rechazo values (
                default, nroTarj, nroComercio, fecha, monto, 'la tarjeta se encuentra suspendida');
                aceptado = false;
            return aceptado;
        end if;
        
        if not exists(
            select * from tarjeta t where t.nroTarjeta = nroTarj and t.estado = '{"v","i","g","e","n","t","e"}') then
            insert into rechazo values (
                default, nroTarj, nroComercio, fecha, monto, 'tarjeta no válida ó no vigente');
            aceptado = false;
            return aceptado;
        end if;
				
        
        if not exists(
            select * from tarjeta t,consumo c
                where t.nroTarjeta = nroTarj and c.codigoSeguridad = t.codigoSeguridad and nroTarj = c.nroTarjeta) and
                    exists (select * from consumo c2 where c2.nroTarjeta = nroTarj) then
						insert into rechazo values (
						default, nroTarj, nroComercio, fecha, monto, 'código de seguridad inválido');
                        aceptado = false;
            return aceptado;
        end if;
        
        if f_vencimiento < fecha then
            insert into rechazo values (
                default, nroTarj, nroComercio, fecha, monto, 'plazo de vigencia expirado');
            aceptado = false;
            return aceptado;
        end if;
        
        if exists (select * from tarjeta t where t.nroTarjeta = nroTarj and t.limiteCompra < monto ) then
            insert into rechazo values (
                default, nroTarj, nroComercio, fecha, monto, 'supera límite de tarjeta');
            aceptado = false;
            return aceptado;
        end if;      

        if aceptado then
            insert into compra values (default ,nroTarj, nroComercio , fecha, monto , pagado);
        end if;
    return aceptado;
    end; 
$$ language plpgsql;
