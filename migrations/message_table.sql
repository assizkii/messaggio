create table if not exists messages (
      id serial primary key,
      phone text,
      text text

);
create index if not exists phonex on messages (phone);
