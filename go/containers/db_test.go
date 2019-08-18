package containers_test

import (
  "log"
  "math/rand"
  "os"
  "testing"
  "time"

  "github.com/go-pg/pg"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"

  "github.com/Liquid-Labs/lc-rdb-service/go/rdb"
  "github.com/Liquid-Labs/terror/go/terror"
  . "github.com/Liquid-Labs/lc-entities-model/go/entities"
  . "github.com/Liquid-Labs/lc-users-model/go/users"
  /* pkg2test */ . "github.com/Liquid-Labs/lc-containers-model/go/containers"
)

func init() {
  terror.EchoErrorLog()
  rand.Seed(time.Now().UnixNano())
}

const runes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_./"
const aznLength = 16

func randStringBytes() string {
    b := make([]byte, aznLength)
    for i := range b {
        b[i] = runes[rand.Int63() % int64(len(runes))]
    }
    return string(b)
}

func retrieveContainer(id EID) (*Container, terror.Terror) {
  c := &Container{Entity:Entity{ID: id}}
  q := rdb.Connect().Model(c).Relation(`Members`).Where(`"container".id=?id`)
  if err := q.Select(); err != nil && err != pg.ErrNoRows {
    return nil, terror.ServerError(`Problem retrieving container.`, err)
  } else if err == pg.ErrNoRows {
    return nil, nil
  } else {
    return c, nil
  }
}

const (
  name = `John Doe`
  desc = `desc`
  legalID = `555-55-5555`
  legalIDType = `SSN`
  active = true
)

type ContainerIntegrationSuite struct {
  suite.Suite
  U      *User
  Thing1 *Entity
  Thing2 *Entity
}
func (s *ContainerIntegrationSuite) SetupSuite() {
  db := rdb.Connect()

  authID := randStringBytes()
  s.U = NewUser(`users`, `User1`, ``, authID, legalID, legalIDType, active)
  require.NoError(s.T(), s.U.CreateRaw(db))
  // log.Printf("User1: %s", s.User1.GetID())

  s.Thing1 = &Entity{ResourceName:`things`, Name:`Thing1`, OwnerID: s.U.GetID()}
  require.NoError(s.T(), CreateEntityRaw(s.Thing1, rdb.Connect()), `Unexpected error creating entity`)

  s.Thing2 = &Entity{ResourceName:`things`, Name:`Thing2`, OwnerID: s.U.GetID()}
  require.NoError(s.T(), CreateEntityRaw(s.Thing2, rdb.Connect()), `Unexpected error creating entity`)
}
func TestContainerIntegrationSuite(t *testing.T) {
  if os.Getenv(`SKIP_INTEGRATION`) == `true` {
    t.Skip()
  } else {
    suite.Run(t, new(ContainerIntegrationSuite))
  }
}

func (s *ContainerIntegrationSuite) TestContainerCreateNoMembers() {
  c := &Container{
    Entity: Entity{ResourceName:`containers`, Name:`Container1`, OwnerID: s.U.GetID()},
    Members: make([]*Entity, 0),
  }
  require.NoError(s.T(), c.CreateRaw(rdb.Connect()), `Unexpected error creating container`)
  assert.Equal(s.T(), `Container1`, c.GetName())
  assert.Equal(s.T(), ResourceName(`containers`), c.GetResourceName())
  assert.Equal(s.T(), s.U.GetID(), c.GetOwnerID())
  assert.Equal(s.T(), 0, len(c.Members))
  // default stuff
  assert.NotEqual(s.T(), EID(``), c.GetID(), `ID should have been set on insert.`)
  assert.NotEqual(s.T(), time.Time{}, c.GetCreatedAt(), `'Created at' should have been set on insert.`)
  assert.NotEqual(s.T(), time.Time{}, c.GetLastUpdated(), `'Last updated' should have been set on insert.`)
  assert.Equal(s.T(), false, c.IsPubliclyReadable())
}

