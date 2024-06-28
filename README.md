# wheres-waldo-api

This is an API I made using Go and PostGreSQL with database hosting on Supabase. This API serves a browser game I made based on _Where's Waldo?_ I chose to use Go for this project because I wanted to try making backend services with a language besides JavaScript.

Refer to [wheres-waldo](https://github.com/ken-ux/wheres-waldo) for the frontend repository.

## Technology Used

- Go
  - Gin: Library used as a web framework but I mostly utilized for easier routing.
  - pgx/pgxpool: Library used as PostGreSQL driver.
- PostGreSQL
- Supabase

## Database Schema

I made a simple relational database for the API to interact with. The statements for creating the tables were executed through Supabase's web interface but I wrote them down in the `creates_tables.sql` and `insert_into_tables.sql` files in this repository for future reference by myself and others. And here's the schema I drew out beforehand:

![Relational database schema.](/database_schema.png)

## Lessons Learned

- Binding incoming form data from a POST request to use in queries.
- Connecting Go server to instance of a Supabase database.
- Correct usage of double vs. single quotes when writing PostGreSQL queries.
- A small issue I ran into early was an access error whenever I tried to connect to my database.
  - The syntax for http authentication is `user:password@domain.com`. I discovered that the admin password for the database contained an "@" symbol, which led to issues identifying the correct domain.
- Another issue I had was an error when querying the database multiple times. The error referenced prepared statements already existing.
  - I read through Supabase documentation and saw that my API was connecting to the database in **transaction** mode rather than **session** mode. Supabase provides two different pooling modes to determine how connections are assigned to clients. Note that prepared statements are not available in transaction mode.
  - This became an issue because pgx, the PostGreSQL driver I used in Go, prepares and caches statements by default. Therefore, all queries being made were prepared statements that weren't compatible with the pooling mode. Connecting to the database in session mode rather than transaction mode fixed my issue.
- Struct fields need to be capitalized to be exported, otherwise they'll be empty when converted into JSON.
  - In the process of troubleshooting an issue related to this, I learned how to use the standard `encoding/json` library to convert structs into JSON format, and then translate JSON from bytes to strings. In the end, I didn't need to do any of that but it could be useful in the future, especially if I need to figure out conversion of complex structs to and from JSON later on.
- Configuring CORS to allow requests from different origins. I was familiar with this concept because of past experience with CORS configuration in Node.js.
