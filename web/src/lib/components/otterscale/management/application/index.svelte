<script lang="ts">
	import { type Facility_Info } from '$gen/api/nexus/v1/nexus_pb';
	import { ManagementApplications } from '$lib/components/otterscale/index';
	import * as Select from '$lib/components/ui/select/index.js';

	let { kuberneteses }: { kuberneteses: Facility_Info[] } = $props();

	function getKey(kubernetes: Facility_Info) {
		return kubernetes.scopeUuid + kubernetes.facilityName;
	}
	function getIdentifier(kubernetes: Facility_Info) {
		return [kubernetes.scopeName, kubernetes.facilityName].join('/');
	}

	let defaultKubernetes = $state(kuberneteses[0] as Facility_Info);
</script>

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
</div>
