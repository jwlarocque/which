# Which?

The beginnings of a simple voting/polling web app.  Svelte frontend, Go backend, PostgreSQL database.  

Current code quality: gradually improving

### Installation Instructions

1. Install `git`, `npm`, `go`, and `postgresql`
1. Get server dependencies: 
    `go get github.com/jackc/pgx github.com/jmoiron/sqlx golang.org/x/oauth2 cloud.google.com/go/compute/metadata`
1. Clone this respository: 
    `git clone https://github.com/jwlarocque/which.git`
1. Change to which directory: 
    `cd which`
1. Create a PostgreSQL DB from `schema.sql` (you might want to change the owner username)
1. Edit `start.sh.sample` with the appropriate paths and variables and rename it `start.sh`.
1. If necessary, allow execution: 
    `chmod u+x start.sh`
1. Build the server executable: 
    `go build -o which_server server`
1. Give `which_server` permission to bind reserved ports: 
    `sudo setcap 'cap_net_bind_service=+ep' which`

Run as daemon:
1. Edit `which_server.service.sample` with the appropriate paths and rename it `which_server.service`
1. Move the systemd service file: 
    `sudo cp which_server.service /etc/systemd/system`
1. Start the service 
    `systemctl start which_server.service`

Or just run `./start.sh`

### Global notes to self

* add a subtitle/place for question authors to add more information

### References

[Golang package docs](https://golang.org/pkg)  
[MDN](https://developer.mozilla.org/en-US/docs/Web)

##### Auth
[Skarlso, Google sign-in Part 1](https://skarlso.github.io/2016/06/12/google-signin-with-go/)  
[Skarlso, Part 2](https://skarlso.github.io/2016/11/02/google-signin-with-go-part2/)  
[Alex Pliutau, stdlib OAuth2](https://itnext.io/getting-started-with-oauth2-in-go-1c692420e03)  
[Jon Calhoun, Securing cookies in Go](https://www.calhoun.io/securing-cookies-in-go/)  

##### Svelte
[Svelte docs](https://svelte.dev/tutorial)  

##### SQL
[PostegreSQL docs](https://www.postgresql.org/docs/12)  
[jmoiron/sqlx](https://github.com/jmoiron/sqlx)  
[VividCortex, DB-driven Go](https://www.vividcortex.com/hubfs/eBooks/The_Ultimate_Guide_To_Building_Database-Driven_Apps_with_Go.pdf)  