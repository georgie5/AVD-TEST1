DROP TABLE IF EXISTS students;

CREATE TABLE students (
    Student_id SERIAL PRIMARY KEY,
    Student_name VARCHAR(255) NOT NULL,
    Email VARCHAR(255) NOT NULL,
    Class VARCHAR(255) NOT NULL,
    Major VARCHAR(255) NOT NULL,
    Age INT NOT NULL
);
