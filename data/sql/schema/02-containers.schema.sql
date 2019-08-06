CREATE TABLE containers (
  id  BIGINT,
  CONSTRAINT containers_key PRIMARY KEY ( id ),
  CONSTRAINT containers_ref_entities FOREIGN KEY ( id ) REFERENCES entities ( id )
);

CREATE TABLE container_contents (
  container BIGINT,
  item      BIGINT,
  CONSTRAINT container_contents_refs_containers FOREIGN KEY ( container ) REFERENCES containers ( id ),
  CONSTRAINT container_contents_item_ref_entities FOREIGN KEY ( item ) REFERENCES entities ( id )
);
