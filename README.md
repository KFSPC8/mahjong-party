# Mahjong Party
<img width="2108" height="1263" alt="image" src="https://github.com/user-attachments/assets/19bbc422-2d97-45ad-8317-2009228aeb6a" />

## To build
1. Generate 2 keys using:
```
$ openssl rand -base64 32
```
2. Create `.env` file in root directory and fill in the following environment variables with the generated keys:
```
PARLOUR_SESSION_AUTH_KEY="YOUR AUTH KEY"
PARLOUR_SESSION_ENC_KEY="YOUR ENC KEY"
```
3. Start the database with:
```
docker compose up --build
``` 
4. Migrate the database with the following commands:
```
docker exec -it mahjong-postgres psql -U mahjong -d mahjong

mahjong=# create table rooms
(
    id      text primary key,
    nonce   int   not null,
    phase   int   not null,
    players jsonb not null,
    round   jsonb
);

mahjong=# alter table rooms
    add column results jsonb;

mahjong=# \q
```
5. Finally, you should be able to start the application using:
```
docker compose up
```

## Folder structure
```
project-root/
├─ docker-compose.yml
├─ mahjong-party/      # Frontend code (React)
├─ mahjong.go/         # Backend code (Go)
├─ .env                # To store sensitive keys
```
