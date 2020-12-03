create or replace function array_de_char_a_date(venc char[]) returns date as $$
    declare
        anio int;
        mes  int;
        result record;
    begin
        anio = venc[1]::int * 1000 + venc[2]::int * 100 + venc[3]::int * 10 +venc[4]::int;
        mes = venc[5]::int * 10 + venc[6]::int;
        select into result format('%s-%s-%s', anio, mes, 1)::date;

        raise notice 'Esta es la fecha que le paso % ', result;
        return result;
    end;
$$ language plpgsql;


create or replace function cargar_consumos_en_compra() returns void as $$
    declare
        v record;
    begin
        for v in select * from consumo loop
            perform autorizacion_de_compra  (v.nro_tarjeta, v.nro_comercio, '2020-11-15',v.monto, 'f');
        end loop;
    end;
$$ language plpgsql;


create or replace function generar_resumenes_del_anio() returns void as $$
    declare
        v record;
    begin
        for v in select * from cliente loop
            for m in 1..12 loop
				perform generacion_de_resumen (v.nro_cliente,'2020',m);
				end loop;
        end loop;
    end;
$$ language plpgsql;

