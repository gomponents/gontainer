package tokens

type mockAliasProvider struct {
	alias string
}

func (m mockAliasProvider) GetAlias(string) string {
	return m.alias
}
