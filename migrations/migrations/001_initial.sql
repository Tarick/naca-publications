-- Write your migrate up statements here

{{ template "migrations/shared/trigger_set_timestamp.sql" . }}

CREATE TABLE publishers (
  uuid uuid primary key,
  name text unique not null,
  url text unique not null,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  modified_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp BEFORE UPDATE ON "publishers" FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

create table publication_types (
  type text not null UNIQUE
  );

insert into publication_types (type) values('rss');
insert into publication_types (type) values('scrapped');
insert into publication_types (type) values('api');

create table publications (
  uuid uuid PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  language_code varchar(2) NOT NULL,
  publisher_uuid uuid NOT NULL REFERENCES publishers(uuid) ON DELETE CASCADE ON UPDATE CASCADE ,
  type text NOT NULL REFERENCES publication_types(type) ON DELETE CASCADE ON UPDATE CASCADE,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  modified_at timestamptz NOT NULL DEFAULT NOW(),
  UNIQUE(name, publisher_uuid)
);

CREATE TRIGGER set_timestamp BEFORE UPDATE ON "publications" FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

---- create above / drop below ----

DROP trigger set_timestamp ON "publications";

DROP trigger set_timestamp ON "publishers";

DROP FUNCTION trigger_set_timestamp;

DROP TABLE "publications";
DROP TABLE "world_languages"
DROP TABLE "publishers";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
