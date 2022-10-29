CREATE TABLE IF NOT EXISTS public.user(
    id bigserial not null primary key,
    email varchar not null unique,
    encrypted_password varchar not null,
    active boolean not null
);

CREATE TABLE IF NOT EXISTS public.transaction(
    id bigserial not null primary key,
    amount decimal(28,3) not null
);

CREATE TABLE IF NOT EXISTS public.user_transaction(
    id bigserial not null primary key,
    user_id bigint not null references public.user(id),
    transaction_id bigint not null references public.transaction(id)
);