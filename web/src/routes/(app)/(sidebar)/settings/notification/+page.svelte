<script lang="ts" context="module">
	import { z } from 'zod';
	export const notificationsFormSchema = z.object({
		type: z.enum(['all', 'mentions', 'none'], {
			required_error: 'You need to select a notification type.'
		}),
		mobile: z.boolean().default(false).optional(),
		communication_emails: z.boolean().default(false).optional(),
		social_emails: z.boolean().default(false).optional(),
		marketing_emails: z.boolean().default(false).optional(),
		security_emails: z.boolean()
	});
	type NotificationFormSchema = typeof notificationsFormSchema;
</script>

<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { tick } from 'svelte';
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';

	import SuperDebug, { type Infer, type SuperValidated, superForm } from 'sveltekit-superforms';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import * as Form from '$lib/components/ui/form/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { browser } from '$app/environment';
	import { toast } from 'svelte-sonner';
	import Icon from '@iconify/svelte';

	type Item = {
		value: string;
		label: string;
	};

	const items: Item[] = [
		{
			value: '0',
			label: 'All new messages'
		},
		{
			value: '1',
			label: 'Direct messages and mentions'
		},
		{
			value: '2',
			label: 'Nothing'
		}
	];

	let open: boolean = false;
	let value = '';

	$: selectedItem = items.find((s) => s.value === value) ?? null;

	// We want to refocus the trigger button when the user selects
	// an item from the list so users can continue navigating the
	// rest of the form with the keyboard.
	function closeAndFocusTrigger(triggerId: string) {
		open = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}

	export let data: SuperValidated<Infer<NotificationFormSchema>>;

	const form = superForm(data, {
		validators: zodClient(notificationsFormSchema)
	});

	const { form: formData, enhance } = form;
</script>

<Card.Root>
	<Card.Header class="pb-3">
		<Card.Title>Notifications</Card.Title>
		<Card.Description>Choose what you want to be notified about.</Card.Description>
	</Card.Header>
	<Card.Content class="grid gap-1">
		<div
			class="-mx-2 flex items-start space-x-4 rounded-md p-2 transition-all hover:bg-accent hover:text-accent-foreground"
		>
			<!-- <Bell class="mt-px h-5 w-5" /> -->
			<Icon icon="line-md:bell-loop" class="mt-px h-5 w-5" />
			<div class="space-y-1">
				<p class="text-sm font-medium leading-none">Everything</p>
				<p class="text-sm text-muted-foreground">Email digest, mentions & all activity.</p>
			</div>
		</div>
		<div
			class="-mx-2 flex items-start space-x-4 rounded-md bg-accent p-2 text-accent-foreground transition-all"
		>
			<Icon icon="line-md:account" class="mt-px h-5 w-5" />
			<div class="space-y-1">
				<p class="text-sm font-medium leading-none">Available</p>
				<p class="text-sm text-muted-foreground">Only mentions and comments.</p>
			</div>
		</div>
		<div
			class="-mx-2 flex items-start space-x-4 rounded-md p-2 transition-all hover:bg-accent hover:text-accent-foreground"
		>
			<Icon icon="line-md:phone-off-loop" class="mt-px h-5 w-5" />
			<div class="space-y-1">
				<p class="text-sm font-medium leading-none">Ignoring</p>
				<p class="text-sm text-muted-foreground">Turn off all notifications.</p>
			</div>
		</div>
	</Card.Content>
</Card.Root>

<Card.Root>
	<Card.Header>
		<Card.Title>Email</Card.Title>
		<Card.Description>You can customize which events trigger notifications.</Card.Description>
	</Card.Header>
	<Card.Content>
		<div class="flex-col items-center space-y-6">
			<div class="flex items-center space-x-4">
				<p class="text-sm text-muted-foreground">Notify me when</p>
				<!-- <Popover.Root bind:open let:ids>
					<Popover.Trigger asChild let:builder>
						<Button builders={[builder]} variant="outline" class="w-[288px] justify-start">
							{selectedItem ? selectedItem.label : items[0].label}
						</Button>
					</Popover.Trigger>
					<Popover.Content class="p-0" align="start">
						<Command.Root>
							<Command.Input placeholder="Change status..." />
							<Command.List>
								<Command.Empty>No results found.</Command.Empty>
								<Command.Group>
									{#each items as item}
										<Command.Item
											value={item.value}
											onSelect={() => {
												closeAndFocusTrigger(ids.trigger);
											}}
										>
											{item.label}
										</Command.Item>
									{/each}
								</Command.Group>
							</Command.List>
						</Command.Root>
					</Popover.Content>
				</Popover.Root> -->
			</div>
			<form method="POST" use:enhance class="space-y-6">
				<Form.Field
					{form}
					name="communication_emails"
					class="flex flex-row items-center justify-between rounded-lg border p-4"
				>
					<!-- <Form.Control let:attrs>
						<div class="space-y-0.5">
							<Form.Label class="text-base">Communication emails</Form.Label>
							<Form.Description>Receive emails about your account activity.</Form.Description>
						</div>
						<Switch {...attrs} bind:checked={$formData.communication_emails} />
					</Form.Control> -->
				</Form.Field>
				<Form.Field
					{form}
					name="marketing_emails"
					class="flex flex-row items-center justify-between rounded-lg border p-4"
				>
					<!-- <Form.Control let:attrs>
						<div class="space-y-0.5">
							<Form.Label class="text-base">Marketing emails</Form.Label>
							<Form.Description>
								Receive emails about new products, features, and more.
							</Form.Description>
						</div>
						<Switch {...attrs} bind:checked={$formData.marketing_emails} />
					</Form.Control> -->
				</Form.Field>
				<Form.Field
					{form}
					name="social_emails"
					class="flex flex-row items-center justify-between rounded-lg border p-4"
				>
					<!-- <Form.Control let:attrs>
						<div class="space-y-0.5">
							<Form.Label class="text-base">Social emails</Form.Label>
							<Form.Description>
								Receive emails for friend requests, follows, and more.
							</Form.Description>
						</div>
						<Switch {...attrs} bind:checked={$formData.social_emails} />
					</Form.Control> -->
				</Form.Field>
				<Form.Field
					{form}
					name="security_emails"
					class="flex flex-row items-center justify-between rounded-lg border p-4"
				>
					<!-- <Form.Control let:attrs>
						<div class="space-y-0.5">
							<Form.Label class="text-base">Security emails</Form.Label>
							<Form.Description>
								Receive emails about your account activity and security.
							</Form.Description>
						</div>
						<Switch {...attrs} bind:checked={$formData.security_emails} />
					</Form.Control> -->
				</Form.Field>
			</form>
			<div></div>
		</div></Card.Content
	>
	<Card.Footer class="border-t px-6 py-4">
		<Button onclick={() => toast.success('Saved!')}>Save</Button>
	</Card.Footer>
</Card.Root>
