# Matchmaker Go

The project will contain a master server, a game server and a game client for
creating rooms to play different games that follow the structure presented
here.

### Quickstart

First you will need to create a `.env.local` file from the `.env.example` file.
This will ensure that all the right env variables will be set for the database
and game server image that will be used.

Start up the postgres container to be able to store rooms for your games.

```console
docker compose up
```

Run the migrations for the database

```console
go run ./cmd/bun/main.go init
go run ./cmd/bun/main.go migrate
```

Start the matchmaker

```console
go run ./cmd/matchmaker/main.go
```

Now you should be able to use the API to create rooms, query the list of public
rooms, or check the details for private rooms. Then you will be able to connect
via the game client to a specific room. For example, let's use the echo game I
provide in the `examples` folder.

First you will need to build the docker image. Also you have to make sure that
the name of the image matches the name of the image from the env file.

```console
docker build -t game-echo .
```

Now if you start the matchmaker server you will be able to create rooms using
the following request (which can be integrated in the game client for example)

```console
curl -X POST localhost:8080/api/v1/rooms -H "Content-Type: application/json" -d '{"name":"name","private":true,"maxPlayers":2}'
```

(The fields from the json body are just an example, you can use any name,
private true or false and any amount of players)

To check the public rooms

```console
curl -X GET localhost:8080/api/v1/rooms
```

And to check a private room using the CODE of a room

```console
curl -X GET localhost:8080/api/v1/rooms/[CODE]
```

Highly encourage to check out the `./examples/echo.go` to see how a game server
would be used.
