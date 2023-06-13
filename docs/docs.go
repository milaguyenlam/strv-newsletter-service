// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Nguyen Thanh Lam",
            "url": "https://github.com/milaguyenlam",
            "email": "milaguyenlam@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/subscription/create": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Create a new subscription with the given name and description",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new subscription",
                "operationId": "create-subscription",
                "parameters": [
                    {
                        "description": "Subscription input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CreateSubscriptionInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Subscription ID",
                        "schema": {
                            "$ref": "#/definitions/model.MessageResponse"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "$ref": "#/definitions/model.MessageResponse"
                        }
                    }
                }
            }
        },
        "/subscription/{subscriptionID}/send": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Send a newsletter to all subscribers",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Send a newsletter",
                "operationId": "send-newsletter",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Subscription ID",
                        "name": "subscriptionID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Email details",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Email"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Message",
                        "schema": {
                            "$ref": "#/definitions/model.MessageResponse"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "$ref": "#/definitions/model.MessageResponse"
                        }
                    }
                }
            }
        },
        "/subscription/{subscriptionID}/subscribe": {
            "get": {
                "description": "Subscribe to a newsletter with a given subscription ID",
                "produces": [
                    "application/json"
                ],
                "summary": "Subscribe to a newsletter",
                "operationId": "subscribe-newsletter",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Subscription ID",
                        "name": "subscriptionID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Email to subscribe",
                        "name": "email",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Message",
                        "schema": {
                            "$ref": "#/definitions/model.MessageResponse"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "$ref": "#/definitions/model.MessageResponse"
                        }
                    }
                }
            }
        },
        "/subscription/{subscriptionID}/unsubscribe": {
            "get": {
                "description": "Unsubscribe from a newsletter with a given subscription ID",
                "produces": [
                    "application/json"
                ],
                "summary": "Unsubscribe from a newsletter",
                "operationId": "unsubscribe-newsletter",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Subscription ID",
                        "name": "subscriptionID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Email to unsubscribe",
                        "name": "email",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Message",
                        "schema": {
                            "$ref": "#/definitions/model.MessageResponse"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "$ref": "#/definitions/model.MessageResponse"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "description": "Logs in a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "User Login",
                "operationId": "login",
                "parameters": [
                    {
                        "description": "Login credentials",
                        "name": "authInput",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.AuthenticationInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Token",
                        "schema": {
                            "$ref": "#/definitions/model.AuthenticationResponse"
                        }
                    },
                    "401": {
                        "description": "Message",
                        "schema": {
                            "$ref": "#/definitions/model.MessageResponse"
                        }
                    }
                }
            }
        },
        "/user/register": {
            "post": {
                "description": "Registers a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "User Registration",
                "operationId": "register",
                "parameters": [
                    {
                        "description": "Registration details",
                        "name": "authInput",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.AuthenticationInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Token",
                        "schema": {
                            "$ref": "#/definitions/model.AuthenticationResponse"
                        }
                    },
                    "401": {
                        "description": "Message",
                        "schema": {
                            "$ref": "#/definitions/model.MessageResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.AuthenticationInput": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "model.AuthenticationResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "model.CreateSubscriptionInput": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "subscriptionName": {
                    "type": "string"
                }
            }
        },
        "model.Email": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "subject": {
                    "type": "string"
                }
            }
        },
        "model.MessageResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "STRV Newsletter Subscription API",
	Description:      "This is a newsletter subscription API service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
