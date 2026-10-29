CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS authors (
    id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    bio TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS books (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    isbn VARCHAR(64) NOT NULL UNIQUE,
    category_id BIGINT NOT NULL,
    published_year INTEGER,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_books_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS book_authors (
    book_id BIGINT NOT NULL,
    author_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (book_id, author_id),
    CONSTRAINT fk_book_authors_book
        FOREIGN KEY (book_id)
        REFERENCES books(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_book_authors_author
        FOREIGN KEY (author_id)
        REFERENCES authors(id)
        ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS copies (
    id BIGSERIAL PRIMARY KEY,
    book_id BIGINT NOT NULL,
    inventory_no VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(32) NOT NULL DEFAULT 'available',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_copies_book
        FOREIGN KEY (book_id)
        REFERENCES books(id)
        ON DELETE CASCADE,
    CONSTRAINT chk_copies_status
        CHECK (status IN ('available', 'lost', 'damaged'))
);

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(32) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_users_role
        CHECK (role IN ('user', 'admin'))
);

CREATE INDEX IF NOT EXISTS idx_books_category_id ON books (category_id);
CREATE INDEX IF NOT EXISTS idx_book_authors_book_id ON book_authors (book_id);
CREATE INDEX IF NOT EXISTS idx_book_authors_author_id ON book_authors (author_id);
CREATE INDEX IF NOT EXISTS idx_copies_book_id ON copies (book_id);
