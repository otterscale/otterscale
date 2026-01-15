<script lang="ts">
	import { type FormState, getValueSnapshot } from '@sjsf/form';

	import { type K8sOpenAPISchema, SchemaForm } from '$lib/components/custom/schema-form';

	import vmSchema from './vm_api.json';

	let form = $state<FormState<Record<string, unknown>> | undefined>();
	let mode = $state<'basic' | 'advance'>('basic');

	const fields = [
		'spec.runStrategy',
		'spec.running',
		'spec.instancetype.name',
		'spec.instancetype.kind',
		'spec.preference.name'
	];

	// Using $derived to reactively get values from the form store
	const formValues = $derived(form ? getValueSnapshot(form) : {});
</script>

<div class="container mx-auto py-10">
	<h1 class="mb-4 text-2xl font-bold">Schema Form Gen Experiment</h1>

	<div class="grid grid-cols-2 gap-8">
		<div class="rounded border bg-white p-4">
			<h2 class="mb-4 text-xl">Generated Form (Mode: {mode})</h2>
			<SchemaForm apiSchema={vmSchema as K8sOpenAPISchema} paths={fields} bind:form bind:mode />
		</div>

		<div class="rounded border bg-gray-50 p-4">
			<h2 class="mb-4 text-xl">Live Values</h2>
			<pre class="overflow-auto rounded bg-gray-900 p-4 text-green-400">
{JSON.stringify(formValues, null, 2)}
            </pre>

			<h2 class="mt-4 mb-2 text-xl">Selected Paths</h2>
			<ul class="list-inside list-disc text-sm">
				{#each fields as field (field)}
					<li><code>{field}</code></li>
				{/each}
			</ul>
		</div>
	</div>
</div>
