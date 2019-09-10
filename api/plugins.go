package api

type PluginParams struct {
	Name    string `url:"plugin_name,omitempty"`
	Enable    bool `url:"enable,omitempty"`
}

func (api *API) EnablePlugin(name string) error {
	params := &PluginParams{Name: name}
	_, err := api.sling.Post("/api/plugins").BodyForm(params).ReceiveSuccess(nil)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) ReadPlugins() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	_, err := api.sling.Get("/api/plugins").ReceiveSuccess(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (api *API) UpdatePlugin(name string, enable bool) error {
	params := &PluginParams{Name: name, Enable: enable}
	_, err := api.sling.Put("api/plugins").BodyForm(params).ReceiveSuccess(nil)
	if err != nil {
		return err
	}
	return nil
}

// alarm_id, type
func (api *API) DisablePlugin(name string) error {
	params := &PluginParams{Name: name}
	_, err := api.sling.Delete("/api/alarms/").BodyForm(params).ReceiveSuccess(nil)
	if err != nil {
		return err
	}
	return nil
}
