<script lang="ts">
	import type { FormState } from '@sjsf/form';
	import { BasicForm, createForm, getValueSnapshot } from '@sjsf/form';
	import { setThemeContext } from '@sjsf/shadcn4-theme';
	import * as components from '@sjsf/shadcn4-theme/new-york';
	import yaml from 'js-yaml';
	import { onMount } from 'svelte';
	import Monaco from 'svelte-monaco';

	import * as Tabs from '$lib/components/ui/tabs';

	import { buildSchemaFromK8s, type K8sOpenAPISchema } from './converter';
	import * as defaults from './defaults';

	interface Props {
		/** The full K8s OpenAPI V3 Schema */
		apiSchema: K8sOpenAPISchema;
		/** Paths to include in basic mode (dot notation, e.g. "spec.running") */
		paths: string[];
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

	// Build subset schema for basic mode
	const formConfig = $derived(buildSchemaFromK8s(apiSchema, paths));

	// Full schema for advance mode (YAML editor)
	let advanceYaml = $state('');

	// Initialize form on mount
	onMount(() => {
		initializeForm(initialData ?? formConfig.initialValue);
	});

	// Helper to create form instance
	function initializeForm(value: Record<string, unknown>) {
		form = createForm<Record<string, unknown>>({
			...defaults,
			idPrefix: 'k8s-schema-form',
			initialValue: value,
			schema: formConfig.schema,
			uiSchema: formConfig.uiSchema,
			onSubmit: (data) => console.log('Form submitted', data)
		});
	}

	// Sync YAML back to form when switching from advance to basic
	function syncYamlToForm() {
		try {
			const parsed = yaml.load(advanceYaml) as Record<string, unknown> | null;
			if (parsed && typeof parsed === 'object') {
				initializeForm(parsed);
			}
		} catch (e) {
			console.error('Failed to parse YAML:', e);
		}
	}

	function handleModeChange(newMode: string) {
		// Sync form values to YAML when switching to advance mode
		if (newMode === 'advance' && form) {
			const currentValues = getValueSnapshot(form);
			advanceYaml = yaml.dump(currentValues, { indent: 2, lineWidth: -1 });
		}
		// Sync YAML back to form when switching to basic mode
		if (newMode === 'basic' && advanceYaml) {
			syncYamlToForm();
		}
		mode = newMode as 'basic' | 'advance';
		onModeChange?.(mode);
	}

	// Set theme context for this component tree
	setThemeContext({ components });
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
				<BasicForm {form} />
			{/if}
		</Tabs.Content>

		<Tabs.Content value="advance" class="h-[50vh]">
			<div class="h-full rounded border">
				<Monaco
					options={{
						language: 'yaml',
						padding: { top: 16, bottom: 8 },
						automaticLayout: true,
						minimap: { enabled: false },
						scrollBeyondLastLine: false
					}}
					theme="vs-dark"
					bind:value={advanceYaml}
				/>
			</div>
		</Tabs.Content>
	</Tabs.Root>
</div>
