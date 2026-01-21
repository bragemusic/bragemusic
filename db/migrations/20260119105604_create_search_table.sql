-- migrate:up
CREATE VIRTUAL TABLE search_items
USING FTS5(name, id, type, link_id, link_type);


-- migrate:down

