package containers

import (
  "github.com/Liquid-Labs/lc-entities-model/go/entities"
)

type Container struct {
  entities.Entity
  Contents []*entities.Entity
}

func NewContainer(
    name string,
    description string,
    ownerPubID entities.PublicID,
    publiclyReadable bool) *Container {
  return &Container{
    *entities.NewEntity(name, description, ownerPubID, publiclyReadable),
    []*entities.Entity{},
  }
}

func (c *Container) Clone() *Container {
  return &Container{*c.Entity.Clone(), c.Contents}
}

func (c *Container) CloneNew() *Container {
  return &Container{*c.Entity.CloneNew(), c.Contents}
}
