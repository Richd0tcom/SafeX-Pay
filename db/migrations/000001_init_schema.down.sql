DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS transfers;

/*drop those two first because theres a foreign key constraint that references the accounts table*/
DROP TABLE IF EXISTS accounts;