package api


func (c *initAPI) initConfig(){
}

func createAPI() (*initAPI, error) {
	c := initAPI{}
	c.initConfig()

	return &c, nil
}