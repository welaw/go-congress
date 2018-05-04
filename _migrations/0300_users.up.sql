CREATE TABLE users (
  uid           UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  username      varchar(255) NOT NULL,
  bioguide_id   varchar(255) NOT NULL,
  lis_id        varchar(255),
  created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at    TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX users_bioguide_idx ON users(bioguide_id) WHERE deleted_at IS NULL;
