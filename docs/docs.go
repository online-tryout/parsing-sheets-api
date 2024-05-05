// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/parsing-sheets/parse": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Creates a new tryout by parsing google sheet with the provided parameters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Parser Sheets"
                ],
                "summary": "Create a new tryout by parsing google sheets",
                "parameters": [
                    {
                        "description": "Request body to create a new tryout by parsing google sheets",
                        "name": "requestBody",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.ParsingSheetsParamRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/api.ParsingSheetsParamResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "api.ModuleResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "moduleOrder": {
                    "type": "integer"
                },
                "questions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.QuestionResponse"
                    }
                },
                "title": {
                    "type": "string"
                },
                "tryoutId": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "api.OptionResponse": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "isTrue": {
                    "type": "boolean"
                },
                "optionOrder": {
                    "type": "integer"
                },
                "questionId": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "api.ParsingSheetsParamRequest": {
            "type": "object",
            "properties": {
                "endedAt": {
                    "type": "string"
                },
                "price": {
                    "type": "string"
                },
                "startedAt": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "api.ParsingSheetsParamResponse": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "endedAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "modules": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.ModuleResponse"
                    }
                },
                "price": {
                    "type": "string"
                },
                "startedAt": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "api.QuestionResponse": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "moduleId": {
                    "type": "string"
                },
                "options": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.OptionResponse"
                    }
                },
                "questionOrder": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Parsing Sheet API Documentation",
	Description:      "This is a documentation for Online Tryout Apps",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
