// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Oleksandr Tolstoi",
            "email": "tolstoi.job@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/refresh-tokens": {
            "post": {
                "description": "Refresh users tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Refresh users tokens",
                "parameters": [
                    {
                        "description": "JSON schema for tokens refresh",
                        "name": "schema",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.userRefreshTokensSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User tokens was successfully refreshed",
                        "schema": {
                            "$ref": "#/definitions/entity.Tokens"
                        }
                    },
                    "400": {
                        "description": "Invalid JSON or wrong type of JSON values",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "422": {
                        "description": "Validation failed through invalid fields",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.ValidationError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/auth/sign-in": {
            "post": {
                "description": "Logins users accounts",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logins users accounts",
                "parameters": [
                    {
                        "description": "JSON schema for user login",
                        "name": "schema",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.userSignInSchema"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User was successfully logged in",
                        "schema": {
                            "$ref": "#/definitions/entity.Tokens"
                        }
                    },
                    "400": {
                        "description": "Invalid JSON or wrong type of JSON values",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "403": {
                        "description": "Your account has been suspended",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "422": {
                        "description": "Validation failed through invalid fields",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.ValidationError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/auth/sign-out": {
            "post": {
                "security": [
                    {
                        "JWT-RS256": []
                    }
                ],
                "description": "Logout users from the server",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout users from the server",
                "responses": {
                    "200": {
                        "description": "User was successfully signed out"
                    },
                    "401": {
                        "description": "Access is denied due to invalid credentials",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "403": {
                        "description": "Your account has been suspended",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/auth/sign-up": {
            "post": {
                "description": "Registers users accounts",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Registers users accounts",
                "parameters": [
                    {
                        "description": "JSON schema for user account registration",
                        "name": "schema",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.userSignUpSchema"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User account was successfully registered"
                    },
                    "400": {
                        "description": "Invalid JSON or wrong type of JSON values",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "422": {
                        "description": "Validation failed through invalid fields",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.ValidationError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/users/me": {
            "get": {
                "security": [
                    {
                        "JWT-RS256": []
                    }
                ],
                "description": "Get users personal information",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get users personal information",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.User"
                        }
                    },
                    "401": {
                        "description": "Access is denied due to invalid credentials",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "403": {
                        "description": "Your account has been suspended",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/v1.errResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/entity.CoreError"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.CoreError": {
            "description": "Basic representation of API call error",
            "type": "object",
            "properties": {
                "code": {
                    "description": "Code is CAPS_CASE constant error code you can programmatically consume to make resolution decisions from",
                    "allOf": [
                        {
                            "$ref": "#/definitions/errorcode.ErrorCode"
                        }
                    ],
                    "example": "ERROR_CODE"
                },
                "message": {
                    "description": "Message indicate a (usually) human-readable description of the error",
                    "type": "string",
                    "example": "error cause description"
                }
            }
        },
        "entity.Tokens": {
            "description": "Pair of access and refresh token which uses for auth operations",
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string",
                    "example": "header.payload.signature"
                },
                "refreshToken": {
                    "type": "string",
                    "example": "header.payload.signature"
                }
            }
        },
        "entity.User": {
            "description": "User entity information",
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string",
                    "example": "2022-12-24T21:49:33.072726+02:00"
                },
                "email": {
                    "type": "string",
                    "example": "tolstoi.job@gmail.com"
                },
                "id": {
                    "type": "string",
                    "example": "63a75a2574ef628a127ee972"
                },
                "suspendedAt": {
                    "description": "SuspendedAt is a date when user was suspended through certain reasons (optional)",
                    "type": "string",
                    "example": "2022-12-25T14:25:58.821989+02:00"
                },
                "updatedAt": {
                    "description": "UpdatedAt is a date of last user personal information modification",
                    "type": "string",
                    "example": "2022-12-24T21:58:27.072726+02:00"
                },
                "username": {
                    "type": "string",
                    "example": "kenplix"
                }
            }
        },
        "entity.ValidationError": {
            "description": "Standardized representation of a validation errors",
            "type": "object",
            "properties": {
                "code": {
                    "description": "Code is CAPS_CASE constant error code you can programmatically consume to make resolution decisions from",
                    "allOf": [
                        {
                            "$ref": "#/definitions/errorcode.ErrorCode"
                        }
                    ],
                    "example": "ERROR_CODE"
                },
                "field": {
                    "description": "Field with which validation error related",
                    "type": "string",
                    "example": "invalid field"
                },
                "message": {
                    "description": "Message indicate a (usually) human-readable description of the error",
                    "type": "string",
                    "example": "error cause description"
                }
            }
        },
        "errorcode.ErrorCode": {
            "description": "Representation of available error codes",
            "type": "string",
            "enum": [
                "PARSING_ERROR",
                "INVALID_SCHEMA",
                "INVALID_FIELD",
                "MISSING_FIELD",
                "ALREADY_EXISTS",
                "INCORRECT_CREDENTIALS",
                "UNAUTHORIZED_ACCESS",
                "CURRENT_USER_SUSPENDED",
                "INTERNAL_ERROR"
            ],
            "x-enum-varnames": [
                "ParsingError",
                "InvalidSchema",
                "InvalidField",
                "MissingField",
                "AlreadyExists",
                "IncorrectCredentials",
                "UnauthorizedAccess",
                "CurrentUserSuspended",
                "InternalError"
            ]
        },
        "v1.errResponse": {
            "description": "Standardized representation of an errors that may occur in API calls",
            "type": "object",
            "properties": {
                "errors": {
                    "description": "Array of errors that occurred performing API call\nUsually have one error inside, except validation errors where we have error for each invalid field",
                    "type": "array",
                    "items": {}
                }
            }
        },
        "v1.userRefreshTokensSchema": {
            "type": "object",
            "required": [
                "refreshToken"
            ],
            "properties": {
                "refreshToken": {
                    "type": "string",
                    "example": "header.payload.signature"
                }
            }
        },
        "v1.userSignInSchema": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "example": "kenplix or tolstoi.job@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "1wE$Rty2"
                }
            }
        },
        "v1.userSignUpSchema": {
            "type": "object",
            "required": [
                "email",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "tolstoi.job@gmail.com"
                },
                "password": {
                    "type": "string",
                    "example": "1wE$Rty2"
                },
                "username": {
                    "type": "string",
                    "example": "kenplix"
                }
            }
        }
    },
    "securityDefinitions": {
        "JWT-RS256": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "localhost:80",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "URL shortener API",
	Description:      "REST API for URL shortener app",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
