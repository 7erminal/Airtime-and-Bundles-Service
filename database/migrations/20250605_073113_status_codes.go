package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type StatusCodes_20250605_073113 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &StatusCodes_20250605_073113{}
	m.Created = "20250605_073113"

	migration.Register("StatusCodes_20250605_073113", m)
}

// Run the migrations
func (m *StatusCodes_20250605_073113) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE status_codes(`status_id` int(11) NOT NULL AUTO_INCREMENT,`status_code` varchar(50) NOT NULL,`status_description` varchar(255) DEFAULT NULL,`active` int(11) DEFAULT 1,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,PRIMARY KEY (`status_id`))")
}

// Reverse the migrations
func (m *StatusCodes_20250605_073113) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `status_codes`")
}
