package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Networks_20250606_125011 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Networks_20250606_125011{}
	m.Created = "20250606_125011"

	migration.Register("Networks_20250606_125011", m)
}

// Run the migrations
func (m *Networks_20250606_125011) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE networks(`network_id` int(11) NOT NULL AUTO_INCREMENT,`name` varchar(80) NOT NULL,`network_code` varchar(80) DEFAULT NULL,`network_reference_id` varchar(250) DEFAULT NULL,`description` varchar(255) DEFAULT NULL,`operator_id` int(11) NOT NULL, `date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,`active` int(11) DEFAULT 1,PRIMARY KEY (`network_id`),  FOREIGN KEY (operator_id) REFERENCES operators(operator_id) ON UPDATE CASCADE ON DELETE NO ACTION)")
}

// Reverse the migrations
func (m *Networks_20250606_125011) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `networks`")
}
