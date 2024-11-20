create table if not exists "open_files" (
  "md5" TEXT NOT NULL,
  "url" INTEGER,
  PRIMARY KEY ("md5")
);

create table if not exists "open_file_frames" (
  "md5" TEXT NOT NULL,
  "url" INTEGER,
  "frames" text,
  PRIMARY KEY ("md5")
);