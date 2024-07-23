create table if not exists "open_files" (
  "md5" TEXT NOT NULL,
  "url" INTEGER,
  PRIMARY KEY ("md5")
);