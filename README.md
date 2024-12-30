### Auth Service

Run `make` to create and start running containers of keycloack and go-auth services

#### Register

```bash
curl --request POST \
  --url http://localhost:8081/register \
  --header 'Content-Type: application/json' \
  --data '{
	"username": "Liuba",
	"password": "*****",
	"email": "l@l"
}'
```

### Login
```bash
curl --request POST \
  --url http://localhost:8081/login \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/10.1.1' \
  --data '{
	"username": "liuba",
	"password": "*****"
}'
```

This method returns access_token and refresh_token.

### Logout
If you want to logout, you should use refresh_token from login response.

```bash
curl --request POST \
  --url http://localhost:8081/logout \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/10.1.1' \
  --data '{
	"refresh_token": "****"
}'
```
