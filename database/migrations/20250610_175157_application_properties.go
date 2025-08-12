package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type ApplicationProperties_20250610_175157 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &ApplicationProperties_20250610_175157{}
	m.Created = "20250610_175157"

	migration.Register("ApplicationProperties_20250610_175157", m)
}

// Run the migrations
func (m *ApplicationProperties_20250610_175157) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE application_properties(`application_property_id` int(11) NOT NULL AUTO_INCREMENT,`property_code` varchar(80) NOT NULL,`property_value` varchar(255) NOT NULL,`date_created` datetime NOT NULL,`date_modified` datetime NOT NULL,`created_by` int(11) DEFAULT NULL,`modified_by` int(11) DEFAULT NULL,`active` int(11) DEFAULT NULL,PRIMARY KEY (`application_property_id`))")
}

// Reverse the migrations
func (m *ApplicationProperties_20250610_175157) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `application_properties`")
}
