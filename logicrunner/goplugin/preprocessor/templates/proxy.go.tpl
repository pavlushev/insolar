package {{ .PackageName }}

import (
    {{- range $import, $i := .Imports }}
        {{$import}}
    {{- end }}
)

{{ range $typeStruct := .Types }}
    {{- $typeStruct }}
{{ end }}

// Reference to class of this contract
var ClassReference = foundation.Reference("{{ .ClassReference }}")

// Contract proxy type
type {{ .ContractType }} struct {
    Reference foundation.Reference
}

type ContractHolder struct {
	data []byte
}

func (r *ContractHolder) AsChild(objRef foundation.Reference) *{{ .ContractType }} {
    ref, err := proxyctx.Current.SaveAsChild(string(objRef), string(ClassReference), r.data)
    if err != nil {
        panic(err)
    }
    return &{{ .ContractType }}{Reference: foundation.Reference(ref)}
}

func (r *ContractHolder) AsDelegate(objRef foundation.Reference) *{{ .ContractType }} {
    ref, err := proxyctx.Current.SaveAsDelegate(string(objRef), string(ClassReference), r.data)
    if err != nil {
        panic(err)
    }
    return &{{ .ContractType }}{Reference: foundation.Reference(ref)}
}

// GetObject
func GetObject(ref foundation.Reference) (r *{{ .ContractType }}) {
    return &{{ .ContractType }}{Reference: ref}
}

{{ range $func := .ConstructorsProxies }}
func {{ $func.Name }}( {{ $func.Arguments }} ) *ContractHolder {
    {{ $func.InitArgs }}

    var argsSerialized []byte
    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    data, err := proxyctx.Current.RouteConstructorCall(string(ClassReference), "{{ $func.Name }}", argsSerialized)
    if err != nil {
		panic(err)
    }

    return &ContractHolder{data: data}
}
{{ end }}

// GetReference
func (r *{{ $.ContractType }}) GetReference() foundation.Reference {
    return r.Reference
}

// GetClass
func (r *{{ $.ContractType }}) GetClass() foundation.Reference {
    return ClassReference
}

{{ range $method := .MethodsProxies }}
func (r *{{ $.ContractType }}) {{ $method.Name }}( {{ $method.Arguments }} ) ( {{ $method.ResultsTypes }} ) {
    {{ $method.InitArgs }}
    var argsSerialized []byte

    err := proxyctx.Current.Serialize(args, &argsSerialized)
    if err != nil {
        panic(err)
    }

    res, err := proxyctx.Current.RouteCall(string(r.Reference), "{{ $method.Name }}", argsSerialized)
    if err != nil {
   		panic(err)
    }

    {{ $method.ResultZeroList }}
    err = proxyctx.Current.Deserialize(res, &resList)
    if err != nil {
        panic(err)
    }

    return {{ $method.Results }}
}
{{ end }}
