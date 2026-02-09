<script lang="ts">
	import { Check } from '@lucide/svelte';
	import type { FormState } from '@sjsf/form';
	import { createForm, getValueSnapshot } from '@sjsf/form';
	import yaml from 'js-yaml';
	import { mode as themeMode } from 'mode-watcher';
	import Monaco from 'svelte-monaco';

	import { Button } from '$lib/components/ui/button';
	import * as Tabs from '$lib/components/ui/tabs';
	import { cn } from '$lib/utils';

	import {
		buildSchemaFromK8s,
		filterDataBySchema,
		formDataToK8s,
		type K8sOpenAPISchema,
		k8sToFormData,
		normalizeArrays,
		type PathOptions,
		type SchemaFormConfig
	} from './converter';
	import * as defaults from './defaults';
	import SchemaFormStep from './schema-form-step.svelte';
	import { deepMerge } from './utils';

	// ── Types ──────────────────────────────────────────────────

	export type GroupedFields = Record<string, Record<string, PathOptions>>;

	interface StepFormData {
		stepName: string;
		paths: Record<string, PathOptions>;
		formConfig: SchemaFormConfig;
		form: FormState<Record<string, unknown>>;
	}

	interface Props {
		apiSchema: K8sOpenAPISchema;
		fields: GroupedFields;
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
		fields,
		initialData,
		mode = $bindable('basic'),
		title,
		onModeChange,
		onSubmit,
		transformData,
		yamlEditable = false
	}: Props = $props();

	let currentStep = $state(0);
	let masterData = $state<Record<string, unknown>>(initialData ?? {});
	let stepForms = $state<StepFormData[]>([]);
	let advanceYaml = $state('');
	let yamlParseError = $state<string | null>(null);
	let formRefs = $state<(HTMLFormElement | undefined)[]>([]);

	// ── Derived ────────────────────────────────────────────────

	const stepNames = $derived(Object.keys(fields));
	const totalSteps = $derived(stepNames.length);
	const isFirstStep = $derived(currentStep === 0);
	const isLastStep = $derived(currentStep === totalSteps - 1);

	// ── Helpers ────────────────────────────────────────────────

	/** Merge form data into masterData after converting to K8s format */
	function mergeStepIntoMaster(
		data: Record<string, unknown>,
		mappings: SchemaFormConfig['transformationMappings']
	) {
		const k8sData = formDataToK8s(data, mappings);
		masterData = normalizeArrays(deepMerge(masterData, k8sData)) as Record<string, unknown>;
	}

	/** Collect current data from all step forms into masterData */
	function collectAllFormData() {
		for (const { form, formConfig } of stepForms) {
			const formData = getValueSnapshot(form) as Record<string, unknown>;
			mergeStepIntoMaster(formData, formConfig.transformationMappings);
		}
	}

	/** Apply transformData callback if provided */
	function applyTransform() {
		if (transformData) {
			masterData = transformData(masterData);
		}
	}

	// ── Form Lifecycle ─────────────────────────────────────────

	function createStepForms(sourceData: Record<string, unknown>) {
		stepForms = Object.entries(fields).map(([stepName, paths]) => {
			const formConfig = buildSchemaFromK8s(apiSchema, paths);
			const transformedData = k8sToFormData(sourceData, formConfig.transformationMappings);
			const initialValue = filterDataBySchema(transformedData, formConfig.schema);

			const form = createForm<Record<string, unknown>>({
				...defaults,
				idPrefix: `k8s-step-${stepName.replace(/\s+/g, '-').toLowerCase()}`,
				initialValue: initialValue,
				schema: formConfig.schema,
				uiSchema: formConfig.uiSchema,
				onSubmit: (data) => handleStepSubmit(data, formConfig)
			});

			return { stepName, paths, formConfig, form };
		});
	}

	createStepForms(masterData);

	// ── Event Handlers ─────────────────────────────────────────

	function handleStepSubmit(data: Record<string, unknown>, formConfig: SchemaFormConfig) {
		mergeStepIntoMaster(data, formConfig.transformationMappings);
		if (isLastStep) {
			handleFinalSubmit();
		} else {
			currentStep++;
		}
	}

	async function handleFinalSubmit() {
		if (mode === 'basic') collectAllFormData();
		applyTransform();
		await onSubmit?.(masterData);
	}

	function handleModeChange(newMode: string) {
		const targetMode = newMode as 'basic' | 'advance';

		if (targetMode === 'basic' && mode === 'advance') {
			syncYamlToMasterData();
			createStepForms(masterData);
		} else if (targetMode === 'advance') {
			collectAllFormData();
			applyTransform();
			syncMasterDataToYaml();
		}

		mode = targetMode;
		onModeChange?.(mode);
	}

	// ── YAML Sync ──────────────────────────────────────────────

	function syncMasterDataToYaml() {
		try {
			advanceYaml = yaml.dump(masterData, { indent: 2, lineWidth: -1 });
		} catch (error) {
			console.error('Error syncing master data to YAML:', error);
		}
	}

	function syncYamlToMasterData() {
		try {
			yamlParseError = null;
			const parsed = yaml.load(advanceYaml) as Record<string, unknown> | null;
			if (parsed && typeof parsed === 'object') masterData = parsed;
		} catch (error) {
			yamlParseError = `Invalid YAML: ${error instanceof Error ? error.message : 'Unknown error'}`;
			console.error('Error parsing YAML:', error);
		}
	}
