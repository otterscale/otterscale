<script lang="ts">
	import { type FormState, getValueSnapshot } from '@sjsf/form';

	import {
		type K8sOpenAPISchema,
		type PathOptions,
		SchemaForm
	} from '$lib/components/custom/schema-form';

	import cronSchema from './cron_api.json';

	let form = $state<FormState<Record<string, unknown>> | undefined>();
	let mode = $state<'basic' | 'advance'>('basic');
	// Default values for CronJob
	const initialData = {
		spec: {
			schedule: '0 0 * * *',
			concurrencyPolicy: 'Allow',
			suspend: false,
			jobTemplate: {
				spec: {
					template: {
						spec: {
							// restartPolicy: 'OnFailure'
							containers: [
								{
									name: 'job-container',
									image: 'busybox'
									// command: ['sh', '-c', 'echo Hello World']
								}
							]
						}
					}
				}
			}
		}
	};

	const fields: Record<string, PathOptions> = {
		// 'metadata.namespace': { title: 'Namespace' },
		'metadata.name': { title: 'Name' },
		// 'metadata.annotations': { title: 'Annotations' },
		// apiVersion: {},
		// 'spec.schedule': { title: 'Cron Schedule', showDescription: true },
		// 'spec.concurrencyPolicy': { title: 'Concurrency Policy' },
		'spec.timeZone': {},
		'spec.startingDeadlineSeconds': {},
		'spec.successfulJobsHistoryLimit': {},
		'spec.failedJobsHistoryLimit': {},
		'spec.suspend': {
			uiSchema: {
				'ui:components': {
					checkboxWidget: 'switchWidget'
				}
			}
		},
		'spec.jobTemplate.spec.template.spec.containers': { title: 'Containers' },
		'spec.jobTemplate.spec.template.spec.containers.name': {},
		'spec.jobTemplate.spec.template.spec.containers.image': {},
		'spec.jobTemplate.spec.template.spec.containers.imagePullPolicy': {}
		// 'spec.jobTemplate.spec.template.spec.containers.resources.requests': {}
		// 'spec.jobTemplate.spec.template.spec.containers.lifecycle.postStart': {}
	};

	// Using $derived to reactively get values from the form store
	const formValues = $derived(form ? getValueSnapshot(form) : {});

	const fieldKeys = Object.keys(fields);
</script>

<div class="container mx-auto py-10">
	<h1 class="mb-4 text-2xl font-bold">CronJob Form</h1>

	<div class="grid grid-cols-2 gap-8">
		<div class="rounded border bg-card p-4 text-card-foreground">
			<h2 class="mb-4 text-xl">Generated Form (Mode: {mode})</h2>
			<SchemaForm
				apiSchema={cronSchema as K8sOpenAPISchema}
				paths={fields}
				bind:form
				bind:mode
				yamlEditable={true}
			/>
		</div>

		<div class="rounded border bg-muted/50 p-4">
			<h2 class="mb-4 text-xl">Live Values</h2>
			<pre class="overflow-auto rounded bg-zinc-950 p-4 text-xs text-zinc-50 dark:bg-zinc-900">
{JSON.stringify(formValues, null, 2)}
            </pre>

			<h2 class="mt-4 mb-2 text-xl">Selected Paths</h2>
			<ul class="list-inside list-disc text-sm">
				{#each fieldKeys as field (field)}
					<li><code>{field}</code></li>
				{/each}
			</ul>
		</div>
	</div>
</div>
