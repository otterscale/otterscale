<script lang="ts">
	import type { FormState } from '@sjsf/form';
	import { createForm, getValueSnapshot } from '@sjsf/form';
	import yaml from 'js-yaml';
	import { mode as themeMode } from 'mode-watcher';
	import Monaco from 'svelte-monaco';

	import { Button } from '$lib/components/ui/button';
	import * as Tabs from '$lib/components/ui/tabs';

	import {
		buildSchemaFromK8s,
		formDataToK8s,
		type K8sOpenAPISchema,
		k8sToFormData,
		normalizeArrays,
		type PathOptions
	} from './converter';
	import * as defaults from './defaults';
	import SchemaFormStep from './schema-form-step.svelte';
	import { deepMerge } from './utils';

	// ── Types ──────────────────────────────────────────────────

	interface Props {
		apiSchema: K8sOpenAPISchema;
		paths: string[] | Record<string, PathOptions>;
		form?: FormState<Record<string, unknown>>;
		initialData?: Record<string, unknown>;
		mode?: 'basic' | 'advance';
		title?: string;
		onModeChange?: (mode: 'basic' | 'advance') => void;
		onSubmit?: (data: Record<string, unknown>) => Promise<void> | void;
		transformData?: (data: Record<string, unknown>) => Record<string, unknown>;
		yamlEditable?: boolean;
	}

	// ── Props & State ──────────────────────────────────────────

	let {
		apiSchema,
		paths,
		form = $bindable(),
		initialData,
		mode = $bindable('basic'),
		title,
		onModeChange,
		onSubmit,
		transformData,
		yamlEditable = false
	}: Props = $props();

	const formConfig = buildSchemaFromK8s(apiSchema, paths);

	let initialValue = $state(k8sToFormData(initialData, formConfig.transformationMappings));
	let masterData = $state<Record<string, unknown>>(initialData ?? {});
	let advanceYaml = $state('');
	let yamlParseError = $state<string | null>(null);
	let ref: HTMLFormElement | undefined;

	// ── Helpers ────────────────────────────────────────────────

	/** Merge form data into masterData after converting to K8s format (preserves data not in form paths) */
	function mergeIntoMaster(data: Record<string, unknown>) {
		const k8sData = formDataToK8s(data, formConfig.transformationMappings);
		masterData = normalizeArrays(deepMerge(masterData, k8sData)) as Record<string, unknown>;
	}

	/** Apply transformData callback if provided */
	function applyTransform() {
		if (transformData) {
			masterData = transformData(masterData);
		}
	}

	// ── Form Lifecycle ─────────────────────────────────────────

	function initForm(data: Record<string, unknown>) {
		return createForm<Record<string, unknown>>({
			...defaults,
			idPrefix: 'k8s-schema-form',
			initialValue: data,
			schema: formConfig.schema,
			uiSchema: formConfig.uiSchema,
			onSubmit: handleFormSubmit
		});
	}

	form = initForm(initialValue);

	// ── Event Handlers ─────────────────────────────────────────

	function handleFormSubmit(data: Record<string, unknown>) {
		const k8sData = formDataToK8s(data, formConfig.transformationMappings);
		submitFinalData(k8sData);
	}

	function submitFinalData(data: Record<string, unknown>) {
		masterData = data;
		applyTransform();

		if (onSubmit) {
			onSubmit(masterData);
		} else {
			onModeChange?.(mode);
		}
	}

	function handleSubmitClick() {
		if (mode === 'advance') {
			syncYamlToForm();
			submitFinalData(masterData);
		} else {
			ref?.requestSubmit();
		}
	}

	function handleModeChange(newMode: string) {
		const targetMode = newMode as 'basic' | 'advance';

		if (targetMode === 'basic' && mode === 'advance') {
			syncYamlToForm();
		}

		mode = targetMode;
		onModeChange?.(mode);
	}

	// ── YAML Sync ──────────────────────────────────────────────

	function syncFormToYaml() {
		try {
			const rawData = form ? (getValueSnapshot(form) as Record<string, unknown>) : initialValue;
			mergeIntoMaster(rawData);
			advanceYaml = yaml.dump(masterData, { indent: 2, lineWidth: -1 });
		} catch (error) {
			console.error('Error syncing form to YAML:', error);
		}
	}

	function syncYamlToForm() {
		try {
			yamlParseError = null;
			const parsed = yaml.load(advanceYaml) as Record<string, unknown> | null;

			if (parsed && typeof parsed === 'object') {
				const formData = k8sToFormData(parsed, formConfig.transformationMappings);
				Object.assign(initialValue, formData);
				form = initForm(initialValue);
				masterData = parsed;
			}
		} catch (error) {
			yamlParseError = `Invalid YAML: ${error instanceof Error ? error.message : 'Unknown error'}`;
			console.error('Error parsing YAML:', error);
		}
	}

	$effect(() => {
		if (form && mode === 'basic') {
			syncFormToYaml();
		}
	});
</script>

<div class="schema-form-container">
	<div class="relative mb-4 flex items-center justify-center py-2">
		{#if title}
			<h1 class="text-2xl font-bold">{title}</h1>
		{/if}
		<div class="absolute right-0">
			<Tabs.Root value={mode} onValueChange={handleModeChange}>
				<Tabs.List>
					<Tabs.Trigger value="basic">Form</Tabs.Trigger>
					<Tabs.Trigger value="advance">YAML</Tabs.Trigger>
				</Tabs.List>
			</Tabs.Root>
		</div>
	</div>

	<Tabs.Root value={mode}>
		<Tabs.Content value="basic">
			{#if form}
				{#key form}
					<SchemaFormStep {form} bind:ref />
				{/key}
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
						scrollBeyondLastLine: false,
						readOnly: !yamlEditable
					}}
					theme={themeMode.current === 'dark' ? 'vs-dark' : 'vs'}
					bind:value={advanceYaml}
				/>
			</div>
		</Tabs.Content>
	</Tabs.Root>

	<Button class="mt-6 w-full" onclick={handleSubmitClick}>Submit</Button>
</div>