</script>

<div class="multi-step-schema-form-container flex h-full flex-col">
	<div class="relative mb-6 flex items-center justify-center py-2">
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

	<Tabs.Root value={mode} class="flex flex-1 flex-col overflow-hidden">
		<Tabs.Content value="basic" class="flex flex-1 flex-col overflow-hidden">
			<div class="mb-6 px-6">
				<div class="flex items-center justify-between">
					{#each stepNames as stepName, index (stepName)}
						<div class="flex items-center">
							<button
								type="button"
								class={cn(
									'flex h-10 w-10 items-center justify-center rounded-full border-2 text-sm font-medium transition-all',
									index < currentStep
										? 'border-primary bg-primary text-primary-foreground'
										: index === currentStep
											? 'border-primary bg-background text-primary'
											: 'border-muted bg-muted text-muted-foreground'
								)}
								onclick={() => {
									if (index <= currentStep) {
										currentStep = index;
									}
								}}
								disabled={index > currentStep}
							>
								{#if index < currentStep}
									<Check class="h-5 w-5" />
								{:else}
									{index + 1}
								{/if}
							</button>

							<span
								class={cn(
									'ml-3 text-sm font-medium',
									index <= currentStep ? 'text-foreground' : 'text-muted-foreground'
								)}
							>
								{stepName}
							</span>
						</div>

						{#if index < totalSteps - 1}
							<div
								class={cn('mx-4 h-0.5 flex-1', index < currentStep ? 'bg-primary' : 'bg-muted')}
							></div>
						{/if}
					{/each}
				</div>
			</div>

			<div class="flex flex-1 flex-col overflow-y-auto rounded-lg p-6">
				{#each stepForms as stepForm, index (stepForm)}
					<div class={currentStep === index ? 'flex flex-1 flex-col' : 'hidden'}>
						{#key stepForm.stepName}
							<div class="multi-step-form-target contents">
								<SchemaFormStep form={stepForm.form} bind:ref={formRefs[index]} />
							</div>
						{/key}
					</div>
				{/each}
			</div>

			<div class="mt-auto flex justify-between px-6 py-4">
				<Button
					variant="outline"
					onclick={() => {
						if (!isFirstStep) {
							currentStep--;
						}
					}}
					disabled={isFirstStep}
					type="button"
				>
					← Previous
				</Button>

				<Button type="button" onclick={() => formRefs[currentStep]?.requestSubmit()}>
					{isLastStep ? 'Submit' : 'Next →'}
				</Button>
			</div>
		</Tabs.Content>

		<Tabs.Content value="advance" class="flex flex-1 flex-col overflow-hidden">
			<div class="mx-6 flex flex-1 flex-col rounded border">
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

			<div class="mt-auto px-6 py-4">
				<Button
					class="w-full"
					onclick={() => {
						syncYamlToMasterData();
						handleFinalSubmit();
					}}
				>
					Submit
				</Button>
			</div>
		</Tabs.Content>
	</Tabs.Root>
</div>

<style>
	:global(.multi-step-form-target form) {
		display: flex;
		flex-direction: column;
		flex-grow: 1;
		min-height: 100%;
	}
</style>
