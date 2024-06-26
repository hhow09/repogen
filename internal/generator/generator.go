package generator

import (
	"go/types"
	"log"

	"github.com/sunboyy/repogen/internal/codegen"
	"github.com/sunboyy/repogen/internal/mongo"
	"github.com/sunboyy/repogen/internal/spec"
)

func GenerateRepositoryImpl(modelPkg, repoPkg, destPkg *types.Package, structModelName,
	repoInterfaceName string) (string, error) {

	namedStruct, intf, err := deriveSourceTypes(modelPkg, repoPkg, structModelName,
		repoInterfaceName)
	if err != nil {
		return "", err
	}

	methodSpecs, err := constructRepositorySpec(repoPkg, namedStruct, intf)
	if err != nil {
		return "", err
	}

	codeBuilder, err := constructCodeBuilder(destPkg, namedStruct,
		repoInterfaceName, methodSpecs)
	if err != nil {
		return "", err
	}

	return codeBuilder.Build()
}

func deriveSourceTypes(modelPkg, repoPkg *types.Package, structModelName string,
	repositoryInterfaceName string) (*types.Named, *types.Interface, error) {

	structModelObj := modelPkg.Scope().Lookup(structModelName)
	if structModelObj == nil {
		return nil, nil, ErrStructNotFound
	}
	namedStruct := structModelObj.Type().(*types.Named)
	if _, ok := namedStruct.Underlying().(*types.Struct); !ok {
		return nil, nil, ErrNotNamedStruct
	}

	intfObj := repoPkg.Scope().Lookup(repositoryInterfaceName)
	if intfObj == nil {
		return nil, nil, ErrInterfaceNotFound
	}
	intf, ok := intfObj.Type().Underlying().(*types.Interface)
	if !ok {
		return nil, nil, ErrNotInterface
	}

	return namedStruct, intf, nil
}

func constructRepositorySpec(pkg *types.Package, namedStruct *types.Named,
	intf *types.Interface) ([]spec.MethodSpec, error) {

	var methodSpecs []spec.MethodSpec
	for i := 0; i < intf.NumMethods(); i++ {
		method := intf.Method(i)
		log.Println("Generating method:", method.Name())

		methodSpec, err := spec.ParseInterfaceMethod(pkg, namedStruct, method)
		if err != nil {
			return nil, err
		}
		methodSpecs = append(methodSpecs, methodSpec)
	}

	return methodSpecs, nil
}

func constructCodeBuilder(pkg *types.Package, namedStruct *types.Named,
	interfaceName string, methodSpecs []spec.MethodSpec) (*codegen.Builder, error) {

	generator := mongo.NewGenerator(pkg, namedStruct, interfaceName)
	codeBuilder := codegen.NewBuilder(
		"repogen",
		pkg.Name(),
		generator.Imports(),
	)

	constructorBuilder, err := generator.GenerateConstructor()
	if err != nil {
		return nil, err
	}

	codeBuilder.AddImplementer(constructorBuilder)
	codeBuilder.AddImplementer(generator.GenerateStruct())

	for _, method := range methodSpecs {
		methodBuilder, err := generator.GenerateMethod(method)
		if err != nil {
			return nil, err
		}
		codeBuilder.AddImplementer(methodBuilder)
	}

	return codeBuilder, nil
}
