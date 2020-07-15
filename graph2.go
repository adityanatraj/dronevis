package dronevis

import (
	"fmt"
	"io"
	"sort"
	"strings"

	yaml2 "github.com/drone/drone-yaml/yaml"
	"github.com/emicklei/dot"
	"github.com/go-yaml/yaml"
)

type (
	node struct {
		stepInfo *yaml2.Container
	}

)

func Graph2(rdr io.Reader) ([]string, error) {

	var p yaml2.Pipeline
	if err := yaml.NewDecoder(rdr).Decode(&p); err != nil {
		return nil, err
	}

	nodes := map[string]node{}
	whens := map[string]struct{}{}
	for _, step := range p.Steps {
		// fmt.Println(step.Name)
		// fmt.Println(step.DependsOn)

		nodes[step.Name] = node{
			stepInfo: step,
		}
		whens[flattenConditions(step.When)] = struct{}{}
	}

	graphs := make([]string, 0, len(whens))
	for when := range whens {
		graph := dot.NewGraph(dot.Directed)
		for name, node := range nodes {
			for _, dep := range node.stepInfo.DependsOn {
				if _, exists := nodes[dep]; !exists {
					fmt.Printf("unknown dep [%s] for step [%s]\n", dep, name)
				} else {
					if cond := flattenConditions(node.stepInfo.When); cond != when {
						continue
					}
					graph.Edge(graph.Node(dep), graph.Node(name))
					// e.Label(flattenConditions(node.stepInfo.When))
				}
			}
		}
		graphs = append(graphs, graph.String())
	}

	return graphs, nil
}

func flattenConditions(conds yaml2.Conditions) string {
	// every condition has N includes and N excludes
	// regardless of what type it is
	// even matrix builds are really just a set of includes
	// and the runtime executes a special build for it
	// (it's own runtime pipeline if you will)
	//

	// each one of these can enumerate a unique pipeline
	// amongst it's excludes and includes
	// the # of rendered pipelines is equal to
	// (# keys in matrix) + PerKey (# includes + # excludes)

	var pieces []string
	// Action   Condition         `json:"action,omitempty"`
	if a := flattenCondition("action", conds.Action); a != "" {
		pieces = append(pieces, a)
	}

	// Cron     Condition         `json:"cron,omitempty"`
	if a := flattenCondition("cron", conds.Cron); a != "" {
		pieces = append(pieces, a)
	}

	// Ref      Condition         `json:"ref,omitempty"`
	if a := flattenCondition("ref", conds.Ref); a != "" {
		pieces = append(pieces, a)
	}

	// Repo     Condition         `json:"repo,omitempty"`
	if a := flattenCondition("repo", conds.Repo); a != "" {
		pieces = append(pieces, a)
	}

	// Instance Condition         `json:"instance,omitempty"`
	if a := flattenCondition("inst", conds.Instance); a != "" {
		pieces = append(pieces, a)
	}

	// Target   Condition         `json:"target,omitempty"`
	if a := flattenCondition("target", conds.Target); a != "" {
		pieces = append(pieces, a)
	}

	// Event    Condition         `json:"event,omitempty"`
	if a := flattenCondition("event", conds.Event); a != "" {
		pieces = append(pieces, a)
	}

	// Branch   Condition         `json:"branch,omitempty"`
	if a := flattenCondition("branch", conds.Branch); a != "" {
		pieces = append(pieces, a)
	}

	// // Status   Condition         `json:"status,omitempty"`
	// if a := flattenCondition("status", conds.Status); a != "" {
	// 	pieces = append(pieces, a)
	// }

	// Paths    Condition         `json:"paths,omitempty"`
	if a := flattenCondition("act", conds.Paths); a != "" {
		pieces = append(pieces, a)
	}

	return strings.Join(pieces, ",")
}

func flattenCondition(name string, cond yaml2.Condition) string {
	if len(cond.Include) == 0 && len(cond.Exclude) == 0 {
		return ""
	}

	bldr := strings.Builder{}
	bldr.WriteString(name)
	bldr.WriteString(":")
	if len(cond.Include) > 0 {
		sort.Strings(cond.Include)
		bldr.WriteString(fmt.Sprintf("+%v", cond.Include))
	}
	if len(cond.Exclude) > 0 {
		sort.Strings(cond.Exclude)
		bldr.WriteString(fmt.Sprintf("-%v", cond.Exclude))
	}
	return bldr.String()
}
