package building

type Builder struct {
	Name string
	F    func(*B)
}

func Main(builders []Builder) {
	for _, builder := range builders {
		builder.F(&B{})
	}
}
