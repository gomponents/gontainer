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
	// my/import.Value
	// my/import.StructName{}.Method
	ServiceValue       = `((?P<import>` + Import + `)\.)?((?P<struct>` + GoToken + `)\{\}\.)?(?P<value>` + GoToken + `)`
	ServiceConstructor = GoFunc
	ServiceCallName    = GoToken
	ServiceFieldName   = GoToken
	ServiceTag         = YamlToken

	TaggedArg = "!tagged +(?P<tag>" + YamlToken + ")"
)
