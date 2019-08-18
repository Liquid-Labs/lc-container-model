package containers

import (
  "github.com/go-pg/pg/orm"

  . "github.com/Liquid-Labs/terror/go/terror"
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

// ModelContainer provides a(n initially empty) Entity receiver and base query.
func ModelContainer(db orm.DB) (*Container, *orm.Query) {
  c := &Container{}
  q := db.Model(c)

  return c, q
}

func (c *Container) createMembers(db orm.DB) Terror {
  for _, member := range c.GetMembers() {
    member := &ContainerMembers{c.GetID(), member.GetID()}
    if _, err := db.Model(member).Insert(); err != nil {
      return ServerError(`There was a problem creating container members.`, err)
    }
  }

  return nil
}

// Create creates (or inserts) a new Container record into the DB. As Containers are logically abstract, one would typically only call this as part of another items create sequence.
func (c *Container) Create(db orm.DB) Terror {
  if err := (&c.Entity).Create(db); err != nil {
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
func (c *Container) Update(db orm.DB) Terror {
  if err := (&c.Entity).Update(db); err != nil {
    return err
  } else { /* If all columns are excluded, go-pg ignores exclusions.
    qu := db.Model(c).
      ExcludeColumn(updateExcludes...).
      Where(`"container".id=?id`)
    qu.GetModel().Table().SoftDeleteField = nil
    if _, err := qu.Update(); err != nil {
      return ServerError(`There was a problem updating the container record.`, err)
    } else {*/
      db.Exec(`DELETE FROM container_members WHERE container_id=?`, c.GetID())
      return c.createMembers(db)
    // }
  }
}

// Archive updates a Container record in the DB. As Containers are logically abstract, one would typically only call this as part of another items archive sequence.
func (c *Container) Archive(db orm.DB) Terror {
  return (&c.Entity).Archive(db)
}
