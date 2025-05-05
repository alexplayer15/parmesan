package test_data

import oas_struct "github.com/alexplayer15/parmesan/data" // adjust the import path

func BaseOAS() oas_struct.OAS {
	return oas_struct.OAS{
		OpenAPI: "3.0.0",
		Info: oas_struct.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Servers: []oas_struct.Server{
			{URL: "http://example.com"},
		},
		Paths: map[string]map[string]oas_struct.Method{
			"/users": {
				"post": {
					Summary: "Create a user",
					RequestBody: oas_struct.RequestBody{
						Content: map[string]oas_struct.Content{
							"application/json": {
								Schema: oas_struct.Schema{
									Type:       "object",
									Properties: map[string]oas_struct.Property{},
								},
							},
						},
					},
					Parameters: []oas_struct.Parameter{},
				},
			},
		},
		Components: oas_struct.Components{
			Schemas: map[string]oas_struct.Schema{},
		},
	}
}
