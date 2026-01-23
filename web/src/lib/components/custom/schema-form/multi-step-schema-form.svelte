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
		formDataToK8s,
		type K8sOpenAPISchema,
		k8sToFormData,
		type PathOptions,
		type SchemaFormConfig
	} from './converter';
	import SchemaFormStep from './SchemaFormStep.svelte';
	import * as defaults from './defaults';

	/** Grouped fields: Record<StepName, Record<Path, PathOptions>> */
	export type GroupedFields = Record<string, Record<string, PathOptions>>;

	interface StepFormData {
		stepName: string;
		paths: Record<string, PathOptions>;
		formConfig: SchemaFormConfig;
		form: FormState<Record<string, unknown>>;
	}

	interface Props {
		/** The full K8s OpenAPI V3 Schema */
		apiSchema: K8sOpenAPISchema;
		/** Grouped paths by step: { 'Step Name': { 'path.to.field': { title: 'Field' } } } */
		fields: GroupedFields;
		/** Optional initial value override */
		initialData?: Record<string, unknown>;
		/** Current mode: 'basic' | 'advance' */
		mode?: 'basic' | 'advance';
		/** Callback when mode changes */
		onModeChange?: (mode: 'basic' | 'advance') => void;
		/** Callback when final submit is clicked with accumulated data */
		onSubmit?: (data: Record<string, unknown>) => Promise<void> | void;
	}

	let {
		apiSchema,
		fields,
		initialData,
		mode = $bindable('basic'),
		onModeChange,
		onSubmit
	}: Props = $props();

	// Set theme context for this component tree
	setThemeContext({ components });

	// Step state management
	let currentStep = $state(0);
	let masterData = $state<Record<string, unknown>>({});
	let stepForms = $state<StepFormData[]>([]);
	let advanceYaml = $state('');
	let yamlParseError = $state<string | null>(null);

	// Derived values
	const stepNames = $derived(Object.keys(fields));
	const totalSteps = $derived(stepNames.length);
	const isFirstStep = $derived(currentStep === 0);
	const isLastStep = $derived(currentStep === totalSteps - 1);

	// Build forms for each step
	$effect(() => {
		const forms: StepFormData[] = [];

		for (const [stepName, paths] of Object.entries(fields)) {
			const formConfig = buildSchemaFromK8s(apiSchema, paths);

			// Convert initial data using the step's transformation mappings
			const stepInitialValue = k8sToFormData(initialData, formConfig.transformationMappings);

			const form = createForm<Record<string, unknown>>({
				...defaults,
				idPrefix: `k8s-step-${stepName.replace(/\s+/g, '-').toLowerCase()}`,
				initialValue: stepInitialValue,
				schema: formConfig.schema,
				uiSchema: formConfig.uiSchema,
				onSubmit: (data) => handleStepSubmit(stepName, data, formConfig)
			});

			forms.push({
				stepName,
				paths,
				formConfig,
				form
			});
		}

		stepForms = forms;

		// Initialize masterData with initial data
		if (initialData) {
			masterData = { ...initialData };
		}
	});

	// Handle individual step submission
	function handleStepSubmit(
		stepName: string,
		data: Record<string, unknown>,
		formConfig: SchemaFormConfig
	) {
		// Convert form data to K8s format
		const k8sData = formDataToK8s(data, formConfig.transformationMappings);

		// Deep merge into master data
		masterData = deepMerge(masterData, k8sData);

		console.log(`Step "${stepName}" submitted:`, k8sData);
		console.log('Master data:', masterData);

		// Move to next step or finish
		if (isLastStep) {
			handleFinalSubmit();
		} else {
			currentStep++;
		}
	}

	// Handle final submission
	async function handleFinalSubmit() {
		console.log('Final submission with data:', masterData);

		if (onSubmit) {
			await onSubmit(masterData);
		}
	}

	// Navigate to previous step
	function goBack() {
		if (!isFirstStep) {
			currentStep--;
		}
	}

	// Navigate to specific step (only if already visited)
	function goToStep(index: number) {
		if (index <= currentStep) {
			currentStep = index;
		}
	}



	// Deep merge utility
	function deepMerge(
		target: Record<string, unknown>,
		source: Record<string, unknown>
	): Record<string, unknown> {
		const result = { ...target };

		for (const key of Object.keys(source)) {
			const sourceValue = source[key];
			const targetValue = result[key];

			if (
				sourceValue &&
				typeof sourceValue === 'object' &&
				!Array.isArray(sourceValue) &&
				targetValue &&
				typeof targetValue === 'object' &&
				!Array.isArray(targetValue)
			) {
				result[key] = deepMerge(
					targetValue as Record<string, unknown>,
					sourceValue as Record<string, unknown>
				);
			} else {
				result[key] = sourceValue;
			}
		}

		return result;
	}

	// Sync master data to YAML editor
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

	// Sync YAML back to master data
	function syncYamlToMasterData() {
		try {
			yamlParseError = null;
			const parsed = yaml.load(advanceYaml) as Record<string, unknown> | null;

			if (parsed && typeof parsed === 'object') {
				masterData = parsed;

				// Note: Form values will be rebuilt on next step navigation
				// The masterData is the source of truth after YAML edit
			}
		} catch (error) {
			const errorMsg = `Invalid YAML: ${error instanceof Error ? error.message : 'Unknown error'}`;
			yamlParseError = errorMsg;
			console.error('Error parsing YAML:', error);
		}
	}

	// Handle mode changes
	function handleModeChange(newMode: string) {
		const targetMode = newMode as 'basic' | 'advance';

		if (targetMode === 'basic' && mode === 'advance') {
			syncYamlToMasterData();
		} else if (targetMode === 'advance') {
			// Collect all current form values into master data before showing YAML
			collectAllFormData();
			syncMasterDataToYaml();
		}

		mode = targetMode;
		onModeChange?.(mode);
	}

	// Collect data from all forms
	function collectAllFormData() {
		for (const stepForm of stepForms) {
			const formData = getValueSnapshot(stepForm.form);
			const k8sData = formDataToK8s(formData, stepForm.formConfig.transformationMappings);
			masterData = deepMerge(masterData, k8sData);
		}
	}

	// Reactive effect to sync YAML when in advance mode
	$effect(() => {
		if (mode === 'advance') {
			syncMasterDataToYaml();
		}
	});
