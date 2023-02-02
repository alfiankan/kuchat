CREATE EXTENSION pgcrypto;

CREATE TABLE vmq_auth_acl
 (
   mountpoint character varying(10) NOT NULL,
   client_id character varying(128) NOT NULL,
   username character varying(128) NOT NULL,
   password character varying(128),
   publish_acl json,
   subscribe_acl json,
   CONSTRAINT vmq_auth_acl_primary_key PRIMARY KEY (mountpoint, client_id, username)
 );

WITH x AS (
    SELECT
        ''::text AS mountpoint,
           'test-client'::text AS client_id,
           'test-user'::text AS username,
           '123'::text AS password,
           gen_salt('bf')::text AS salt,
           '[{"pattern": "a/b/c"}, {"pattern": "c/b/#"}]'::json AS publish_acl,
           '[{"pattern": "a/b/c"}, {"pattern": "c/b/#"}]'::json AS subscribe_acl
    )
INSERT INTO vmq_auth_acl (mountpoint, client_id, username, password, publish_acl, subscribe_acl)
    SELECT
        x.mountpoint,
        x.client_id,
        x.username,
        crypt(x.password, x.salt),
        publish_acl,
        subscribe_acl
    FROM x;


---- USER

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
 (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    email character varying(128) NOT NULL,
    password character varying(128),
    CONSTRAINT users_primary_key PRIMARY KEY (id)
 );
