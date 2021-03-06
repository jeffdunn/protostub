package gen

import (
	"fmt"

	"github.com/arachnys/protostub"
)

// This is more or less the same as the types in types.go, however I would like
// to preserve an amount of distinction between the two - in case of future
// diversion
type classData struct {
	name      string
	members   []protostub.Member
	functions []protostub.Function
	types     []protostub.ProtoType
	extend    bool
	comments  []string
	hasValue  bool
	hasName   bool
}

func messageToClass(m *protostub.Message) *classData {
	return &classData{
		name:     m.Typename(),
		members:  m.Members,
		types:    m.Types,
		extend:   m.IsExtend,
		comments: m.Comment,
	}
}

func enumToClass(e *protostub.Enum) *classData {
	return &classData{
		name:     e.Typename(),
		members:  e.Members,
		hasValue: true,
		hasName:  true,
	}
}

func serviceToClass(s *protostub.Service) *classData {
	return &classData{
		name:      s.Name(),
		functions: s.Functions,
		types:     s.Types,
		extend:    false,
		comments:  s.Comment,
	}
}

// generate a mypy/python class
func (g *generator) genClass(c *classData) error {
	if _, err := g.bw.WriteRune('\n'); err != nil {
		return err
	}

	if c.extend {
		fmt.Println("Extensions are not yet supported")
		return nil
	}

	if len(c.comments) > 0 {
		for _, i := range c.comments {
			g.indent()
			g.bw.WriteString(fmt.Sprintf("#%s\n", i))
		}
	}

	if err := g.indent(); err != nil {
		return err
	}

	if _, err := g.bw.WriteString(fmt.Sprintf("class %s:\n", c.name)); err != nil {
		return err
	}

	// we're working on members of this class, so indent - ensure to remove
	// the indent when done
	g.depth++

	defer func() {
		g.depth--
	}()

	for n, i := range c.members {
		for _, j := range i.Comment {
			g.indent()
			g.bw.WriteString(fmt.Sprintf("#%s\n", j))
		}

		if err := g.indent(); err != nil {
			return err
		}

		if _, err := g.bw.WriteString(fmt.Sprintf("%s: %s", i.Name(), i.Typename())); err != nil {

			return err
		}

		if n < len(c.members)-1 {
			if _, err := g.bw.WriteRune('\n'); err != nil {
				return err
			}
		}
	}

	for n, i := range c.functions {
		for _, j := range i.Comment {
			g.indent()
			g.bw.WriteString(fmt.Sprintf("#%s\n", j))
		}

		if err := g.indent(); err != nil {
			return err
		}

		if _, err := g.bw.WriteString(fmt.Sprintf("def %s: ...", i.Typename())); err != nil {
			return err
		}

		if n < len(c.functions)-1 {
			if _, err := g.bw.WriteRune('\n'); err != nil {
				return err
			}
		}
	}

	// let's make that constructor
	g.bw.WriteRune('\n')
	g.indent()
	g.bw.WriteString("def __init__(self, ")

	for n, i := range c.members {
		if n < len(c.members)-1 {
			g.bw.WriteString(fmt.Sprintf("%s: %s = None, ", i.Name(), i.Typename()))
			continue
		}

		g.bw.WriteString(fmt.Sprintf("%s: %s = None", i.Name(), i.Typename()))
	}

	g.bw.WriteString(fmt.Sprintf(") -> %s: ...\n", c.name))

	if c.hasName {
		g.indent()
		g.bw.WriteString(fmt.Sprintf("def Name(enumClass: %s) -> Any: ...\n", c.name))
	}

	if c.hasValue {
		g.indent()
		g.bw.WriteString("def Value(memberName: str) -> Any: ...\n")
	}

	for _, i := range c.types {
		// enums need to be treated differently
		if e, ok := i.(*protostub.Enum); ok {
			for _, j := range e.Members {
				g.indent()
				g.bw.WriteString(fmt.Sprintf("%s: Any\n", j.Name()))
			}
		}

		if err := g.indent(); err != nil {
			return err
		}

		if err := g.gen(i); err != nil {
			return err
		}
	}

	// then we just need to generate all the default methods that protoc adds to
	// python classes
	g.indent()
	g.bw.WriteString(fmt.Sprintf("def CopyFrom(self, other: %s) -> None: ...\n", c.name))
	g.indent()
	g.bw.WriteString("def SerializeToString(self) -> str: ...\n")
	g.indent()
	g.bw.WriteString("def ParseFromString(self, data: str) -> None: ...\n")
	g.indent()
	g.bw.WriteString("@staticmethod\n")
	g.indent()
	g.bw.WriteString("def ListFields() -> Tuple[FieldDescriptor, value]: ...\n")

	return nil
}
