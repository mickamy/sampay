DO
$$
    DECLARE
        r               RECORD;
        reader_username text := current_setting('sampay.reader_username', false);
        error_message   text;
    BEGIN
        IF reader_username IS NULL THEN
            RAISE EXCEPTION 'Environment variable for reader username is not set.';
        END IF;

        BEGIN
            EXECUTE format('GRANT USAGE ON SCHEMA public TO %I', reader_username);
            RAISE NOTICE 'Granted USAGE on schema public to "%".', reader_username;
        EXCEPTION
            WHEN others THEN
                GET STACKED DIAGNOSTICS error_message = MESSAGE_TEXT;
                RAISE EXCEPTION 'Failed to connect to database: %', error_message;
        END;

        BEGIN
            FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public')
                LOOP
                    EXECUTE format('GRANT SELECT ON TABLE public.%I TO %I', r.tablename, reader_username);
                END LOOP;

            FOR r IN (SELECT table_name FROM information_schema.views WHERE table_schema = 'public')
                LOOP
                    EXECUTE format('GRANT SELECT ON TABLE public.%I TO %I', r.table_name, reader_username);
                END LOOP;

            RAISE NOTICE 'SELECT privileges granted to reader role "%" on all existing tables and views.', reader_username;
        EXCEPTION
            WHEN others THEN
                GET STACKED DIAGNOSTICS error_message = MESSAGE_TEXT;
                RAISE EXCEPTION 'Error granting SELECT privileges to reader role "%": %', reader_username, error_message;
        END;

    EXCEPTION
        WHEN others THEN
            GET STACKED DIAGNOSTICS error_message = MESSAGE_TEXT;
            RAISE EXCEPTION 'Unexpected error occurred: %', error_message;
    END
$$
;
