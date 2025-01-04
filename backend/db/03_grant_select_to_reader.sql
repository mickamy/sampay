\c sampay

DO
$$
    DECLARE
        r RECORD;
    BEGIN
        FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public')
            LOOP
                EXECUTE format('GRANT SELECT ON TABLE public.%I TO sampay_reader;', r.tablename);
            END LOOP;
    END
$$;
