CREATE TABLE IF NOT EXISTS public.users(
    id bigserial not null primary key,
    email varchar not null unique,
    encrypted_password varchar not null
);

CREATE TABLE IF NOT EXISTS public.transaction(
    id bigserial not null primary key,
    amount decimal(28,3) not null,
    date_time timestamp not null
);

CREATE TABLE IF NOT EXISTS public.user_transaction(
    id bigserial not null primary key,
    user_id bigint not null references public.users(id),
    transaction_id bigint not null references public.transaction(id)
);

CREATE TABLE IF NOT EXISTS public.balance(
    id bigserial not null primary key,
    active boolean not null
);

CREATE TABLE IF NOT EXISTS public.user_balance(
    user_id bigint not null references public.users(id),
    balance_id bigint not null references public.balance(id)
);

CREATE TABLE IF NOT EXISTS public.balance_audit(
    balance_audit_id bigserial not null primary key,
    balance_id bigint not null references public.balance(id),
    balance decimal(28,3) not null,
    last_audit_time timestamp not null
);