package regex

const (
	BaseImport = `[A-Za-z](\/?[A-Z-a-z0-9._-])*`
	Import     = `((` + BaseImport + `)|("` + BaseImport + `"))`
	GoToken    = `[A-Za-z][A-Za-z0-9_]*`
	YamlToken  = `[A-Za-z](\.?[A-Za-z0-9_])*`
	GoFunc     = `((?P<import>` + Import + `)\.)?(?P<fn>` + GoToken + `)`

	MetaPkg           = GoToken
	MetaContainerType = GoToken
	MetaImport        = Import
	MetaImportAlias   = YamlToken
	MetaFn            = GoToken
	MetaGoFn          = GoFunc

	ParamName = YamlToken

	ServiceName   = YamlToken
	ServiceGetter = GoToken
	ServiceType   = `(?P<ptr>\*)?` + `((?P<import>` + Import + `)\.)?(?P<type>` + GoToken + `)`
	// ServiceValue supports
	// Value
	// &Value
	// my/import.Value
	// "my/import".GlobalVar.Field
	// &"my/import".GlobalVar.Field
	// MyStruct{}
	// &MyStruct{}
	// my/import.MyStruct{}
	// &my/import.MyStruct{}
	ServiceValue = `((?P<v1>(?P<ptr>\&)?((?P<import>` + Import + `)\.)?(?P<value>` + GoToken + `(\.` + GoToken + `)*` + `))` +
		`|(?P<v2>(?P<ptr2>\&)?((?P<import2>` + Import + `)\.)?(?P<struct2>` + GoToken + `)\{\}))`
	ServiceConstructor = GoFunc
	ServiceCallName    = GoToken
	ServiceFieldName   = GoToken
	ServiceTag         = YamlToken

	ArgTagged = "!tagged +(?P<tag>" + YamlToken + ")"
	ArgValue  = "!value +(" + ServiceValue + ")"
)
