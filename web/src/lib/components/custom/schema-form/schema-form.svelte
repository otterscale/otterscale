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
		type PathOptions
	} from './converter';
	import * as defaults from './defaults';
	import SchemaFormStep from './SchemaFormStep.svelte';

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
		/** Form title */
		title?: string;
		/** Callback when mode changes */
		onModeChange?: (mode: 'basic' | 'advance') => void;
		/** Callback when submit button is clicked */
		onSubmit?: (data: Record<string, unknown>) => Promise<void> | void;
		/** Transform data before submission */
		transformData?: (data: Record<string, unknown>) => Record<string, unknown>;
		/** Whether YAML editor is editable */
		yamlEditable?: boolean;
	}

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
	let masterData = $state<Record<string, unknown>>({});
	let advanceYaml = $state('');
	let yamlParseError = $state<string | null>(null);
	let ref: HTMLFormElement | undefined;

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

	function syncFormToYaml() {
		try {
			const rawData = form ? getValueSnapshot(form) : initialValue;
			const k8sData = formDataToK8s(rawData, formConfig.transformationMappings);
			masterData = k8sData;

			advanceYaml = yaml.dump(k8sData, {
				indent: 2,
				lineWidth: -1
			});
		} catch (error) {
			console.error(`Error during syncing form to YAML:`, error);
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
			const errorMsg = `Invalid YAML: ${error instanceof Error ? error.message : 'Unknown error'}`;
			yamlParseError = errorMsg;
			console.error(`Error during parsing YAML:`, error);
		}
	}

	function submitFinalData(k8sData: Record<string, unknown>) {
		try {
			if (transformData) {
				k8sData = transformData(k8sData);
			}
			masterData = k8sData;

			if (onSubmit) {
				onSubmit(k8sData);
			} else {
				onModeChange?.(mode);
			}
		} catch (error) {
			console.error(`Error during form submission:`, error);
		}
	}

	function handleFormSubmit(data: Record<string, unknown>) {
		try {
			const k8sData = formDataToK8s(data, formConfig.transformationMappings);
			submitFinalData(k8sData);
		} catch (error) {
			console.error(`Error during form submission:`, error);
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

	<Button
		class="mt-6 w-full"
		onclick={() => {
			if (mode === 'advance') {
				syncYamlToForm();
				if (masterData) {
					submitFinalData(masterData);
				}
			} else {
				ref?.requestSubmit();
			}
		}}
	>
		Submit
	</Button>
</div>
