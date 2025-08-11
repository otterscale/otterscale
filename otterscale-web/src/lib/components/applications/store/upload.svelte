<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/ui/button';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { cn } from '$lib/utils';

	let open = $state(false);
	function close() {
		open = false;
		step = 0;
	}
	let { class: className }: { class?: string } = $props();
	let step = $state(0);
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class={cn('flex items-center gap-1', className)}>
		<Button>
			<Icon icon="ph:upload" />
			Upload
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		{#if step == 0}
			<AlertDialog.Header>
				<AlertDialog.Title>
					Basic Information
					<div class="text-muted-foreground text-sm font-normal">
						Setup and upload your custom application template
					</div>
				</AlertDialog.Title>

				<div class="grid gap-2 space-y-2 py-4">
					<div class="grid gap-2">
						<Label>Name</Label>
						<Input value="Foo" />
					</div>
					<div class="grid gap-2">
						<Label>Description</Label>
						<Input value="Bar" />
					</div>
					<div class="grid gap-2">
						<Label>Container Image</Label>
						<Input value="docker.io/otterscale/foo:bar" />
					</div>
					<div class="grid gap-2">
						<Label>Logo URL</Label>
						<Input value="otterscale.io/logo.svg" />
					</div>
					<div class="grid grid-cols-2 gap-4">
						<div class="grid gap-2">
							<Label>Tags</Label>
							<div>
								<Badge>LLM</Badge>
								<Badge>MCP</Badge>
								<Badge>AI</Badge>
								<Badge>TIC</Badge>
							</div>
						</div>
						<div class="grid gap-2">
							<Label>Labels</Label>
							<div>
								<Badge>OtterScale</Badge>
								<Badge>Enterprise</Badge>
							</div>
						</div>
					</div>
				</div>
			</AlertDialog.Header>
		{:else if step == 1}
			<AlertDialog.Header>
				<AlertDialog.Title>
					Service
					<div class="text-muted-foreground text-sm font-normal">
						Setup and configure your service details
					</div>
				</AlertDialog.Title>

				<div class="grid gap-2 space-y-2 py-4">
					<div class="grid grid-cols-7 gap-4">
						<div class="grid gap-2">
							<Label>Enable</Label>
							<Checkbox checked={false} />
						</div>
						<div class="col-span-3 grid gap-2">
							<Label>Service Name</Label>
							<Input value="Web Page" />
						</div>
						<div class="col-span-3 grid gap-2">
							<Label>Service Port</Label>
							<Input value="8443" />
						</div>
					</div>
					<div class="grid grid-cols-7 gap-4">
						<div class="grid gap-2">
							<Label>Enable</Label>
							<Checkbox checked={false} />
						</div>
						<div class="col-span-3 grid gap-2">
							<Label>Service Name</Label>
							<Input value="Web API" />
						</div>
						<div class="col-span-3 grid gap-2">
							<Label>Service Port</Label>
							<Input value="9000" />
						</div>
					</div>
					<div class="grid gap-2">
						<Button>
							<Icon icon="ph:plus" />
						</Button>
					</div>
				</div>
			</AlertDialog.Header>
		{:else}
			<AlertDialog.Header>
				<AlertDialog.Title>
					Confirm
					<div class="text-muted-foreground text-sm font-normal">
						Setup and upload your custom application template
					</div>
				</AlertDialog.Title>

				<div>foo</div>
			</AlertDialog.Header>
		{/if}

		<AlertDialog.Footer>
			<AlertDialog.Cancel
				onclick={() => {
					close();
				}}
				class="mr-auto">Cancel</AlertDialog.Cancel
			>
			<AlertDialog.Action
				onclick={() => {
					step++;
					if (step >= 3) {
						close();
					}
				}}
			>
				{#if step == 2}
					Upload
				{:else}
					Next
				{/if}
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
