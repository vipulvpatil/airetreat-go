package utilities

type IdGeneratorMockConstant struct {
	Id string
}

func (c *IdGeneratorMockConstant) Generate() string {
	return c.Id
}
