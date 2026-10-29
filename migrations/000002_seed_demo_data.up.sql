INSERT INTO categories (name) VALUES
    ('Fantasy'),
    ('Science Fiction'),
    ('History'),
    ('Drama')
ON CONFLICT (name) DO NOTHING;

INSERT INTO authors (full_name, bio) VALUES
    ('George Orwell', 'Author of dystopian literature.'),
    ('Isaac Asimov', 'Science fiction and popular science writer.'),
    ('Mary Shelley', 'Author of Frankenstein.'),
    ('Percy Bysshe Shelley', 'Romantic poet and writer.'),
    ('J.R.R. Tolkien', 'Fantasy and mythology author.'),
    ('Yuval Noah Harari', 'Historian and author.'),
    ('Jane Austen', 'English novelist.'),
    ('Frank Herbert', 'Science fiction novelist.')
ON CONFLICT DO NOTHING;

INSERT INTO books (title, isbn, category_id, published_year, description)
SELECT '1984', '9780451524935', c.id, 1949, 'Classic dystopian novel.'
FROM categories c WHERE c.name = 'Science Fiction'
ON CONFLICT (isbn) DO NOTHING;

INSERT INTO books (title, isbn, category_id, published_year, description)
SELECT 'Animal Farm', '9780451526342', c.id, 1945, 'Political allegory about power.'
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

INSERT INTO books (title, isbn, category_id, published_year, description)
SELECT 'Pride and Prejudice', '9780141439518', c.id, 1813, 'Classic romantic novel.'
FROM categories c WHERE c.name = 'Drama'
ON CONFLICT (isbn) DO NOTHING;

INSERT INTO books (title, isbn, category_id, published_year, description)
SELECT 'Dune', '9780441172719', c.id, 1965, 'Epic desert planet saga.'
FROM categories c WHERE c.name = 'Science Fiction'
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
JOIN authors a ON a.full_name = 'George Orwell'
WHERE b.isbn = '9780451526342'
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
JOIN authors a ON a.full_name = 'Percy Bysshe Shelley'
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

INSERT INTO book_authors (book_id, author_id)
SELECT b.id, a.id
FROM books b
JOIN authors a ON a.full_name = 'Jane Austen'
WHERE b.isbn = '9780141439518'
ON CONFLICT DO NOTHING;

INSERT INTO book_authors (book_id, author_id)
SELECT b.id, a.id
FROM books b
JOIN authors a ON a.full_name = 'Frank Herbert'
WHERE b.isbn = '9780441172719'
ON CONFLICT DO NOTHING;

INSERT INTO copies (book_id, inventory_no, status)
SELECT b.id, 'INV-' || v.n, CASE WHEN v.n % 5 = 0 THEN 'lost' ELSE 'available' END
FROM (
    VALUES
        (1, 'INV-001'), (2, 'INV-002'), (3, 'INV-003'), (4, 'INV-004'), (5, 'INV-005'),
        (6, 'INV-006'), (7, 'INV-007'), (8, 'INV-008'), (9, 'INV-009'), (10, 'INV-010'),
        (11, 'INV-011'), (12, 'INV-012'), (13, 'INV-013'), (14, 'INV-014'), (15, 'INV-015'),
        (16, 'INV-016'), (17, 'INV-017'), (18, 'INV-018'), (19, 'INV-019'), (20, 'INV-020')
) AS v(n, inventory_no)
JOIN books b ON b.id = ((v.n - 1) % (SELECT count(*) FROM books)) + 1
ON CONFLICT (inventory_no) DO NOTHING;

UPDATE copies
SET status = CASE WHEN id % 4 = 0 THEN 'damaged' ELSE status END
WHERE id % 4 = 0;
