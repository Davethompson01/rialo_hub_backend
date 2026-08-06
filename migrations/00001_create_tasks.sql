
CREATE TYPE status AS ENUM('Ongoing', 'Accepted', 'Completed');

CREATE Table tasks(
    task_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    reward INTEGER NOT NULL,
    status status,
    deadline TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE task_application_status AS ENUM('Ongoing', 'Accepted', 'Rejected');

CREATE TABLE task_application(
    task_application_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    task_id INTEGER REFERENCES tasks(task_id),
    employee_id INTEGER REFERENCES users(user_id),
    employer_id INTEGER REFERENCES users(user_id),
    skills TEXT,
    status task_application_status,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)