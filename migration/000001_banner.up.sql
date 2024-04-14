create table if not exists public.tag
(
    id integer primary key,
    name text not null
);

create table if not exists public.feature
(
    id integer primary key,
    name text not null
);

create table if not exists public.banner
(
    id serial primary key,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    is_active bool not null,
    last_version integer not null,
    active_version integer not null
    );

create table if not exists public.banner_content
(
    id serial primary key,
    banner_id integer not null,
    version integer not null,
    content jsonb not null,
    updated_at timestamp with time zone not null,
    constraint fk_banner foreign key (banner_id) references public.banner (id) on delete cascade
    );

create table if not exists public.banner_feature_tag
(
    id serial primary key,
    banner_id integer references banner(id) on delete cascade,
    feature_id integer references feature(id),
    tag_id integer references tag(id),
    version integer not null,
    updated_at timestamp with time zone not null
);




