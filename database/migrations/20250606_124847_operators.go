package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Operators_20250606_124847 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Operators_20250606_124847{}
	m.Created = "20250606_124847"

	migration.Register("Operators_20250606_124847", m)
}

// Run the migrations
func (m *Operators_20250606_124847) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE operators(`operator_id` int(11) NOT NULL AUTO_INCREMENT,`operator_name` varchar(80) NOT NULL,`description` varchar(255) DEFAULT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT NULL,`modified_by` int(11) DEFAULT 1,`active` int(11) DEFAULT 1,PRIMARY KEY (`operator_id`))")
}

// Reverse the migrations
func (m *Operators_20250606_124847) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `operators`")
}
