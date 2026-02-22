CREATE USER sampay_writer WITH PASSWORD 'password';
CREATE USER sampay_reader WITH PASSWORD 'password';

GRANT USAGE ON SCHEMA public TO
    sampay_writer,
    sampay_reader;

\echo '✅ Successfully created users.'
