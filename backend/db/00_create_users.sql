DO
$do$
    DECLARE
        writer_username text := current_setting('sampay.writer_username', false);
        writer_password text := current_setting('sampay.writer_password', false);
        reader_username text := current_setting('sampay.reader_username', false);
        reader_password text := current_setting('sampay.reader_password', false);
        error_message   text;
    BEGIN
        IF writer_username IS NULL THEN
            RAISE EXCEPTION 'Environment variable for writer username is not set.';
        END IF;
        IF writer_password IS NULL THEN
            RAISE EXCEPTION 'Environment variable for writer password is not set.';
        END IF;
        IF reader_username IS NULL THEN
            RAISE EXCEPTION 'Environment variable for reader username is not set.';
        END IF;
        IF reader_password IS NULL THEN
            RAISE EXCEPTION 'Environment variable for reader password is not set.';
        END IF;

        BEGIN
            EXECUTE format('CREATE ROLE %I WITH LOGIN PASSWORD %L', writer_username, writer_password);
            RAISE NOTICE 'Writer role "%" created successfully.', writer_username;
        EXCEPTION
            WHEN duplicate_object THEN
                RAISE NOTICE 'Writer role "%" already exists. Skipping creation.', writer_username;
        END;

        BEGIN
            EXECUTE format('
            GRANT CREATE, USAGE ON SCHEMA public TO %I;
            ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO %I;
            GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO %I;
        ', writer_username, writer_username, writer_username);
            RAISE NOTICE 'Privileges granted to writer role "%".', writer_username;
        EXCEPTION
            WHEN others THEN
                GET STACKED DIAGNOSTICS error_message = MESSAGE_TEXT;
                RAISE EXCEPTION 'Error granting privileges to writer role: %', error_message;
        END;

        BEGIN
            EXECUTE format('CREATE USER %I PASSWORD %L', reader_username, reader_password);
            RAISE NOTICE 'Reader role "%" created successfully.', reader_username;
        EXCEPTION
            WHEN duplicate_object THEN
                RAISE NOTICE 'Reader role "%" already exists. Skipping creation.', reader_username;
        END;

        BEGIN
            EXECUTE format('
            GRANT USAGE ON SCHEMA public TO %I;
            GRANT SELECT ON ALL TABLES IN SCHEMA public TO %I;
            ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO %I;
        ', reader_username, reader_username, reader_username);
            RAISE NOTICE 'Privileges granted to reader role "%".', reader_username;
        EXCEPTION
            WHEN others THEN
                GET STACKED DIAGNOSTICS error_message = MESSAGE_TEXT;
                RAISE EXCEPTION 'Error granting privileges to reader role: %', error_message;
        END;

        BEGIN
            EXECUTE format('
            ALTER DEFAULT PRIVILEGES FOR ROLE %I IN SCHEMA public
                GRANT SELECT ON TABLES TO %I;
        ', writer_username, reader_username);
            RAISE NOTICE 'Default privileges set for writer role "%" to grant SELECT to reader role "%".', writer_username, reader_username;
        EXCEPTION
            WHEN others THEN
                GET STACKED DIAGNOSTICS error_message = MESSAGE_TEXT;
                RAISE EXCEPTION 'Error setting default privileges: %', error_message;
        END;

    EXCEPTION
        WHEN others THEN
            GET STACKED DIAGNOSTICS error_message = MESSAGE_TEXT;
            RAISE EXCEPTION 'Unexpected error occurred: %', error_message;
    END
$do$;
