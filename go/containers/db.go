package containers

import (
  "github.com/go-pg/pg/orm"

  . "github.com/Liquid-Labs/terror/go/terror"
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

var ContainerFields = append(EntityFields, `members`)

// ModelContainer provides a(n initially empty) Entity receiver and base query.
func ModelContainer(db orm.DB) (*Container, *orm.Query) {
  c := &Container{}
  q := db.Model(c)

  return c, q
}

func (c *Container) createMembers(db orm.DB) Terror {
  if 0 < len(c.GetMembers()) {
    newMembers := make([]*ContainerMembers, 0)
    for _, member := range c.GetMembers() {
      newMembers = append(newMembers,
        &ContainerMembers{c.GetID(), member.GetID()})
    }

    if _, err := db.Model(&newMembers).Insert(); err != nil {
      return ServerError(`There was a problem creating container members.`, err)
    }
  }
  return nil
}

// Create creates (or inserts) a new Container record into the DB. As Containers are logically abstract, one would typically only call this as part of another items create sequence.
func (c *Container) CreateRaw(db orm.DB) Terror {
  if err := CreateEntityRaw(c, db); err != nil {
    return err
  } else {
    qs := db.Model(c).ExcludeColumn(EntityFields...)
    if _, err := qs.Insert(); err != nil {
      return ServerError(`There was a problem creating the container record.`, err)
    } else {
      return c.createMembers(db)
    }
  }
}

var updateExcludes = make([]string, len(EntityFields))
func init() {
  copy(updateExcludes, EntityFields)
  updateExcludes = append(updateExcludes, "id")
}
// Update updates a container record in the DB. As Containers are logically abstract, one would typically only call this as part of another items update sequence.
func (c *Container) UpdateRaw(db orm.DB) Terror {
  if err := (&c.Entity).UpdateRaw(db); err != nil {
    return err
  } else {
    // We don't need to update the 'container' record itself; it's just an ID. Besides, ff all columns are excluded, go-pg (v8.0.5) ignores exclusions.
    if _, err := db.Exec(`DELETE FROM container_members WHERE container_id=?`, c.GetID()); err != nil {
      return ServerError(`Problem updataing container members.`, err)
    }
    return c.createMembers(db)
  }
}

// Archive updates a Container record in the DB. As Containers are logically abstract, one would typically only call this as part of another items archive sequence.
func (c *Container) ArchiveRaw(db orm.DB) Terror {
  return (&c.Entity).ArchiveRaw(db)
}
