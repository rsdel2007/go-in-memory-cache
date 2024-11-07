# Simple In-Memory Database

This project implements an in-memory key/value database similar to Redis. The database supports basic data operations and transaction commands, allowing for nested transactions.

 
## Data Commands

The database accepts the following commands to operate on keys:

* `SET name value` – Set the variable `name` to the value `value`. For
  simplicity `value` may be an integer.
* `GET name` – Print out the value of the variable `name`, or `NULL` if that
  variable is not set.
* `UNSET name` – Unset the variable `name`, making it just like that variable
  was never set.

```
INPUT	            OUTPUT
--------------------------
SET ex 10
GET ex              10
UNSET ex
GET ex              NULL


INPUT	            OUTPUT
--------------------------
SET b 10
SET b 30
GET b               30
```

## Transaction Commands

In addition to the above data commands, the database also supports transactions with the following commands:

* `BEGIN` – Open a new transaction block. **Transactions can be nested;** a
  `BEGIN` can be issued inside of an existing block.
* `ROLLBACK` – Undo commands issued in the current transaction, and closes it
  Returns an error if no transaction is in progress.
* `COMMIT` – Close **all** open transactions, permanently applying the changes
  made in them. Returns an error if no transaction is in progress.

Any data command that is run outside of a transaction should commit
immediately. Here are some example command sequences:

```
INPUT	          OUTPUT
------------------------
BEGIN
SET a 10
GET a             10
BEGIN
SET a 20
GET a             20
ROLLBACK
GET a             10
ROLLBACK
GET a             NULL
END

INPUT	          OUTPUT
------------------------
BEGIN
SET a 30
BEGIN
SET a 40
COMMIT
GET a             40
ROLLBACK          NO TRANSACTION
END


INPUT	          OUTPUT
------------------------
SET a 50
BEGIN
GET a             50
SET a 60
BEGIN
UNSET a
GET a             NULL
ROLLBACK
GET a             60
COMMIT
GET a             60
END
```

## Implementaion
The database is implemented in Go, with the following key components:

* `SimpleDb` interface: Defines the methods for setting, getting, unsetting keys, and managing transactions.
* `simpleDB` struct: Implements the `SimpleDb` interface, using a map to store data and a slice of maps to manage transactions.
