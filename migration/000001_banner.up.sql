create table if not exists public.tag
(
    id serial primary key,
    name text not null
);

create table if not exists public.feature
(
    id serial primary key,
    name text not null
);

create table if not exists public.banner
(
    id serial primary key,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null,
    tag_ids integer[] not null,
    feature_id integer not null,
    is_active bool not null,
    version integer not null,
    constraint fk_feature_id foreign key (feature_id) references public.feature (id),
    constraint fk_tag_ids foreign key (tag_ids) references public.tag (id)
);

create table if not exists public.banner_version
(
    id serial primary key,
    banner_id integer not null,
    version integer not null,
    content jsonb not null,
    updated_at timestamp with time zone not null,
    is_acrive bool not null,
    constraint fk_banner_id foreign key (banner_id) references public.banner (id) on delete cascade
);

insert into public.banner (is_active, time_created)
select true, NOW()
FROM generate_series(1, 100);

INSERT INTO public.tag (name)
SELECT 'test_tag_' || i
FROM generate_series(1, 100) AS s(i);

INSERT INTO public.feature (name)
SELECT 'test_feature_' || i
FROM generate_series(1, 100) AS s(i);


