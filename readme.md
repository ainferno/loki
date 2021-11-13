# Loki

## Create users table

```sql
CREATE TABLE users  (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255),
    token VARCHAR(255)
);
```

## Seed users

```sql
INSERT INTO users (name, email, password) VALUES
("Mary", "mary@example.com", "secret"),
("Vasya", "vasya@example.com", "secret"),
("Alex", "alex@example.com", "secret");
```
