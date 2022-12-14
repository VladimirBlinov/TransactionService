CREATE DATABASE "TransactionDev"
    WITH
    OWNER = admin
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;

CREATE DATABASE "TransactionTest"
    WITH
    OWNER = admin
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;

CREATE DATABASE "Transaction"
    WITH
    OWNER = admin
    ENCODING = 'UTF8'
    CONNECTION LIMIT = -1;

GRANT ALL PRIVILEGES ON DATABASE "TransactionDev" TO admin;
GRANT ALL PRIVILEGES ON DATABASE "TransactionTest" TO admin;
GRANT ALL PRIVILEGES ON DATABASE "Transaction" TO admin;


