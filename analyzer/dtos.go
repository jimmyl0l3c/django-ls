package analyzer

type ModelField struct {
	Name     string `json:"name"`
	TypeName string `json:"type_name"`
}

type DjangoModel struct {
	App    string       `json:"app"`
	Name   string       `json:"name"`
	Fields []ModelField `json:"fields"`

	lookups []FieldLookup
}

func (dm *DjangoModel) syncLookups() {
	dm.lookups = nil
	for _, f := range dm.Fields {
		dm.lookups = append(dm.lookups, FieldLookup{Name: f.Name, TypeName: f.TypeName})

		// TODO: extend the lookups
	}
}

func (dm *DjangoModel) GetLookups() []FieldLookup {
	return dm.lookups
}