</script>

<div class="multi-step-schema-form-container">
	<!-- Mode Toggle -->
	<div class="mb-4 flex items-center justify-end">
		<Tabs.Root value={mode} onValueChange={handleModeChange}>
			<Tabs.List>
				<Tabs.Trigger value="basic">Basic</Tabs.Trigger>
				<Tabs.Trigger value="advance">Advance</Tabs.Trigger>
			</Tabs.List>
		</Tabs.Root>
	</div>

	<Tabs.Root value={mode}>
		<!-- Basic Mode: Multi-Step Form -->
		<Tabs.Content value="basic">
			<!-- Step Indicators -->
			<div class="mb-8">
				<div class="flex items-center justify-between">
					{#each stepNames as stepName, index (stepName)}
						<div class="flex flex-1 items-center">
							<!-- Step Circle -->
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

							<!-- Step Label -->
							<span
								class={cn(
									'ml-3 text-sm font-medium',
									index <= currentStep ? 'text-foreground' : 'text-muted-foreground'
								)}
							>
								{stepName}
							</span>

							<!-- Connector Line -->
							{#if index < totalSteps - 1}
								<div
									class={cn('mx-4 h-0.5 flex-1', index < currentStep ? 'bg-primary' : 'bg-muted')}
								></div>
							{/if}
						</div>
					{/each}
				</div>
			</div>

			<!-- Step Content -->
			<div class="rounded-lg border bg-card p-6">
				{#each stepForms as stepForm, index (stepForm.stepName)}
					<div class={currentStep === index ? 'block' : 'hidden'}>
						{#key stepForm.stepName}
							<SchemaFormStep form={stepForm.form}>
								<div class="mt-6 flex justify-between">
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
						{/key}
					</div>
				{/each}
			</div>
		</Tabs.Content>

		<!-- Advance Mode: YAML Editor -->
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
