# Envoyer

## Development environment

Go to `envoyer_backend` directory 

```shell
cd envoyer_backend
```

Copy `.env.template` and create `.env` file and write your local config in it.

```shell
cp .env.template .env
```

Run migration.

```shell
go run main.go migrate
```

Run Envoyer backend

Run rabitmq, mysql database and run the backend server with this command.

```shell
go run main.go server
```

Go to `envoyer_frontend` directory 

```shell
cd ..
cd envoyer_frontend
```

Run Envoyer frontend

Copy `.env.template` and create `.env` file and write your local config in it.

```shell
cp .env.template .env
```

Run the frontend server with this command.

```shell
yarn dev
```
Check http://localhost:8081/ping

Check http://localhost:3000

## Docker environment

Go to `envoyer_backend` directory.

```shell
cd envoyer_backend
```

Copy `.env.template` and create `.env` file and write your local config in it.

```shell
cp .env.template .env
```

Go to `envoyer_frontend` directory 

```shell
cd ..
cd envoyer_frontend
```

Copy `.env.template` and create `.env` file and write your local config in it.

```shell
cp .env.template .env
```

Go back to main directory

```shell
cd ..
```

Run Envoyer

```bash
docker-compose up
```

Check http://localhost:8081/ping

Check http://localhost:3000
