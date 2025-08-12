package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type CreateAirtimeRequestTable_20250604_044624 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateAirtimeRequestTable_20250604_044624{}
	m.Created = "20250604_044624"

	migration.Register("CreateAirtimeRequestTable_20250604_044624", m)
}

// Run the migrations
func (m *CreateAirtimeRequestTable_20250604_044624) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE IF NOT EXISTS `requests` (" +
		"`request_id` int(11) NOT NULL AUTO_INCREMENT," +
		"`cust_id` int(11) NOT NULL," +
		"`request_type` varchar(255) NOT NULL," +
		"`request` text NOT NULL," +
		"`request_status` varchar(255) NOT NULL," +
		"`request_amount` decimal(10,2) NOT NULL," +
		"`request_date` datetime NOT NULL," +
		"`request_response` text DEFAULT NULL," +
		"`created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"`updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
		"PRIMARY KEY (`request_id`)," +
		"FOREIGN KEY (cust_id) REFERENCES customers(customer_id) ON UPDATE CASCADE ON DELETE CASCADE," +
		"UNIQUE KEY `request_id` (`request_id`)" +
		");")

}

// Reverse the migrations
func (m *CreateAirtimeRequestTable_20250604_044624) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update

}
