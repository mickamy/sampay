CREATE EXTENSION IF NOT EXISTS dblink;

DO
$$
    DECLARE
        db_name         text := current_setting('sampay.db_name', false);
        writer_username text := current_setting('sampay.writer_username', false);
        error_message   text;
    BEGIN
        IF db_name IS NULL THEN
            RAISE EXCEPTION 'Environment variable for database name (sampay.db_name) is not set.';
        END IF;
        IF writer_username IS NULL THEN
            RAISE EXCEPTION 'Environment variable for writer username is not set.';
        END IF;

        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = db_name) THEN
                PERFORM dblink_exec('dbname=' || current_database(),
                                    'CREATE DATABASE ' || quote_ident(db_name));
                RAISE NOTICE 'Database "%" created successfully.', db_name;
            ELSE
                RAISE NOTICE 'Database "%" already exists. Skipping creation.', db_name;
            END IF;
        EXCEPTION
            WHEN insufficient_privilege THEN
                RAISE EXCEPTION 'Insufficient privileges to create database "%". Please check your permissions.', db_name;
            WHEN duplicate_database THEN
                RAISE NOTICE 'Database "%" already exists (race condition detected). Skipping creation.', db_name;
            WHEN others THEN
                GET STACKED DIAGNOSTICS error_message = MESSAGE_TEXT;
                RAISE EXCEPTION 'Error creating database "%": %', db_name, error_message;
        END;

        BEGIN
            PERFORM dblink_exec('dbname=' || quote_ident(db_name),
                                format('
                GRANT USAGE ON SCHEMA public TO %I;
                GRANT CREATE ON SCHEMA public TO %I;
                ALTER DEFAULT PRIVILEGES IN SCHEMA public
                    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO %I;
                GRANT ALL ON SCHEMA public TO %I;
            ', writer_username, writer_username, writer_username, writer_username)
                    );
            RAISE NOTICE 'Privileges granted to writer role "%" in database "%".', writer_username, db_name;
        EXCEPTION
            WHEN others THEN
                GET STACKED DIAGNOSTICS error_message = MESSAGE_TEXT;
                RAISE EXCEPTION 'Error granting privileges to writer role "%" in database "%": %', writer_username, db_name, error_message;
        END;

    EXCEPTION
        WHEN others THEN
            GET STACKED DIAGNOSTICS error_message = MESSAGE_TEXT;
            RAISE EXCEPTION 'Unexpected error occurred: %', error_message;
    END
$$
;
