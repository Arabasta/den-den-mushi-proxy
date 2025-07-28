package puppet_trusted

func FromModel(model *Model) *Record {
	return &Record{
		Certname: model.Certname,
	}
}
