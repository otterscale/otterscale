<script lang="ts">
	import {
		type GroupedFields,
		type K8sOpenAPISchema,
		MultiStepSchemaForm
	} from '$lib/components/custom/schema-form';

	import cronSchema from './cron_api.json';

	// Grouped fields for multi-step form
	const groupedFields: GroupedFields = {
		// Step 1: Basic Metadata
		'Basic Metadata': {
			'metadata.namespace': { title: 'Namespace' },
			'metadata.name': { title: 'Name' },
			'metadata.labels.app': { title: 'App Label' },
		},
		// Step 2: CronJob Spec
		'CronJob Spec': {
			'spec.schedule': { title: 'Cron Schedule', showDescription: true },
			'spec.timeZone': { title: 'Time Zone' },
			'spec.suspend': {
				title: 'Suspend',
				uiSchema: {
					'ui:components': {
						checkboxWidget: 'switchWidget'
					}
				}
			},
			'spec.concurrencyPolicy': { title: 'Concurrency Policy' },
			'spec.successfulJobsHistoryLimit': { title: 'Successful History Limit' },
			'spec.failedJobsHistoryLimit': { title: 'Failed History Limit' }
		},
		// Step 3: Job Template
		'Job Template': {
			'spec.jobTemplate.spec.template.spec.restartPolicy': { title: 'Restart Policy' }
		},
		// Step 4: Container Spec
		'Container Spec': {
			'spec.jobTemplate.spec.template.spec.containers': { title: 'Containers' },
			'spec.jobTemplate.spec.template.spec.containers.name': { title: 'Container Name' },
			'spec.jobTemplate.spec.template.spec.containers.image': { title: 'Image' },
			'spec.jobTemplate.spec.template.spec.containers.command': { title: 'Command' },
			'spec.jobTemplate.spec.template.spec.containers.args': { title: 'Arguments' },
			'spec.jobTemplate.spec.template.spec.containers.env': { title: 'Environment Variables' },
			'spec.jobTemplate.spec.template.spec.containers.resources.requests.cpu': { title: 'Requests CPU' },
			'spec.jobTemplate.spec.template.spec.containers.resources.requests.memory': { title: 'Requests Memory' },
			'spec.jobTemplate.spec.template.spec.containers.resources.limits.cpu': { title: 'Limits CPU' },
			'spec.jobTemplate.spec.template.spec.containers.resources.limits.memory': { title: 'Limits Memory' },
			'spec.jobTemplate.spec.template.spec.containers.imagePullPolicy': {
				title: 'Image Pull Policy'
			}
		}
	};

	function handleMultiStepSubmit(data: Record<string, unknown>) {
		console.log('Multi-step form submitted:', data);
		alert('Form submitted! Check console for data.');
	}
</script>

<div class="container mx-auto py-10">
	<h1 class="mb-4 text-2xl font-bold">CronJob Form</h1>

	<!-- Multi-Step Form -->
	<div class="mb-12">
		<h2 class="mb-4 text-xl font-semibold">Multi-Step Schema Form</h2>
		<div class="rounded border bg-card p-4 text-card-foreground">
			<MultiStepSchemaForm
				apiSchema={cronSchema as K8sOpenAPISchema}
				fields={groupedFields}
				onSubmit={handleMultiStepSubmit}
			/>
		</div>
	</div>
</div>
