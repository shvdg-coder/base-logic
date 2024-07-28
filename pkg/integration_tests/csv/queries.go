package csv

// createContactsTableQuery creates the contacts table.
const createContactsTableQuery = `CREATE TABLE contacts (
		id int NOT NULL PRIMARY KEY,
		name varchar(255),
		phone varchar(255)
    );`

// getContactsQuery retrieves the contacts.
const getContactsQuery = `SELECT id, name, phone FROM contacts;`
