<script lang="ts">
	import { type Essential } from '$gen/api/essential/v1/essential_pb';
	import { ManagementApplications } from '$lib/components/otterscale/index';
	import * as Select from '$lib/components/ui/select/index.js';

	let { kuberneteses }: { kuberneteses: Essential[] } = $props();

	function getKey(kubernetes: Essential) {
		return kubernetes.scopeUuid + kubernetes.name;
	}
	function getIdentifier(kubernetes: Essential) {
		return [kubernetes.scopeName, kubernetes.name].join('/');
	}

	let defaultKubernetes = $state(kuberneteses[0] as Essential);
</script>

<ManagementApplications />
<!-- 
<div class="grid gap-2">
	<span class="flex items-center gap-2 p-4">
		Kubernetes
		<Select.Root type="single">
			<Select.Trigger class="border-0 border-b bg-muted focus:ring-0 focus:ring-offset-0">
				{getIdentifier(defaultKubernetes)}
			</Select.Trigger>
			<Select.Content>
				{#each kuberneteses as kubernetes}
					{@const selection = getKey(kubernetes)}
					<Select.Item
						value={selection}
						onclick={() => {
							defaultKubernetes = kubernetes;
						}}
					>
						{getIdentifier(kubernetes)}
					</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</span>
	{#key defaultKubernetes}
		<ManagementApplications
			scopeUuid={defaultKubernetes.scopeUuid}
			facilityName={defaultKubernetes.facilityName}
		/>
	{/key}
</div> -->
