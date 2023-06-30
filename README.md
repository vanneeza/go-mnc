# Installation
- Clone the repository https://github.com/vanneeza/go-mnc.git
- Create a database in PostgreSQL.
- Open the folder where you have cloned the repository. Navigate to ../utils/document/sql
- Import ddl.sql and dml.sql into the previously created database.
- Make sure to import ddl.sql first, followed by dml.sql.
- Import the project into VSCode or a similar application.
- Create a file named ".env". You can refer to the sample.env file as an example.
- Type "go mod tidy" in the terminal to ensure that all packages are installed.
- You can also use "go build" to make sure that all packages are installed.
- Type "go run main.go" to run the server.
- If you are using "go build", you can open the executable file (go-mnc.exe).
- Open Postman or a similar application.
- Import the Postman collection file located in the folder ../utils/document/api/.

# System Flow
1. Merchant registers, logs in, adds products, and adds banks.
2. Customer registers, logs in, and performs transactions such as creating orders and making payments.
3. Admin confirms payments, and the transactions are completed.

# Requirements
- PostgreSQL installed.
- Application like VSCode.
- Application like Postman.

# Note
For complete documentation, please refer to the following file: utils/document/api/Dokumentasi.pdf