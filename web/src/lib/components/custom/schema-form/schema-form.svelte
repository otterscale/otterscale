<script lang="ts">
	import type { FormState } from '@sjsf/form';
	import { Content, createForm, Form, getValueSnapshot, setFormContext } from '@sjsf/form';
	import { setThemeContext } from '@sjsf/shadcn4-theme';
	import * as components from '@sjsf/shadcn4-theme/new-york';
	import yaml from 'js-yaml';
	import Monaco from 'svelte-monaco';

	import { Button } from '$lib/components/ui/button';
	import * as Tabs from '$lib/components/ui/tabs';
	import { mode as themeMode } from 'mode-watcher';

	import { buildSchemaFromK8s, type K8sOpenAPISchema, type PathOptions } from './converter';
	import * as defaults from './defaults';

	interface Props {
		/** The full K8s OpenAPI V3 Schema */
		apiSchema: K8sOpenAPISchema;
		/** Paths to include in basic mode (dot notation, e.g. "spec.running") */
		paths: string[] | Record<string, PathOptions>;
		/** Allow binding the form instance back to parent */
		form?: FormState<Record<string, unknown>>;
		/** Optional initial value override */
		initialData?: Record<string, unknown>;
		/** Current mode: 'basic' | 'advance' */
		mode?: 'basic' | 'advance';
		/** Callback when mode changes */
		onModeChange?: (mode: 'basic' | 'advance') => void;
	}

	let {
		apiSchema,
		paths,
		form = $bindable(),
		initialData,
		mode = $bindable('basic'),
		onModeChange
	}: Props = $props();

	// Set theme context for this component tree
	setThemeContext({ components });

	// Build subset schema for basic mode
	const formConfig = $derived(buildSchemaFromK8s(apiSchema, paths));

	// Reactive states
	let initialValue = $state(initialData ?? formConfig.initialValue);
	let advanceYaml = $state('');
	let yamlParseError = $state<string | null>(null);
	let ref: HTMLFormElement | undefined;

	// Create form instance
	form = createForm<Record<string, unknown>>({
		...defaults,
		idPrefix: 'k8s-schema-form',
		initialValue,
		schema: formConfig.schema,
		uiSchema: formConfig.uiSchema,
		onSubmit: handleFormSubmit
	});
	setFormContext(form);
	syncFormToYaml();

	// Sync form values to YAML editor
	function syncFormToYaml() {
		try {
			advanceYaml = yaml.dump(form ? getValueSnapshot(form) : initialValue, {
				indent: 2,
				lineWidth: -1
			});
		} catch (error) {
			console.error(`Error during syncing form to YAML:`, error);
		}
	}

	// Sync YAML back to form
	function syncYamlToForm() {
		try {
			yamlParseError = null;
			const parsed = yaml.load(advanceYaml) as Record<string, unknown> | null;

			if (parsed && typeof parsed === 'object') {
				Object.assign(initialValue, parsed);
			}
		} catch (error) {
			const errorMsg = `Invalid YAML: ${error instanceof Error ? error.message : 'Unknown error'}`;
			yamlParseError = errorMsg;
			console.error(`Error during parsing YAML:`, error);
		}
	}

	// Handle form submission
	function handleFormSubmit(data: Record<string, unknown>) {
		try {
			// Display submitted form data
			console.log('Form submitted with data:', data);

			// Custom submission logic can be added here
			onModeChange?.(mode);
		} catch (error) {
			console.error(`Error during form submission:`, error);
		}
	}

	// Handle mode changes
	function handleModeChange(newMode: string) {
		const targetMode = newMode as 'basic' | 'advance';

		if (targetMode === 'basic' && mode === 'advance') {
			syncYamlToForm();
		}

		mode = targetMode;
		onModeChange?.(mode);
	}

	// Reactive effects
	$effect(() => {
		if (form && mode === 'basic') {
			syncFormToYaml();
		}
	});
</script>

<div class="schema-form-container">
	<div class="mb-4 flex items-center justify-end">
		<Tabs.Root value={mode} onValueChange={handleModeChange}>
			<Tabs.List>
				<Tabs.Trigger value="basic">Basic</Tabs.Trigger>
				<Tabs.Trigger value="advance">Advance</Tabs.Trigger>
			</Tabs.List>
		</Tabs.Root>
	</div>

	<Tabs.Root value={mode}>
		<Tabs.Content value="basic">
			{#if form}
				<Form bind:ref>
					<Content />
				</Form>
			{/if}
		</Tabs.Content>

		<Tabs.Content value="advance">
			<div class="h-[70vh] rounded border">
				{#if yamlParseError}
					<div
						class="mb-2 rounded bg-red-100 p-2 text-sm text-red-700 dark:bg-red-900/20 dark:text-red-400"
					>
						{yamlParseError}
					</div>
				{/if}

				<Monaco
					options={{
						language: 'yaml',
						padding: { top: 16, bottom: 8 },
						automaticLayout: true,
						minimap: { enabled: false },
						scrollBeyondLastLine: false
					}}
					theme={themeMode.current === 'dark' ? 'vs-dark' : 'vs'}
					bind:value={advanceYaml}
				/>
			</div>
		</Tabs.Content>
	</Tabs.Root>

	<Button
		class="mt-6 w-full"
		onclick={() => {
			if (mode === 'advance') {
				syncYamlToForm();
			}
			ref?.requestSubmit();
		}}
	>
		Submit
	</Button>
</div>
