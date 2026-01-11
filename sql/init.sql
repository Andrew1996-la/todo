DROP TABLE IF EXISTS tasks;

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tasks (title, description, completed) VALUES
    ('Покушать', 'Покушать вкусный обед', false),
    ('Пребрать коробку шахи', 'Выполнить разбор и определить причину неисправности', true),
    ('Найти второй носок', 'Снова пропал второй носок, провести расследование', false);