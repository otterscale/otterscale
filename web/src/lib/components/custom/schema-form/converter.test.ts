import { describe, expect, it } from 'vitest';

import { cronjobSchema, workspaceSchema } from './__fixtures__';
import { buildSchemaFromK8s, type PathOptions } from './converter';

/**
 * Test paths extracted from workspace-form/+page.svelte groupedFields.
 * Note: Svelte components (like UserSelectWidget) are omitted as they cannot be
 * imported in node-based tests. The uiSchema config is preserved.
 */
const workspaceFormPaths: Record<string, PathOptions> = {
	// Step 1: Workspace & Users
	'metadata.name': { title: 'Workspace Name' },
	'spec.namespace': { title: 'Namespace', showDescription: true },
	'spec.users': {
		title: 'Users',
		uiSchema: {
			items: {
				'ui:components': {
					objectField: 'UserSelectWidget' // String placeholder for Svelte component
				}
			}
		}
	},
	// Step 2: Network Isolation
	'spec.networkIsolation': { title: 'Network Isolation' },
	'spec.networkIsolation.enabled': {
		title: 'Enable Network Isolation',
		uiSchema: {
			'ui:components': {
				checkboxWidget: 'switchWidget'
			}
		}
	},
	'spec.networkIsolation.allowedNamespaces': { title: 'Allowed Namespaces' },
	// Step 3: Default Resource Settings
	'spec.resourceQuota.hard.requests.cpu': { title: 'Requests CPU', disabled: true },
	'spec.resourceQuota.hard.requests.memory': { title: 'Requests Memory', disabled: true },
	'spec.resourceQuota.hard.requests.nvidia.com/gpu': { title: 'Requests GPU', disabled: true },
	'spec.resourceQuota.hard.limits.cpu': { title: 'Limits CPU', disabled: true },
	'spec.resourceQuota.hard.limits.memory': { title: 'Limits Memory', disabled: true },
	'spec.resourceQuota.hard.limits.nvidia.com/gpu': { title: 'Limits GPU', disabled: true },
	'spec.limitRange.limits': {
		title: 'Limit Range',
		uiSchema: {
			'ui:options': {
				addable: false,
				removable: false,
				orderable: false
			}
		}
	},
	'spec.limitRange.limits.type': { title: 'Type', disabled: true },
	'spec.limitRange.limits.default.cpu': { title: 'Default CPU Limit', disabled: true },
	'spec.limitRange.limits.default.memory': { title: 'Default Memory Limit', disabled: true },
	'spec.limitRange.limits.defaultRequest.cpu': { title: 'Default CPU Request', disabled: true },
	'spec.limitRange.limits.defaultRequest.memory': {
		title: 'Default Memory Request',
		disabled: true
	}
};

describe('buildSchemaFromK8s', () => {
	it('should generate expected schema, uiSchema and transformationMappings for workspace form', () => {
		const result = buildSchemaFromK8s(workspaceSchema, workspaceFormPaths);

		// Snapshot tests - these will fail if output changes between versions
		expect(result.schema).toMatchSnapshot('schema');
		expect(result.uiSchema).toMatchSnapshot('uiSchema');
		expect(result.transformationMappings).toMatchSnapshot('transformationMappings');
	});

	it('should handle simple array of paths', () => {
		const paths = ['metadata.name', 'spec.namespace'];
		const result = buildSchemaFromK8s(workspaceSchema, paths);

		expect(result.schema).toMatchSnapshot('simple-paths-schema');
		expect(result.uiSchema).toMatchSnapshot('simple-paths-uiSchema');
		expect(result.transformationMappings).toMatchSnapshot('simple-paths-transformationMappings');
	});

	it('should handle nested array paths correctly', () => {
		const paths: Record<string, PathOptions> = {
			'spec.users': { title: 'Users' },
			'spec.users.subject': { title: 'Subject' },
			'spec.users.role': { title: 'Role' }
		};
		const result = buildSchemaFromK8s(workspaceSchema, paths);

		expect(result.schema).toMatchSnapshot('array-paths-schema');
		expect(result.uiSchema).toMatchSnapshot('array-paths-uiSchema');
	});

	it('should handle map (additionalProperties) fields correctly', () => {
		const paths: Record<string, PathOptions> = {
			'spec.resourceQuota.hard.requests.cpu': { title: 'Requests CPU' }
		};
		const result = buildSchemaFromK8s(workspaceSchema, paths);

		expect(result.schema).toMatchSnapshot('map-paths-schema');
		expect(result.uiSchema).toMatchSnapshot('map-paths-uiSchema');
		expect(result.transformationMappings).toMatchSnapshot('map-paths-transformationMappings');
	});
});

/**
 * CronJob form paths from create-form.svelte groupedFields
 */
const cronjobFormPaths: Record<string, PathOptions> = {
	// Step 1: General Settings
	'metadata.name': { title: 'Name' },
	'spec.schedule': { title: 'Schedule', showDescription: true },
	'spec.concurrencyPolicy': { title: 'Concurrency Policy' },
	'spec.suspend': {
		title: 'Suspend execution',
		uiSchema: {
			'ui:components': {
				checkboxWidget: 'switchWidget'
			}
		}
	},
	// Step 2: Container Settings
	'spec.jobTemplate.spec.template.spec.containers.name': { title: 'Name' },
	'spec.jobTemplate.spec.template.spec.containers.image': { title: 'Image' },
	'spec.jobTemplate.spec.template.spec.containers.command': { title: 'Command' },
	'spec.jobTemplate.spec.template.spec.containers.args': { title: 'Arguments' },
	'spec.jobTemplate.spec.template.spec.containers.env': { title: 'Environment Variables' },
	// Step 3: Resources
	'spec.jobTemplate.spec.template.spec.containers.resources.requests.cpu': {
		title: 'Requests CPU'
	},
	'spec.jobTemplate.spec.template.spec.containers.resources.requests.memory': {
		title: 'Requests Memory'
	},
	'spec.jobTemplate.spec.template.spec.containers.resources.limits.cpu': {
		title: 'Limits CPU'
	},
	'spec.jobTemplate.spec.template.spec.containers.resources.limits.memory': {
		title: 'Limits Memory'
	}
};

describe('buildSchemaFromK8s - CronJob', () => {
	it('should generate expected schema, uiSchema and transformationMappings for cronjob form', () => {
		const result = buildSchemaFromK8s(cronjobSchema, cronjobFormPaths);

		expect(result.schema).toMatchSnapshot('cronjob-schema');
		expect(result.uiSchema).toMatchSnapshot('cronjob-uiSchema');
		expect(result.transformationMappings).toMatchSnapshot('cronjob-transformationMappings');
	});
});
