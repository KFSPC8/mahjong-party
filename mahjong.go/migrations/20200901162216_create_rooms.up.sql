create table rooms
(
    id      text primary key,
    nonce   int   not null,
    phase   int   not null,
    players jsonb not null,
    round   jsonb
);
