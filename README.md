# webapp
This is a repository contains RESTful API developed in Go lang as a part of CSYE6225, Network Structures and Cloud Computing course at Northeastern University, Boston, MA

### Prerequisite and Installations
Before you get started, ensure you have the following prerequisites installed on your system:
- [Go](https://golang.org/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)

### Code and Configurations
1. Clone this repository:

   ```bash
   git clone git@github.com:cloud-computing-csye6225/webapp.git
   ```

2. Install required code dependencies using,

   ```bash
   go mod download
   ```

3. Set up .env file with required configurations, below is the format for .env

   ```bash
    DBHOST= <Host Address>
    DBUSER= <DB Username>
    DBPASSWORD= <DB Password>
    DBNAME= <DB name>
    DBPORT= <DB PORT>
    SERVERPORT= <API Server port, :8080 is default>
    DEFAULTUSERS= <Path to the default users file>
   ```
   
4. Start Postgres server if it is not active, and create required DB users & password if does not exists
5. Run the application
