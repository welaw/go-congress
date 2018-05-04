CREATE TABLE votes (
  uid           UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
  username      varchar(255) NOT NULL,
  vote          varchar(255) NOT NULL,
  upstream      varchar(255) NOT NULL,
  ident         varchar(255) NOT NULL,
  branch        varchar(255) NOT NULL,
  created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at    TIMESTAMP WITHOUT TIME ZONE
);
