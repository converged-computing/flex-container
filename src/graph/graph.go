package graph

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"

	"github.com/converged-computing/jsongraph-go/jsongraph/v1/graph"
	"github.com/flux-framework/fluxion-go/pkg/fluxcli"

	"fmt"
)

type ClusterGraph struct {

	// Clients needed for aws and fluxion
	cli *fluxcli.ReapiClient

	// User preferences
	MatchPolicy string

	// These get reset between topology generations
	graph   *graph.JsonGraph
	counter int32

	// Lookup for unique id objects
	seen map[string]UniqueId

	// Lookup of nodes and edges for graph (to create at once)
	nodes map[string]*graph.Node
	edges map[string]*graph.Edge
}

// Reset the topology graph to a "zero" count and no nodes seen or created
func (t *ClusterGraph) Reset() {
	// Start counting at 1, the root is 0
	t.counter = 1
	t.seen = map[string]UniqueId{}
	t.nodes = map[string]*graph.Node{}
	t.edges = map[string]*graph.Edge{}

	// prepare a graph to load targets into
	t.graph = graph.NewGraph()

}

// AddNode adds a node to the graph
func (t *ClusterGraph) AddNode(node *graph.Node) {
	t.nodes[node.Id] = node
}

// AddEdge adds a bidirectional edge to the graph
func (t *ClusterGraph) AddEdge(source string, dest string) {

	// Since we create bi-directional, always sort
	ids := []string{source, dest}
	sort.Strings(ids)
	key := strings.Join(ids, "-")

	_, ok := t.edges[key]
	if !ok {
		t.edges[key] = &graph.Edge{Source: source, Target: dest}
	}

}

// A NewClusterGraph is associated with a region and match policy
func NewClusterGraph(matchPolicy string) *ClusterGraph {

	// Set default match policy
	if matchPolicy == "" {
		matchPolicy = "first"
	}

	// Alert the user to all the chosen parameters
	// Note that "grug" == "graphml" but probably nobody knows what grug means
	// We are using JGF for now because XML is slightly evil
	fmt.Printf(" Match policy: %s\n", matchPolicy)
	fmt.Println(" Load format: JSON Graph Format (JGF)")

	t := ClusterGraph{MatchPolicy: matchPolicy}
	t.Reset()

	// instantiate fluxion
	t.cli = fluxcli.NewReapiClient()
	fmt.Printf("Created flex resource graph %s\n", t.cli)
	return &t
}

// A unique id can hold the id and return string and other derivates of it
type UniqueId struct {
	Uid  int32
	Name string
}

// String converts the int uid to a string
func (u *UniqueId) String() string {
	return fmt.Sprintf("%d", u.Uid)
}

// Get a unique id for a node (instance or network node)
// We need both int and string, so we return a struct
func (t *ClusterGraph) GetUniqueId(name string) *UniqueId {

	// Have we seen it before?
	uid, ok := t.seen[name]

	// Nope, create a node for it!
	if !ok {
		fmt.Printf("%s is not yet seen, adding with uid %d\n", name, t.counter)
		uid = UniqueId{Uid: t.counter, Name: name}
		t.seen[name] = uid
		t.counter += 1
	}
	return &uid
}

// ReadNodeJsonGraph reads in the node JGF
// We read it in just to validate, but serialize as string
func ReadNodeJsonGraph(jsonFile string) (graph.JsonGraph, string, error) {

	g := graph.JsonGraph{}

	file, err := os.ReadFile(jsonFile)
	if err != nil {
		return g, "", fmt.Errorf("error reading %s:%s", jsonFile, err)
	}

	err = json.Unmarshal([]byte(file), &g)
	if err != nil {
		return g, "", fmt.Errorf("error unmarshalling %s:%s", jsonFile, err)
	}
	return g, string(file), nil
}

// This is to say that nn-ec17* is at the top, and the instance is connected directly
// to nn-a59. This means that two instances connected to that node are close together.
// The closer two instances are in the graph, overall, the closer. That is all of
// the information that we have!
func (t *ClusterGraph) CreateGraph(jgfFile string) error {

	// Reset counter and ids
	// Note we aren't currently using this - we give the graph from
	// json directly to fluxion. It only supports JGF v1.
	t.Reset()
	return t.initFluxionContext(jgfFile)

}

// initFluxionContext, and also save the graph to file if desired.
// If a saveFile is not provided, we save to temporary file (and clean up)
// I'm not sure why fluxion requires both the bytes and the filename path, it seems redundant.
func (t *ClusterGraph) initFluxionContext(jgfFile string) error {

	conf, err := os.ReadFile(jgfFile)
	if err != nil {
		return fmt.Errorf("error reading %s:%s", jgfFile, err)
	}

	// 2. Create the context, the default format is JGF
	// 3. Remainder of defaults should work out of the box
	// Note that the options get passed as a json string to here:
	// https://github.com/flux-framework/flux-sched/blob/master/resource/reapi/bindings/c%2B%2B/reapi_cli_impl.hpp#L412
	opts := `{"matcher_policy": "%s", "load_file": "%s", "load_format": "jgf", "match_format": "jgf"}`
	p := fmt.Sprintf(opts, t.MatchPolicy, jgfFile)

	// 4. Then pass in the JGF as a string of bytes
	err = t.cli.InitContext(string(conf), p)
	if err != nil {
		return fmt.Errorf("Error creating context: %s", err)
	}
	fmt.Printf("\n✨️ Init context complete!\n")
	return nil
}

func (f *ClusterGraph) Match(specFile string) error {
	fmt.Printf("  💻️  Request: %s\n", specFile)

	spec, err := os.ReadFile(specFile)
	if err != nil {
		return errors.New("Error reading jobspec")
	}

	_, match, _, _, number, err := f.cli.MatchAllocate(false, string(spec))
	if err != nil {
		msg := f.cli.GetErrMsg()
		fmt.Println(msg)
		return err
	}
	fmt.Printf("        Match: %s\n", match)
	fmt.Printf("        Number: %d\n", number)
	return nil
}
