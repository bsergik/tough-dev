CREATE TABLE "tasks" (
  "id" INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "title" varchar(300) NOT NULL,
  "description" text,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "user_has_tasks" (
  "id" INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "task_id" integer NOT NULL,
  "user_id" integer NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "statuses" (
  "id" INTEGER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "name" varchar(50) NOT NULL,
  "description" text,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "task_has_statuses" (
  "task_id" integer NOT NULL,
  "status_id" integer NOT NULL,
  "user_id" integer NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "schema_version" (
  "version" varchar(20),
  "dirty" boolean DEFAULT false
);

COMMENT ON COLUMN "task_has_statuses"."user_id" IS 'Who changed status';

ALTER TABLE "user_has_tasks" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

ALTER TABLE "task_has_statuses" ADD FOREIGN KEY ("task_id") REFERENCES "tasks" ("id");

ALTER TABLE "task_has_statuses" ADD FOREIGN KEY ("status_id") REFERENCES "statuses" ("id");
