-- สร้าง Enum สำหรับ Role
CREATE TYPE user_role AS ENUM ('admin', 'user');

-- สร้าง Table Users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL, 
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (username, password, role) 
VALUES ('admin', '$2a$10$ZNs8QQLX.u7V3IPSTONcL.NaypZp81kxYzdsgiTRGr6UBwsomuY/q', 'admin');