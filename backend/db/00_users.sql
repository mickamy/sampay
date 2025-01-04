-- for development purposes only
CREATE OR REPLACE FUNCTION terminate_backend(pid integer)
    RETURNS void
    LANGUAGE plpgsql
    SECURITY DEFINER AS
$$
BEGIN
    PERFORM pg_terminate_backend(pid);
END;
$$;

DO
$do$
    BEGIN
        IF EXISTS (SELECT usename FROM pg_user WHERE usename = 'sampay_writer') THEN

            RAISE NOTICE 'Role "sampay_writer" already exists. Skipping.';
        ELSE
            CREATE USER sampay_writer CREATEDB PASSWORD 'password';

            GRANT CREATE ON SCHEMA public TO sampay_writer;
            GRANT USAGE ON SCHEMA public TO sampay_writer;
            ALTER DEFAULT PRIVILEGES IN SCHEMA public
                GRANT ALL ON TABLES TO sampay_writer;
            GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO sampay_writer;

            -- for development purposes only
            GRANT SELECT ON pg_stat_activity TO sampay_writer;
            GRANT EXECUTE ON FUNCTION terminate_backend(integer) TO sampay_writer;
        END IF;
    END
$do$;

DO
$do$
    BEGIN
        IF EXISTS (SELECT usename FROM pg_user WHERE usename = 'sampay_reader') THEN
            RAISE NOTICE 'Role "sampay_reader" already exists. Skipping.';
        ELSE
            CREATE USER sampay_reader PASSWORD 'password';

            GRANT USAGE ON SCHEMA public TO sampay_reader;
            GRANT SELECT ON ALL TABLES IN SCHEMA public TO sampay_reader;
            ALTER DEFAULT PRIVILEGES IN SCHEMA public
                GRANT SELECT ON TABLES TO sampay_reader;
        END IF;
    END
$do$;

ALTER DEFAULT PRIVILEGES FOR ROLE sampay_writer IN SCHEMA public
    GRANT SELECT ON TABLES TO sampay_reader;
ALTER DEFAULT PRIVILEGES FOR ROLE sampay_writer IN SCHEMA public
    GRANT SELECT ON TABLES TO sampay_reader;
