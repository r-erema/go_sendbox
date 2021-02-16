package exmaple1

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"
	"go_sendbox/config"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	err         error
	driverNeo4j neo4j.Driver
	resultNeo4j neo4j.Result
	pgDB        *sql.DB
)

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	shutdown()

	os.Exit(exitVal)
}

func setup() {
	var (
		session    neo4j.Session
		dumpBuffer []byte
	)
	if driverNeo4j, err = neo4j.NewDriver(config.Neo4jDSN(), neo4j.BasicAuth("neo4j", "123", ""), func(c *neo4j.Config) {
		c.Encrypted = false
	}); err != nil {
		log.Fatal(err)
	}

	if session, err = driverNeo4j.Session(neo4j.AccessModeWrite); err != nil {
		log.Fatal(err)
	}
	if dumpBuffer, err = ioutil.ReadFile("./neo4j.dump"); err != nil {
		log.Fatal(err)
	}
	if resultNeo4j, err = session.Run(string(dumpBuffer), nil); err != nil {
		log.Fatal(err)
	}

	if pgDB, err = sql.Open("postgres", config.PostgresDSN()); err != nil {
		log.Fatal(err)
	}
	if dumpBuffer, err = ioutil.ReadFile("./postgres_dump.sql"); err != nil {
		log.Fatal(err)
	}
	if _, err = pgDB.Query(string(dumpBuffer)); err != nil {
		log.Fatal(err)
	}
}

func shutdown() {
	var session neo4j.Session
	if session, err = driverNeo4j.Session(neo4j.AccessModeWrite); err != nil {
		log.Fatal(err)
	}
	if resultNeo4j, err = session.Run("MATCH (n) DETACH DELETE n", nil); err != nil {
		log.Fatal(err)
	}
	if _, err = pgDB.Query("DROP TABLE public.edges; DROP TABLE public.vertices;"); err != nil {
		log.Fatal(err)
	}
}

func TestNeo4j(t *testing.T) {
	var session neo4j.Session
	query := `MATCH
 (person) -[:BORN_IN]-> () -[:WITHIN*0..]-> (us:Location {name:'United States'}),
 (person) -[:LIVES_IN]-> () -[:WITHIN*0..]-> (eu:Location {name:'Europe'})
RETURN person.name`

	if session, err = driverNeo4j.Session(neo4j.AccessModeRead); err != nil {
		log.Fatal(err)
	}
	if resultNeo4j, err = session.Run(query, nil); err != nil {
		t.Fatal(err)
	}

	resultNeo4j.Next()
	assert.Equal(t, "Lucy", resultNeo4j.Record().GetByIndex(0))
}

func TestPostgres(t *testing.T) {
	var rows *sql.Rows
	query := `
	WITH RECURSIVE
	 -- in_usa is the set of vertex IDs of all locations WITHIN the United States
	 in_usa(vertex_id) AS (
		 SELECT vertex_id FROM vertices WHERE properties->>'name' = 'United States'
		 UNION
		 SELECT edges.tail_vertex FROM edges
		 JOIN in_usa ON edges.head_vertex = in_usa.vertex_id
		 WHERE edges.label = 'WITHIN'
	 ),
	 -- in_europe is the set of vertex IDs of all locations WITHIN Europe
	 in_europe(vertex_id) AS (
		 SELECT vertex_id FROM vertices WHERE properties->>'name' = 'Europe'
		 UNION
		 SELECT edges.tail_vertex FROM edges
		 JOIN in_europe ON edges.head_vertex = in_europe.vertex_id
		 WHERE edges.label = 'WITHIN'
	 ),
	 -- BORN_IN_usa is the set of vertex IDs of all people born in the US
	 BORN_IN_usa(vertex_id) AS (
		 SELECT edges.tail_vertex FROM edges
		 JOIN in_usa ON edges.head_vertex = in_usa.vertex_id
		 WHERE edges.label = 'BORN_IN'
	 ),
	 LIVES_IN_europe(vertex_id) AS (
		 SELECT edges.tail_vertex FROM edges
		 JOIN in_europe ON edges.head_vertex = in_europe.vertex_id
		 WHERE edges.label = 'LIVES_IN'
	 )
	SELECT vertices.properties->>'name'
	FROM vertices
	-- join to find those people who were both born in the US *and* live in Europe
	JOIN BORN_IN_usa ON vertices.vertex_id = BORN_IN_usa.vertex_id
	JOIN LIVES_IN_europe ON vertices.vertex_id = LIVES_IN_europe.vertex_id;
`

	if rows, err = pgDB.Query(query); err != nil {
		t.Fatal(err)
	}

	var name string
	rows.Next()
	if err = rows.Scan(&name); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Lucy", name)
}
