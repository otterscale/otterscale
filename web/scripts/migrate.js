import { drizzle } from 'drizzle-orm/node-postgres';
import { migrate } from 'drizzle-orm/node-postgres/migrator';
import path from 'path';
import { Pool } from 'pg';

const config = {
	databaseUrl: process.env.DATABASE_URL,
	migrationsFolder: path.resolve(process.cwd(), './drizzle')
};

async function runMigrations() {
	const { databaseUrl, migrationsFolder } = config;

	if (!databaseUrl) {
		console.log('DATABASE_URL is not set. Skipping migrations.');
		return;
	}

	console.log('Starting database migrations...');
	console.log(`Migrations folder path: ${migrationsFolder}`);

	const pool = new Pool({
		connectionString: databaseUrl
	});

	try {
		const db = drizzle(pool);
		await migrate(db, { migrationsFolder });
		console.log('Database migrations completed successfully!');
	} catch (error) {
		console.error('Database migrations failed!');
		console.error(error);
		process.exit(1);
	} finally {
		await pool.end();
		console.log('Database connection pool closed.');
	}
}

(async () => {
	await runMigrations();
})();
