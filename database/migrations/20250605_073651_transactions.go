package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Transactions_20250605_073651 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Transactions_20250605_073651{}
	m.Created = "20250605_073651"

	migration.Register("Transactions_20250605_073651", m)
}

// Run the migrations
func (m *Transactions_20250605_073651) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE bil_transactions(`transaction_id` int(11) NOT NULL AUTO_INCREMENT,`service_id` int(11) NOT NULL,`request_id` int(11) NOT NULL,`transaction_by` int(11) NOT NULL,`amount` int(11) DEFAULT NULL,`transacting_currency` varchar(255) DEFAULT NULL,`source_channel` varchar(255) DEFAULT NULL,`source` varchar(255) NOT NULL,`destination` varchar(255) NOT NULL,`charge` float DEFAULT NULL,`external_reference_number` varchar(255) DEFAULT NULL,`status_id` int(11) NOT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) DEFAULT 1,`modified_by` int(11) DEFAULT 1,`active` int(11) DEFAULT 1,PRIMARY KEY (`transaction_id`), FOREIGN KEY (service_id) REFERENCES services(service_id) ON UPDATE CASCADE ON DELETE NO ACTION, FOREIGN KEY (request_id) REFERENCES requests(request_id) ON UPDATE CASCADE ON DELETE NO ACTION, FOREIGN KEY (status_id) REFERENCES status_codes(status_id) ON UPDATE CASCADE ON DELETE NO ACTION, FOREIGN KEY (transaction_by) REFERENCES customers(customer_id) ON UPDATE CASCADE ON DELETE NO ACTION)")
}

// Reverse the migrations
func (m *Transactions_20250605_073651) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `transactions`")
}
