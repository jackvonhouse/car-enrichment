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
        "/car": {
            "get": {
                "description": "Получение автомобилей с возможностью фильтрации",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Автомобиль"
                ],
                "summary": "Получить автомобилей",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Лимит",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Смещение",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Гос. номер",
                        "name": "regNum",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Марка",
                        "name": "mark",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Модель",
                        "name": "model",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Год",
                        "name": "year",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Имя владельца",
                        "name": "ownerName",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фамилия владельца",
                        "name": "ownerSurname",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Отчество владельца",
                        "name": "ownerPatronymic",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_jackvonhouse_car-enrichment_internal_dto.Car"
                            }
                        }
                    },
                    "404": {
                        "description": "Автомобили отсутствуют",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Неизвестная ошибка",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Создание и обогощение автомобиля",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Автомобиль"
                ],
                "summary": "Создание автомобиля",
                "parameters": [
                    {
                        "description": "Массив гос. номеров",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_jackvonhouse_car-enrichment_internal_dto.CreateCar"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "result": {
                                    "type": "boolean"
                                }
                            }
                        }
                    },
                    "409": {
                        "description": "Автомобиль или владелец уже существует",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Неизвестная ошибка",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        },
        "/car/{id}": {
            "put": {
                "description": "Обновление автомобиля",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Автомобиль"
                ],
                "summary": "Обновить автомобиль",
                "parameters": [
                    {
                        "description": "Данные об автомобиле",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_jackvonhouse_car-enrichment_internal_dto.Car"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "Идентификатор автомобиля",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "result": {
                                    "type": "boolean"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Автомобиль или владелец не найдены",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "409": {
                        "description": "Автомобиль уже существует",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Неизвестная ошибка",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаление автомобиля",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Автомобиль"
                ],
                "summary": "Удалить автомобиль",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор автомобиля",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "result": {
                                    "type": "boolean"
                                }
                            }
                        }
                    },
                    "404": {
                        "description": "Автомобиль не найден",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Неизвестная ошибка",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "error": {
                                    "type": "string"
                                }
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_jackvonhouse_car-enrichment_internal_dto.Car": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "mark": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "owner": {
                    "$ref": "#/definitions/github_com_jackvonhouse_car-enrichment_internal_dto.Owner"
                },
                "regNum": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "github_com_jackvonhouse_car-enrichment_internal_dto.CreateCar": {
            "type": "object",
            "properties": {
                "reg_numbers": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "github_com_jackvonhouse_car-enrichment_internal_dto.Owner": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "patronymic": {
                    "type": "string"
                },
                "surname": {
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
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Каталог автомобилей",
	Description:      "Простейшее API для каталога автомобилей",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
