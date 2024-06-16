ALTER TABLE
    "public"."projects"
ADD
    COLUMN "stops" jsonb NOT NULL DEFAULT '[]' :: jsonb;