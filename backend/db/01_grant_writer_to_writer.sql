DO
$$
    DECLARE
        writer_username TEXT := current_setting('sampay.writer_username', false);
    BEGIN
        IF writer_username IS NULL THEN
            RAISE EXCEPTION 'Environment variable for writer username is not set.';
        END IF;

        EXECUTE format('
        GRANT USAGE, CREATE ON SCHEMA public TO %I;
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO %I;
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO %I;
        GRANT ALL PRIVILEGES ON SCHEMA public TO %I;
    ', writer_username, writer_username, writer_username, writer_username);

        RAISE NOTICE 'Privileges granted to writer role "%s" in database "%s".', writer_username, current_database();
    END
$$;
