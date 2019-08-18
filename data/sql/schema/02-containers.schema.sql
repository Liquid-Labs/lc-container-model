CREATE TABLE containers (
  id UUID,

  CONSTRAINT containers_key PRIMARY KEY ( id ),
  CONSTRAINT containers_ref_entities FOREIGN KEY ( id ) REFERENCES entities ( id )
);

CREATE TABLE container_members (
  container_id UUID,
  member       UUID,

  CONSTRAINT container_members_primary_key PRIMARY KEY (container_id, member),
  CONSTRAINT container_members_refs_containers FOREIGN KEY ( container_id ) REFERENCES containers ( id ),
  CONSTRAINT container_members_member_ref_entities FOREIGN KEY ( member ) REFERENCES entities ( id )
);

CREATE VIEW containers_join_entity AS
  SELECT e.*
    FROM containers c JOIN entities e ON c.id=e.id;
