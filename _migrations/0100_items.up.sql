CREATE TABLE items (
  uid           UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  ident         varchar(255) NOT NULL,
  bioguide_id   varchar(255) NOT NULL,
  loc           varchar(255) NOT NULL,
  last_mod      TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at    TIMESTAMP WITHOUT TIME ZONE
);
