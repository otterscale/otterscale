<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import * as Tabs from '$lib/components/ui/tabs';
	import { type Scope } from '$gen/api/nexus/v1/nexus_pb';
	import { ManagementScopes, ManagementFacilities } from '$lib/components/otterscale/index';

	let { scopes }: { scopes: Scope[] } = $props();

	let scopeByName = $derived(Object.fromEntries(scopes.map((scope) => [scope.name, scope])));
	let chosenScopeName = $state('controller');
</script>

<main class="p-4">
	<Tabs.Root value="facility" class="w-full">
		<Tabs.List>
			<Tabs.Trigger value="facility">FACILITY</Tabs.Trigger>
			<Tabs.Trigger value="scope">SCOPE</Tabs.Trigger>
		</Tabs.List>
		<Tabs.Content value="facility">
			<div class="grid gap-2">
				{@render ChooseScope()}
				{#key chosenScopeName}
					<ManagementFacilities scopeUuid={scopeByName[chosenScopeName].uuid} />
				{/key}
				<div></div>
			</div>
		</Tabs.Content>
		<Tabs.Content value="scope">
			<div class="grid gap-2">
				<ManagementScopes {scopes} />
			</div>
		</Tabs.Content>
	</Tabs.Root>
</main>

{#snippet ChooseScope()}
	<span class="flex items-center gap-1 text-xl">
		Scope
		<Select.Root type="single">
			<Select.Trigger
				class="rounded-background border-0 border-b bg-muted focus:ring-0 focus:ring-offset-0"
			>
				<h1 class=" text-xl">{chosenScopeName}</h1>
			</Select.Trigger>
			<Select.Content>
				{#each scopes as model}
					<Select.Item
						value={model.name}
						onclick={() => {
							chosenScopeName = model.name;
						}}
					>
						{model.name}
					</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</span>
{/snippet}
