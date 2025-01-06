CREATE TABLE usage_categories
(
    category_type usage_category_type NOT NULL PRIMARY KEY,
    display_order SMALLINT            NOT NULL,
    created_at    TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_usage_categories_display_order UNIQUE (display_order) DEFERRABLE INITIALLY DEFERRED
);
