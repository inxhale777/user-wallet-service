create table transactions
(
    id          serial primary key,
    user_id     int  not null,
    status      text not null,
    amount      int  not null,
    description text not null
);