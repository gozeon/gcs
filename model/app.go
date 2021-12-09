package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type App struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" binding:"required" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Api         string             `json:"api,omitempty" bson:"api"`
	Logo        string             `json:"logo,omitempty" bson:"logo"`
	ClassName   string             `json:"className,omitempty" bson:"class_name"`
	Pages       string             `json:"pages" bson:"pages"`

	// 目前使用amis的editor控件，请求为字符串
	//Pages       []Page             `json:"pages" bson:"pages"`
}

type Page struct {
	Label     string  `json:"label" binding:"required" bson:"label"`
	Url       string  `json:"url" bson:"url"`
	Icon      string  `json:"icon,omitempty" bson:"icon"`
	SchemaApi string  `json:"schemaApi,omitempty" bson:"schema_api"`
	ClassName string  `json:"className,omitempty" bson:"class_name"`
	Redirect  string  `json:"redirect,omitempty" bson:"redirect"`
	Children  []*Page `json:"children,omitempty" bson:"children"`
}
