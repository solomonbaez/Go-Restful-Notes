DROP TABLE IF EXISTS notes;
CREATE TABLE notes (
    id       INT AUTO_INCREMENT NOT NULL,
    title    VARCHAR(128) NOT NULL,
    content  TEXT NOT NULL,
    PRIMARY KEY (`id`)
);

INSERT INTO notes
    (title, content)
VALUES
    ("Test", "This is a test note!")