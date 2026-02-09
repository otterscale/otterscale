/**
 * Test fixtures for schema-form tests.
 *
 * This module provides helper functions to load JSON schema fixtures
 * for testing buildSchemaFromK8s and related functions.
 */

import { readFileSync } from 'fs';
import { dirname, resolve } from 'path';
import { fileURLToPath } from 'url';

import type { K8sOpenAPISchema } from '../converter';

const __dirname = dirname(fileURLToPath(import.meta.url));

/**
 * Loads a JSON schema fixture by filename.
 * @param filename - The fixture filename (e.g., 'workspace_api.json')
 * @returns The parsed K8s OpenAPI Schema
 */
export function loadSchemaFixture(filename: string): K8sOpenAPISchema {
	const filepath = resolve(__dirname, filename);
	return JSON.parse(readFileSync(filepath, 'utf-8'));
}

// Pre-loaded fixtures for convenience
export const workspaceSchema = loadSchemaFixture('workspace_api.json');
export const cronjobSchema = loadSchemaFixture('cron_api.json');

