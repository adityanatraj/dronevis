package dronevis

import (
	"fmt"
	"io"

	"github.com/go-yaml/yaml"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type (
	Condition struct {
		Branch string   `yaml:"branch"`
		Event  []string `yaml:"event"`
	}

	Step struct {
		Name      string    `yaml:"name"`
		DependsOn []string  `yaml:"depends_on"`
		When      Condition `yaml:"when"`

		// todo: maybe use "failure" for nice vis of termination
		// Failure string `yaml:"failure"`
	}

	Pipeline struct {
		Name  string `yaml:"name"`
		Steps []Step `yaml:"steps"`
	}
)

func Graph(rdr io.Reader) (string, error) {
	var p Pipeline
	if err := yaml.NewDecoder(rdr).Decode(&p); err != nil {
		return "", err
	}

	whens := map[string][]Step{}
	for _, step := range p.Steps {
		addCondition(whens, step)
	}

	for cond, steps := range whens {
		err := drawGraph(cond, steps)
		if err != nil {
			return "", err
		}
	}

	return "done", nil
}

func addCondition(m map[string][]Step, step Step) {
	cond := normalizeCondition(step.When)
	m[cond] = append(m[cond], step)
}

func normalizeCondition(c Condition) string {
	if c.Event[0] == "tag" {
		return "tag"
	}
	if c.Branch == "" {
		return fmt.Sprintf("%s-%v", "master", c.Event)
	}
	return fmt.Sprintf("%s-%v", c.Branch, c.Event)
}

func drawGraph(cond string, steps []Step) error {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			panic(err)
		}
		g.Close()
	}()

	nodes := map[string]*cgraph.Node{}
	for _, step := range steps {
		n, err := graph.CreateNode(step.Name)
		if err != nil {
			return err
		}
		nodes[step.Name] = n
	}

	for _, step := range steps {
		for _, dep := range step.DependsOn {
			if _, exists := nodes[dep]; !exists {
				return fmt.Errorf("dependency [%s] for [%s] unknown", dep, step.Name)
			}

			_, err := graph.CreateEdge(cond, nodes[dep], nodes[step.Name])
			if err != nil {
				return err
			}
		}
	}

	// render it to command line dot file:
	// var buf bytes.Buffer
	// if err := g.Render(graph, "dot", &buf); err != nil {
	// 	return err
	// }
	// fmt.Println(buf.String())

	// write to file directly
	outpath := fmt.Sprintf("%s.png", cond)
	if err := g.RenderFilename(graph, graphviz.PNG, outpath); err != nil {
		return err
	}

	return nil
}
