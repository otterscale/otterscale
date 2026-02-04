<script lang="ts">
	import '@sjsf/basic-theme/css/basic.css';

	import Layout from '@sjsf/basic-theme/components/layout.svelte';
	import {
		Content,
		createForm,
		Form,
		getValueSnapshot,
		type Schema,
		setFormContext,
		type UiSchema
	} from '@sjsf/form';
	import { overrideByRecord } from '@sjsf/form/lib/resolver';
	import { setThemeContext } from '@sjsf/shadcn4-theme';
	import * as components from '@sjsf/shadcn4-theme/new-york';

	import * as defaults from '$lib/components/custom/schema-form/defaults';

	setThemeContext({ components });

	const schema: Schema = {
		type: 'object',
		properties: {
			spec: {
				type: 'object',
				properties: {
					jobTemplate: {
						type: 'object',
						properties: {
							spec: {
								type: 'object',
								properties: {
									template: {
										type: 'object',
										properties: {
											spec: {
												type: 'object',
												properties: {
													restartPolicy: {
														type: 'string',
														enum: ['Always', 'Never', 'OnFailure'],
														title: 'Restart Policy'
													},
													containers: {
														type: 'array',
														items: {
															type: 'object',
															properties: {
																name: {
																	type: 'string',
																	default: '',
																	title: 'Container Name'
																},
																image: {
																	type: 'string',
																	title: 'Image'
																},
																command: {
																	type: 'array',
																	items: {
																		default: '',
																		type: 'string'
																	},
																	title: 'Command'
																}
															},
															required: ['name']
														},
														title: 'Containers'
													}
												},
												required: ['containers']
											}
										}
									}
								},
								required: ['template']
							}
						}
					}
				},
				required: ['jobTemplate']
			}
		},
		required: []
	};

	const uiSchema: UiSchema = {};

	const initialValue = {
		spec: {
			jobTemplate: {
				spec: {
					template: {
						spec: {
							restartPolicy: 'OnFailure',
							containers: [
								{
									name: 'aaa',
									image: 'bbb',
									command: ['ccc']
								}
							]
						}
					}
				}
			}
		}
	};

	const form = createForm({
		...defaults,
		schema,
		uiSchema,
		initialValue,
		theme: overrideByRecord(defaults.theme, {
			layout: Layout
		})
	});
	setFormContext(form);

	// Reactive form values for display
	const formValues = $derived(getValueSnapshot(form));

	let ref: HTMLFormElement | undefined = $state();
</script>

<div class="container mx-auto py-10">
	<h1 class="mb-4 text-2xl font-bold">Simple Form Test</h1>

	<div class="grid grid-cols-2 gap-8">
		<div class="rounded border bg-card p-4 text-card-foreground">
			<h2 class="mb-4 text-xl">Form</h2>

			<Form bind:ref>
				<Content />
			</Form>

			<div class="mt-4 flex gap-2">
				<button
					class="rounded bg-primary px-4 py-2 text-primary-foreground hover:bg-primary/90"
					onclick={() => ref?.requestSubmit()}
				>
					Submit
				</button>
				<button
					class="rounded bg-secondary px-4 py-2 text-secondary-foreground hover:bg-secondary/90"
					onclick={() => ref?.reset()}
				>
					Reset
				</button>
			</div>
		</div>

		<div class="rounded border bg-muted/50 p-4">
			<h2 class="mb-4 text-xl">Live Values</h2>
			<pre class="overflow-auto rounded bg-zinc-950 p-4 text-xs text-zinc-50 dark:bg-zinc-900">
{JSON.stringify(formValues, null, 2)}
			</pre>
		</div>
	</div>
</div>
