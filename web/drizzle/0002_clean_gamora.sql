ALTER TABLE "sessions" ADD COLUMN "access_token" text;--> statement-breakpoint
ALTER TABLE "sessions" ADD COLUMN "refresh_token" text;--> statement-breakpoint
ALTER TABLE "sessions" ADD COLUMN "refresh_token_expires_at" timestamp with time zone;