<script lang="ts">
	import { MultipleStep as MultipleStepModal } from '$lib/components/custom/modal/index';
	import { buttonVariants } from '$lib/components/ui/button';
	import { cn } from '$lib/utils';
</script>

<MultipleStepModal.Root steps={3}>
	<MultipleStepModal.Trigger
		class={cn('hover:cursor-pointer', buttonVariants({ variant: 'outline' }))}
	>
		Trigger
	</MultipleStepModal.Trigger>
	<MultipleStepModal.Content>
		<MultipleStepModal.Header>Header</MultipleStepModal.Header>
		<MultipleStepModal.Stepper>
			<MultipleStepModal.Steps>
				<MultipleStepModal.Step icon="ph:sun">
					{#snippet text({ value })}
						Step {value}
					{/snippet}
				</MultipleStepModal.Step>
				<MultipleStepModal.Step icon="ph:moon">
					{#snippet text({ value })}
						Step {value}
					{/snippet}
				</MultipleStepModal.Step>
				<MultipleStepModal.Step>
					{#snippet text({ value })}
						Step {value}
					{/snippet}
				</MultipleStepModal.Step>
			</MultipleStepModal.Steps>
			<MultipleStepModal.Models>
				<MultipleStepModal.Model>Model 1</MultipleStepModal.Model>
				<MultipleStepModal.Model>Model 2</MultipleStepModal.Model>
				<MultipleStepModal.Model>Model 3</MultipleStepModal.Model>
			</MultipleStepModal.Models>
		</MultipleStepModal.Stepper>
		<MultipleStepModal.Footer>
			<MultipleStepModal.Cancel>Cancel</MultipleStepModal.Cancel>
			<MultipleStepModal.Controllers>
				<MultipleStepModal.Back>Back</MultipleStepModal.Back>
				<MultipleStepModal.Next>Next</MultipleStepModal.Next>
				<MultipleStepModal.Confirm>Confirm</MultipleStepModal.Confirm>
			</MultipleStepModal.Controllers>
		</MultipleStepModal.Footer>
	</MultipleStepModal.Content>
</MultipleStepModal.Root>
