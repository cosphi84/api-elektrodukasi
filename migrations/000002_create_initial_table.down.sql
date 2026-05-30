DROP INDEX IF EXISTS idx_article_hits_article_created;
DROP INDEX IF EXISTS idx_comments_article;
DROP INDEX IF EXISTS idx_article_tags_tag;
DROP INDEX IF EXISTS idx_articles_category;
DROP INDEX IF EXISTS idx_articles_published_at;

DROP TABLE IF EXISTS article_hits;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS article_tags;
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;