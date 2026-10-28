INSERT INTO categories (name) VALUES
    ('Fantasy'),
    ('Science Fiction'),
    ('History')
ON CONFLICT (name) DO NOTHING;

INSERT INTO authors (full_name, bio) VALUES
    ('George Orwell', 'Author of dystopian literature.'),
    ('Isaac Asimov', 'Science fiction and popular science writer.'),
    ('Mary Shelley', 'Author of Frankenstein.'),
    ('J.R.R. Tolkien', 'Fantasy and mythology author.'),
    ('Yuval Noah Harari', 'Historian and author.')
ON CONFLICT DO NOTHING;

INSERT INTO books (title, isbn, category_id, published_year, description)
SELECT '1984', '9780451524935', c.id, 1949, 'Classic dystopian novel.'
FROM categories c WHERE c.name = 'Science Fiction'
ON CONFLICT (isbn) DO NOTHING;

INSERT INTO books (title, isbn, category_id, published_year, description)
SELECT 'Foundation', '9780553293357', c.id, 1951, 'Epic science fiction novel.'
FROM categories c WHERE c.name = 'Science Fiction'
ON CONFLICT (isbn) DO NOTHING;

INSERT INTO books (title, isbn, category_id, published_year, description)
SELECT 'Frankenstein', '9780486282114', c.id, 1818, 'Classic Gothic novel.'
FROM categories c WHERE c.name = 'Fantasy'
ON CONFLICT (isbn) DO NOTHING;

INSERT INTO books (title, isbn, category_id, published_year, description)
SELECT 'The Hobbit', '9780547928227', c.id, 1937, 'Fantasy adventure.'
FROM categories c WHERE c.name = 'Fantasy'
ON CONFLICT (isbn) DO NOTHING;

INSERT INTO books (title, isbn, category_id, published_year, description)
SELECT 'Sapiens', '9780062316097', c.id, 2014, 'Big history narrative.'
FROM categories c WHERE c.name = 'History'
ON CONFLICT (isbn) DO NOTHING;

INSERT INTO book_authors (book_id, author_id)
SELECT b.id, a.id
FROM books b
JOIN authors a ON a.full_name = 'George Orwell'
WHERE b.isbn = '9780451524935'
ON CONFLICT DO NOTHING;

INSERT INTO book_authors (book_id, author_id)
SELECT b.id, a.id
FROM books b
JOIN authors a ON a.full_name = 'Isaac Asimov'
WHERE b.isbn = '9780553293357'
ON CONFLICT DO NOTHING;

INSERT INTO book_authors (book_id, author_id)
SELECT b.id, a.id
FROM books b
JOIN authors a ON a.full_name = 'Mary Shelley'
WHERE b.isbn = '9780486282114'
ON CONFLICT DO NOTHING;

INSERT INTO book_authors (book_id, author_id)
SELECT b.id, a.id
FROM books b
JOIN authors a ON a.full_name = 'J.R.R. Tolkien'
WHERE b.isbn = '9780547928227'
ON CONFLICT DO NOTHING;

INSERT INTO book_authors (book_id, author_id)
SELECT b.id, a.id
FROM books b
JOIN authors a ON a.full_name = 'Yuval Noah Harari'
WHERE b.isbn = '9780062316097'
ON CONFLICT DO NOTHING;
