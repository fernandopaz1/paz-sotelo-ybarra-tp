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
