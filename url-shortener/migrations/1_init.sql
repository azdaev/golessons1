CREATE TABLE links (
    id bigserial primary key,
    long_link text NOT NULL,
    short_link text NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL
);

CREATE TABLE redirects (
    id bigserial primary key,
    long_link text NOT NULL,
    short_link text NOT NULL,
    user_agent text NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL
);