func (s *ContainerIntegrationSuite) TestContainerCreateWithMembers() {
  c := &Container{
    Entity: Entity{ResourceName:`containers`, Name:`Container2`, OwnerID: s.U.GetID()},
    Members: []*Entity{s.Thing1, s.Thing2},
  }
  require.NoError(s.T(), c.CreateRaw(rdb.Connect()), `Unexpected error creating container`)

  assert.Equal(s.T(), `Container2`, c.GetName())
  assert.Equal(s.T(), ResourceName(`containers`), c.GetResourceName())
  assert.Equal(s.T(), s.U.GetID(), c.GetOwnerID())
  assert.Equal(s.T(), 2, len(c.Members))
  assert.Equal(s.T(), s.Thing1, c.Members[0])
  assert.Equal(s.T(), s.Thing2, c.Members[1])
  // default stuff
  assert.NotEqual(s.T(), EID(``), c.GetID(), `ID should have been set on insert.`)
  assert.NotEqual(s.T(), time.Time{}, c.GetCreatedAt(), `'Created at' should have been set on insert.`)
  assert.NotEqual(s.T(), time.Time{}, c.GetLastUpdated(), `'Last updated' should have been set on insert.`)
  assert.Equal(s.T(), false, c.IsPubliclyReadable())
}

func (s *ContainerIntegrationSuite) TestContainerRetrieveWithMembers() {
  c := &Container{
    Entity: Entity{ResourceName:`containers`, Name:`Container3`, OwnerID: s.U.GetID()},
    Members: []*Entity{s.Thing1, s.Thing2},
  }
  require.NoError(s.T(), c.CreateRaw(rdb.Connect()), `Unexpected error creating container`)

  cCopy, err := retrieveContainer(c.GetID())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), c, cCopy)
}

func (s *ContainerIntegrationSuite) TestContainersUpdate() {
  c := &Container{
    Entity: Entity{ResourceName:`containers`, Name:`Update Container`, OwnerID: s.U.GetID()},
    Members: []*Entity{s.Thing1},
  }
  require.NoError(s.T(), c.CreateRaw(rdb.Connect()), `Unexpected error creating container`)

  c.SetName(`New Name`)
  c.SetDescription(`bar`)
  c.SetPubliclyReadable(true)
  c.SetMembers([]*Entity{s.Thing2})
  require.NoError(s.T(), c.UpdateRaw(rdb.Connect()))
  assert.Equal(s.T(), `New Name`, c.GetName())
  assert.Equal(s.T(), `bar`, c.GetDescription())
  assert.Equal(s.T(), true, c.IsPubliclyReadable())
  assert.Equal(s.T(), 1, len(c.GetMembers()))
  assert.Equal(s.T(), s.Thing2, c.GetMembers()[0])

  cCopy, err := retrieveContainer(c.GetID())
  require.NoError(s.T(), err)
  assert.Equal(s.T(), c, cCopy)
}

func (s *ContainerIntegrationSuite) TestContainerArchive() {
  c := &Container{
    Entity: Entity{ResourceName:`containers`, Name:`Update Container`, OwnerID: s.U.GetID()},
    Members: []*Entity{s.Thing1},
  }
  require.NoError(s.T(), c.CreateRaw(rdb.Connect()), `Unexpected error creating container`)
  require.NoError(s.T(), c.ArchiveRaw(rdb.Connect()))
  log.Printf(`c: %+v`, c)

  eCopy, err := retrieveContainer(c.GetID())
  require.NoError(s.T(), err)
  assert.Nil(s.T(), eCopy)

  /* TODO: https://github.com/go-pg/pg/issues/1350
     `Relations()` not working with `Deleted()`
  archived := &Container{}
  assert.NoError(s.T(), rdb.Connect().Model(archived).Deleted().Relation(`Members`).Where(`"container".id=?`, c.GetID()).Select())
  assert.Equal(s.T(), c, archived)
  */
}
