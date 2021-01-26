package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/johejo/dbtypegen"
)

var (
	pkg        string
	jsonType   string
	uuidType   string
	typeSuffix string
	typePrefix string
	tag        string

	schemaFile string
	out        string
)

func init() {
	flag.StringVar(&pkg, "package", "dbtype", "package")
	flag.StringVar(&jsonType, "json-type", "json.RawMessage", "json type")
	flag.StringVar(&uuidType, "uuid-type", "string", "uuid type")
	flag.StringVar(&typeSuffix, "type-suffix", "", "type prefix")
	flag.StringVar(&typePrefix, "type-prefix", "", "type suffix")
	flag.StringVar(&tag, "tag", "db", "tag")

	flag.StringVar(&schemaFile, "schema", "", "schema file")
	flag.StringVar(&out, "out", "", "output file")
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	var r io.Reader
	if schemaFile == "" {
		r = os.Stdin
	} else {
		f, err := os.Open(schemaFile)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	opts := []dbtypegen.Option{
		dbtypegen.WithPackage(pkg),
		dbtypegen.WithJSONType(jsonType),
		dbtypegen.WithUUIDType(uuidType),
		dbtypegen.WithTypePrefix(typePrefix),
		dbtypegen.WithTypeSuffix(typeSuffix),
		dbtypegen.WithTag(tag),
	}

	data, err := dbtypegen.Generate(ctx, r, opts...)
	if err != nil {
		return err
	}

	if out == "" || out == "-" {
		if _, err := os.Stdout.Write(data); err != nil {
			return err
		}
		return nil
	}
	return ioutil.WriteFile(out, data, 0o644)
}
