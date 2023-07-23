package main

import (
	"belt/compiler"
	"belt/reporter"
)

func main() {
	{
		file := compiler.FileFromString(`let a = 100

a = 200
println a
exit 0`, "test.bl")
		err := reporter.Error(
			reporter.WhereNew(3, 3, 13, 20),
			"name `a` is immutable but assigns twice",
		)
		reporter.Report(&err, &file)
	}
	{
		file := compiler.FileFromString(`let a = 100

a =
    200`, "test.bl")
		err := reporter.Error(
			reporter.WhereNew(3, 4, 13, 23),
			"name `a` is immutable but assigns twice",
		)
		reporter.Report(&err, &file)
	}
	{
		file := compiler.FileFromString(`let a = 100

a
    =
        200`, "test.bl")
		err := reporter.Error(
			reporter.WhereNew(3, 5, 13, 31),
			"name `a` is immutable but assigns twice",
		)
		reporter.Report(&err, &file)
	}
}
