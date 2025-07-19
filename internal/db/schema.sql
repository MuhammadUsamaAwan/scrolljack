CREATE TABLE `modlists` (
	`id` text PRIMARY KEY NOT NULL,
	`name` text NOT NULL,
	`author` text,
	`description` text,
	`website` text,
	`image` text,
	`readme` text,
	`game_type` text,
	`version` text,
	`is_nsfw` integer,
	`created_at` text DEFAULT (CURRENT_TIMESTAMP) NOT NULL
);

CREATE TABLE `profiles` (
	`id` text PRIMARY KEY NOT NULL,
	`modlist_id` text NOT NULL,
	`name` text NOT NULL,
	FOREIGN KEY (`modlist_id`) REFERENCES `modlists`(`id`) ON UPDATE no action ON DELETE cascade
);

CREATE INDEX `idx_profiles_modlist_id` ON `profiles` (`modlist_id`);

CREATE TABLE `profile_files` (
	`id` text PRIMARY KEY NOT NULL,
	`profile_id` text NOT NULL,
	`name` text NOT NULL,
	`file_path` text NOT NULL,
	FOREIGN KEY (`profile_id`) REFERENCES `profiles`(`id`) ON UPDATE no action ON DELETE cascade
);

CREATE INDEX `idx_profile_files_profile_id` ON `profile_files` (`profile_id`);

CREATE TABLE `mods` (
	`id` text PRIMARY KEY NOT NULL,
	`profile_id` text NOT NULL,
	`name` text NOT NULL,
	`is_separator` integer NOT NULL,
	`order` integer NOT NULL,
	`mod_order` integer NOT NULL,
	`is_active` integer NOT NULL,
	FOREIGN KEY (`profile_id`) REFERENCES `profiles`(`id`) ON UPDATE no action ON DELETE cascade
);

CREATE INDEX `idx_mods_profile_id_order` ON `mods` (`profile_id`,`order`);

CREATE TABLE `mod_files` (
	`id` text PRIMARY KEY NOT NULL,
	`mod_id` text NOT NULL,
    `hash` text NOT NULL,
	`type` text NOT NULL,
	`path` text NOT NULL,
	`source_file_path` text,
	`patch_file_path` text,
	`bsa_files` text,
	FOREIGN KEY (`mod_id`) REFERENCES `mods`(`id`) ON UPDATE no action ON DELETE cascade
);

CREATE INDEX `idx_mod_files_mod_id` ON `mod_files` (`mod_id`);

CREATE TABLE `mod_archives` (
	`id` text PRIMARY KEY NOT NULL,
	`mod_id` text NOT NULL,
	`archive_type` text,
	`nexus_game_name` text,
	`nexus_mod_id` integer,
	`nexus_file_id` integer,
	`direct_url` text,
	`version` text,
	`size` integer,
	`description` text,
	FOREIGN KEY (`mod_id`) REFERENCES `mods`(`id`) ON UPDATE no action ON DELETE cascade
);

CREATE INDEX `idx_mod_archives_mod_id` ON `mod_archives` (`mod_id`);

CREATE TABLE `mod_file_archives` (
	`modlist_id` text NOT NULL,
	`mod_file_id` text NOT NULL,
	`mod_archive_id` text NOT NULL,
	PRIMARY KEY(`mod_file_id`, `mod_archive_id`, `modlist_id`),
	FOREIGN KEY (`modlist_id`) REFERENCES `modlists`(`id`) ON UPDATE no action ON DELETE cascade,
	FOREIGN KEY (`mod_file_id`) REFERENCES `mod_files`(`id`) ON UPDATE no action ON DELETE cascade,
	FOREIGN KEY (`mod_archive_id`) REFERENCES `mod_archives`(`id`) ON UPDATE no action ON DELETE cascade
);
