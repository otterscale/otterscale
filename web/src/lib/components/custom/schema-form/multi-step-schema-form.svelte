<script lang="ts">
	import { Check } from '@lucide/svelte';
	import type { FormState } from '@sjsf/form';
	import { createForm, getValueSnapshot } from '@sjsf/form';
	import { setThemeContext } from '@sjsf/shadcn4-theme';
	import * as components from '@sjsf/shadcn4-theme/new-york';
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
	import SchemaFormStep from './SchemaFormStep.svelte';
	import { deepMerge } from './utils';

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
	}

	let {
		apiSchema,
		fields,
		initialData,
		mode = $bindable('basic'),
		title,
		onModeChange,
		onSubmit
	}: Props = $props();

	setThemeContext({ components });

	let currentStep = $state(0);
	let masterData = $state<Record<string, unknown>>({});
	let stepForms = $state<StepFormData[]>([]);
	let advanceYaml = $state('');
	let yamlParseError = $state<string | null>(null);

	const stepNames = $derived(Object.keys(fields));
	const totalSteps = $derived(stepNames.length);
	const isFirstStep = $derived(currentStep === 0);
	const isLastStep = $derived(currentStep === totalSteps - 1);

	function createStepForms(sourceData: Record<string, unknown> | undefined) {
		const forms: StepFormData[] = [];
		const source = sourceData || {};

		for (const [stepName, paths] of Object.entries(fields)) {
			const formConfig = buildSchemaFromK8s(apiSchema, paths);
			// Filter initialData to only include fields defined in this step's schema
			const filteredInitialData = filterDataBySchema(source, formConfig.schema);
			const stepInitialValue = k8sToFormData(
				filteredInitialData,
				formConfig.transformationMappings
			);

			const form = createForm<Record<string, unknown>>({
				...defaults,
				idPrefix: `k8s-step-${stepName.replace(/\s+/g, '-').toLowerCase()}`,
				initialValue: stepInitialValue,
				schema: formConfig.schema,
				uiSchema: formConfig.uiSchema,
				onSubmit: (data) => handleStepSubmit(stepName, data, formConfig)
			});

			forms.push({ stepName, paths, formConfig, form });
		}

		stepForms = forms;
	}

	$effect(() => {
		if (initialData) {
			masterData = { ...initialData };
		}
		createStepForms(initialData);
	});

	function handleStepSubmit(
		stepName: string,
		data: Record<string, unknown>,
		formConfig: SchemaFormConfig
	) {
		const k8sData = formDataToK8s(data, formConfig.transformationMappings);
		masterData = normalizeArrays(deepMerge(masterData, k8sData)) as Record<string, unknown>;

		console.log(`Step "${stepName}" submitted:`, k8sData);
		console.log('Master data:', masterData);

		if (isLastStep) {
			handleFinalSubmit();
		} else {
			currentStep++;
		}
	}

	async function handleFinalSubmit() {
		collectAllFormData();
		console.log('Final submission with data:', masterData);

		if (onSubmit) {
			await onSubmit(masterData);
		}
	}

	function goBack() {
		if (!isFirstStep) {
			currentStep--;
		}
	}

	function goToStep(index: number) {
		if (index <= currentStep) {
			currentStep = index;
		}
	}

	function syncMasterDataToYaml() {
		try {
			advanceYaml = yaml.dump(masterData, {
				indent: 2,
				lineWidth: -1
			});
		} catch (error) {
			console.error('Error syncing master data to YAML:', error);
		}
	}

	function syncYamlToMasterData() {
		try {
			yamlParseError = null;
			const parsed = yaml.load(advanceYaml) as Record<string, unknown> | null;

			if (parsed && typeof parsed === 'object') {
				masterData = parsed;
			}
		} catch (error) {
			const errorMsg = `Invalid YAML: ${error instanceof Error ? error.message : 'Unknown error'}`;
			yamlParseError = errorMsg;
			console.error('Error parsing YAML:', error);
		}
	}

	function handleModeChange(newMode: string) {
		const targetMode = newMode as 'basic' | 'advance';

		if (targetMode === 'basic' && mode === 'advance') {
			syncYamlToMasterData();
			createStepForms(masterData);
		} else if (targetMode === 'advance') {
			collectAllFormData();
			syncMasterDataToYaml();
		}

		mode = targetMode;
		onModeChange?.(mode);
	}

	function collectAllFormData() {
		for (const stepForm of stepForms) {
			const formData = getValueSnapshot(stepForm.form);
			const k8sData = formDataToK8s(formData, stepForm.formConfig.transformationMappings);
			masterData = deepMerge(masterData, k8sData);
		}
		// Normalize all numeric-keyed objects to arrays
		masterData = normalizeArrays(masterData) as Record<string, unknown>;
	}

	$effect(() => {
		if (mode === 'advance') {
			syncMasterDataToYaml();
		}
	});
</script>

<div class="multi-step-schema-form-container flex h-full flex-col">
	<div class="relative mb-6 flex items-center justify-center py-2">
		{#if title}
			<h1 class="text-2xl font-bold">{title}</h1>
		{/if}
		<div class="absolute right-0">
			<Tabs.Root value={mode} onValueChange={handleModeChange}>
				<Tabs.List>
					<Tabs.Trigger value="basic">Basic</Tabs.Trigger>
					<Tabs.Trigger value="advance">Advance</Tabs.Trigger>
				</Tabs.List>
			</Tabs.Root>
		</div>
	</div>

	<Tabs.Root value={mode} class="flex flex-1 flex-col overflow-hidden">
		<Tabs.Content value="basic" class="flex flex-1 flex-col overflow-hidden">
			<div class="mb-6">
				<div class="flex items-center justify-between">
					{#each stepNames as stepName, index (stepName)}
						<div class="flex flex-1 items-center">
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
								onclick={() => goToStep(index)}
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

							{#if index < totalSteps - 1}
								<div
									class={cn('mx-4 h-0.5 flex-1', index < currentStep ? 'bg-primary' : 'bg-muted')}
								></div>
							{/if}
						</div>
					{/each}
				</div>
			</div>

			<div class="flex flex-1 flex-col overflow-y-auto rounded-lg border bg-card p-6">
				{#each stepForms as stepForm, index (stepForm)}
					<div class={currentStep === index ? 'flex flex-1 flex-col' : 'hidden'}>
						{#key stepForm.stepName}
							<div class="contents multi-step-form-target">
								<SchemaFormStep form={stepForm.form}>
									<div class="mt-auto flex justify-between pt-6">
										<Button variant="outline" onclick={goBack} disabled={isFirstStep} type="button">
											← Previous
										</Button>

										<div class="flex gap-2">
											{#if !isLastStep}
												<Button type="submit">Next →</Button>
											{:else}
												<Button type="submit">Submit</Button>
											{/if}
										</div>
									</div>
								</SchemaFormStep>
							</div>
						{/key}
					</div>
				{/each}
			</div>
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

			<Button
				class="mt-6 w-full"
				onclick={() => {
					syncYamlToMasterData();
					handleFinalSubmit();
				}}
			>
				Submit
			</Button>
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
