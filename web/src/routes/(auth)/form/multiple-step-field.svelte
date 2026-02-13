<script lang="ts" module>
	import { createContext } from 'svelte';

	export const [getStepperContext, setStepperContext] = createContext<Ref<number>>();
</script>

<script lang="ts">
	import {
		type ComponentProps,
		type FieldValue,
		getChildPath,
		getFieldComponent,
		getFormContext,
		getValueSnapshot,
		retrieveTranslate,
		retrieveUiOption,
		retrieveUiSchema,
		uiTitleOption,
		updateErrors,
		validate
	} from '@sjsf/form';
	import { isSchemaObject } from '@sjsf/form/lib/json-schema';
	import type { Ref } from '@sjsf/form/lib/svelte.svelte';
	import { mode as themeMode } from 'mode-watcher';
	import Monaco from 'svelte-monaco';
	import { stringify } from 'yaml';

	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import Button, { buttonVariants } from '$lib/components/ui/button/button.svelte';
	import * as Item from '$lib/components/ui/item';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import Switch from '$lib/components/ui/switch/switch.svelte';

	let { config, value = $bindable() }: ComponentProps['tupleField'] = $props();

	const ctx = getFormContext();
	const stepperCtx = getStepperContext();

	const stepSchemas = $derived.by(() => {
		const items = config.schema.items;
		return Array.isArray(items) && items.every(isSchemaObject) ? items : [];
	});
	const stepUiSchemas = $derived.by(() => {
		const items = config.uiSchema.items ?? {};
		return (Array.isArray(items) ? items : stepSchemas.map(() => items)).map((s) => {
			const retrieved = retrieveUiSchema(ctx, s);
			return {
				...retrieved,
				'ui:options': {
					...retrieved['ui:options'],
					hideTitle: true
				}
			};
		});
	});
	const stepTitles = $derived(
		stepUiSchemas.map((s, i) => uiTitleOption(ctx, s) ?? stepSchemas[i].title ?? `Step ${i + 1}`)
	);

	const stepConfig = $derived({
		path: getChildPath(ctx, config.path, stepperCtx.current),
		required: true,
		schema: stepSchemas[stepperCtx.current],
		uiSchema: stepUiSchemas[stepperCtx.current],
		title: stepTitles[stepperCtx.current]
	});

	const Form = $derived(getFieldComponent(ctx, stepConfig));

	let open = $state(false);

	let isBasicMode = $state(true);
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class={buttonVariants({ variant: 'outline' })}>
		Multiple Step Form
	</AlertDialog.Trigger>
	<AlertDialog.Content class="flex max-h-[90vh] min-h-[90vh] min-w-[50vw] flex-col gap-0 p-4">
		{#if isBasicMode}
			<AlertDialog.Header>
				<Item.Root>
					<Item.Content class="h-15 text-left">
						<Item.Title class="text-lg font-bold"
							>{stepSchemas[stepperCtx.current].title}</Item.Title
						>
						<Item.Description>{stepSchemas[stepperCtx.current].description}</Item.Description>
					</Item.Content>
					<Item.Actions>
						<Switch bind:checked={isBasicMode} />
					</Item.Actions>
					<Progress value={(stepperCtx.current + 1) / stepSchemas.length} max={1} />
				</Item.Root>
			</AlertDialog.Header>
			<!-- Form -->
			<div class="h-full overflow-y-auto p-4">
				<Form
					config={stepConfig}
					translate={retrieveTranslate(ctx, stepConfig)}
					type="field"
					uiOption={(opt) => retrieveUiOption(ctx, stepConfig, opt)}
					bind:value={
						() => value?.[stepperCtx.current] as undefined,
						(v) => {
							if (value) {
								value[stepperCtx.current] = v;
							} else {
								const arr = new Array<FieldValue>(stepSchemas.length);
								arr[stepperCtx.current] = v;
								value = arr;
							}
						}
					}
				/>
			</div>
			<!-- Actions -->
			<AlertDialog.Footer class="mt-auto flex items-center justify-between p-4">
				<Button
					size="sm"
					class="mr-auto"
					type="button"
					onclick={() => {
						stepperCtx.current--;
					}}
				>
					Back
				</Button>
				{#if stepperCtx.current < stepSchemas.length - 1}
					<Button
						size="sm"
						type="button"
						onclick={() => {
							const { errors } = validate(ctx);
							const s = stepperCtx.current;
							const currentErrors = errors?.filter((e) => e.path[0] === s);
							if (currentErrors?.length) {
								updateErrors(ctx, currentErrors);
							} else {
								stepperCtx.current++;
							}
						}}
					>
						Continue
					</Button>
				{/if}
				{#if stepperCtx.current === stepSchemas.length - 1}
					<Button
						size="sm"
						type="submit"
						onclick={() => {
							open = false;
						}}
					>
						Submit
					</Button>
				{/if}
			</AlertDialog.Footer>
		{:else}
			<AlertDialog.Header>
				<Item.Root>
					<Item.Content class="h-15 text-left">
						<Item.Title class="text-lg font-bold">Dynamic Form</Item.Title>
						<Item.Description>Description for dynamic form.</Item.Description>
					</Item.Content>
					<Item.Actions>
						<Switch bind:checked={isBasicMode} />
					</Item.Actions>
				</Item.Root>
			</AlertDialog.Header>
			<div class="mt-auto h-[80vh]">
				<Monaco
					options={{
						language: 'yaml',
						padding: { top: 16, bottom: 8 },
						automaticLayout: true,
						minimap: { enabled: false },
						scrollBeyondLastLine: false,
						readOnly: !isBasicMode
					}}
					theme={themeMode.current === 'dark' ? 'vs-dark' : 'vs'}
					value={JSON.stringify(value, null, 2)}
				/>
			</div>
		{/if}
	</AlertDialog.Content>
</AlertDialog.Root>
