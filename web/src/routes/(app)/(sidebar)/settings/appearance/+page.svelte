<script lang="ts" context="module">
	import { z } from 'zod';

	export const appearanceFormSchema = z.object({
		theme: z.enum(['light', 'dark'], {
			required_error: 'Please select a theme.'
		}),
		font: z.enum(['inter', 'manrope', 'system'], {
			invalid_type_error: 'Select a font',
			required_error: 'Please select a font.'
		})
	});

	export type AppearanceFormSchema = typeof appearanceFormSchema;
</script>

<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	// import ChevronDown from 'svelte-radix/ChevronDown.svelte';
	import SuperDebug, { type Infer, type SuperValidated, superForm } from 'sveltekit-superforms';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import * as Form from '$lib/components/ui/form/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { cn } from '$lib/utils.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { toast } from 'svelte-sonner';
	import { mode } from 'mode-watcher';
	export let data: SuperValidated<Infer<AppearanceFormSchema>>;

	const form = superForm(data, {
		validators: zodClient(appearanceFormSchema)
	});

	const { form: formData, enhance } = form;
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Font</Card.Title>
		<Card.Description>Set the font you want to use in the dashboard.</Card.Description>
	</Card.Header>
	<Card.Content>
		<Form.Field {form} name="font">
			<!-- <Form.Control let:attrs>
				<div class="relative w-max">
					<select
						{...attrs}
						class={cn(
							buttonVariants({ variant: 'outline' }),
							'w-[200px] appearance-none font-normal'
						)}
						bind:value={$formData.font}
					>
						<option value="inter">Inter</option>
						<option value="manrope">Manrope</option>
						<option value="system">System</option>
					</select>
				</div>
			</Form.Control> -->
		</Form.Field>
	</Card.Content>
	<Card.Footer class="border-t px-6 py-4">
		<Button onclick={() => toast.success('Saved!')}>Save</Button>
	</Card.Footer>
</Card.Root>

<Card.Root>
	<Card.Header>
		<Card.Title>Theme</Card.Title>
		<Card.Description>Select the theme for the dashboard.</Card.Description>
	</Card.Header>
	<Card.Content>
		<Form.Fieldset {form} name="theme">
			<RadioGroup.Root
				class="grid max-w-md grid-cols-2 gap-8 pt-2"
				orientation="horizontal"
				bind:value={$formData.theme}
			>
				<Form.Control let:attrs>
					<Label class="[&:has([data-state=checked])>div]:border-primary">
						<RadioGroup.Item {...attrs} value="light" class="sr-only" />
						<div class="items-center rounded-md border-2 border-muted p-1 hover:border-accent">
							<div class="space-y-2 rounded-sm bg-[#ecedef] p-2">
								<div class="space-y-2 rounded-md bg-white p-2 shadow-sm">
									<div class="h-2 w-[80px] rounded-lg bg-[#ecedef]"></div>
									<div class="h-2 w-[100px] rounded-lg bg-[#ecedef]"></div>
								</div>
								<div class="flex items-center space-x-2 rounded-md bg-white p-2 shadow-sm">
									<div class="h-4 w-4 rounded-full bg-[#ecedef]"></div>
									<div class="h-2 w-[100px] rounded-lg bg-[#ecedef]"></div>
								</div>
								<div class="flex items-center space-x-2 rounded-md bg-white p-2 shadow-sm">
									<div class="h-4 w-4 rounded-full bg-[#ecedef]"></div>
									<div class="h-2 w-[100px] rounded-lg bg-[#ecedef]"></div>
								</div>
							</div>
						</div>
						<span class="block w-full p-2 text-center font-normal"> Light </span>
					</Label>
				</Form.Control>
				<Form.Control let:attrs>
					<Label class="[&:has([data-state=checked])>div]:border-primary">
						<RadioGroup.Item {...attrs} value="dark" class="sr-only" />
						<div
							class="items-center rounded-md border-2 border-muted bg-popover p-1 hover:bg-accent hover:text-accent-foreground"
						>
							<div class="space-y-2 rounded-sm bg-slate-950 p-2">
								<div class="space-y-2 rounded-md bg-slate-800 p-2 shadow-sm">
									<div class="h-2 w-[80px] rounded-lg bg-slate-400"></div>
									<div class="h-2 w-[100px] rounded-lg bg-slate-400"></div>
								</div>
								<div class="flex items-center space-x-2 rounded-md bg-slate-800 p-2 shadow-sm">
									<div class="h-4 w-4 rounded-full bg-slate-400"></div>
									<div class="h-2 w-[100px] rounded-lg bg-slate-400"></div>
								</div>
								<div class="flex items-center space-x-2 rounded-md bg-slate-800 p-2 shadow-sm">
									<div class="h-4 w-4 rounded-full bg-slate-400"></div>
									<div class="h-2 w-[100px] rounded-lg bg-slate-400"></div>
								</div>
							</div>
						</div>
						<span class="block w-full p-2 text-center font-normal"> Dark </span>
					</Label>
				</Form.Control>
				<!-- <RadioGroup.Input name="theme" /> -->
			</RadioGroup.Root>
		</Form.Fieldset>
	</Card.Content>
	<Card.Footer class="border-t px-6 py-4">
		<Button onclick={() => toast.success('Saved!')}>Save</Button>
	</Card.Footer>
</Card.Root>
