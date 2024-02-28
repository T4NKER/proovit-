# AdCash Internship assignment

## How to run

` chmod +x scripts/run.sh && bash scripts/run.sh ` 

or use  

` chmod +x scripts/databaseInit.sh && bash scripts/databaseInit.sh && go run cmd/main.go ` 

I used WSL to create the project.

### The page is running on localhost:8080.

#### Additional information

The goal was to make three separate API endpoints: one to see all the transfers, one to see the balance in EUR and BTC and one to make transfers. 

From the imported packages I used GORM and gin web framework. GORM, because it facilitates working with databases, especially when the queries are rather short and gin because it makes the errorhandling and routing of different endpoints very comfortable.

In general I followed a modular approach to make every query have a separate function and every additional check be done in the handler by helpers.
