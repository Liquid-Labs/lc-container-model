package containers

import (
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

type Container struct {
  tableName struct{} `sql:"select:containers_join_entity"`
  Entity
  Members []*Entity `pg:"many2many:container_members,joinFK:member"`
}

func NewContainer(
    exemplar Identifiable,
    name string,
    description string,
    ownerID EID,
    publiclyReadable bool,
    members []*Entity) *Container {
  return &Container{
    struct{}{},
    *NewEntity(exemplar, name, description, ownerID, publiclyReadable),
    members,
  }
}

func (c *Container) Clone() *Container {
  return &Container{struct{}{}, *c.Entity.Clone(), c.Members}
}

func (c *Container) CloneNew() *Container {
  return &Container{struct{}{}, *c.Entity.CloneNew(), c.Members}
}

func (c *Container) GetMembers() []*Entity { return c.Members }
func (c *Container) SetMembers(m []*Entity) { c.Members = m }
func (c *Container) AddMember(e *Entity) []*Entity {
  c.Members = append(c.Members, e)
  return c.Members
}
func (c *Container) RemoveMember(t *Entity) (bool, []*Entity) {
  for i, e := range c.Members {
    if e.GetID() == t.GetID() {
      return c.RemoveMemberAt(i)
    }
  }
  return false, c.Members
}
func (c *Container) RemoveMemberAt(i int) (bool, []*Entity) {
  if (i < 0 || i > len(c.Members)) {
    return false, c.Members
  } else {
    c.Members = append(c.Members[:i], c.Members[i+1:]...)
    return true, c.Members
  }
}

type ContainerMembers struct {
  ContainerID EID
  Member      EID
}
