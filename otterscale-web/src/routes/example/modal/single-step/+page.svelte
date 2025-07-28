<script lang="ts">
	import { SingleStep as SingleStepModal } from '$lib/components/custom/modal';
	import Icon from '@iconify/svelte';
</script>

<SingleStepModal.Root>
	<SingleStepModal.Trigger variant="default">
		<Icon icon="ph:mouse" />
		Trigger
	</SingleStepModal.Trigger>
	<SingleStepModal.Trigger variant="creative">
		<Icon icon="ph:mouse" />
		Trigger
	</SingleStepModal.Trigger>
	<SingleStepModal.Trigger variant="destructive">
		<Icon icon="ph:mouse" />
		Trigger
	</SingleStepModal.Trigger>
	<SingleStepModal.Content>
		<SingleStepModal.Header>Header</SingleStepModal.Header>
		<SingleStepModal.Footer>
			<SingleStepModal.Cancel>Cancel</SingleStepModal.Cancel>
			<SingleStepModal.ActionsGroup>
				<SingleStepModal.Actions>
					<SingleStepModal.ActionGroup>
						<SingleStepModal.ActionGroupHeading>Group</SingleStepModal.ActionGroupHeading>
						<SingleStepModal.ActionItem>
							Debug
							<SingleStepModal.ActionShortcut>⌘S</SingleStepModal.ActionShortcut>
						</SingleStepModal.ActionItem>
					</SingleStepModal.ActionGroup>
					<SingleStepModal.ActionGroup>
						<SingleStepModal.ActionGroupHeading>Group</SingleStepModal.ActionGroupHeading>
						<SingleStepModal.ActionItem>
							Action
							<SingleStepModal.ActionShortcut>⌘S</SingleStepModal.ActionShortcut>
						</SingleStepModal.ActionItem>
						<SingleStepModal.ActionItem>
							Action
							<SingleStepModal.ActionShortcut>⌘S</SingleStepModal.ActionShortcut>
						</SingleStepModal.ActionItem>
					</SingleStepModal.ActionGroup>
				</SingleStepModal.Actions>
				<SingleStepModal.Action>Confirm</SingleStepModal.Action>
			</SingleStepModal.ActionsGroup>
		</SingleStepModal.Footer>
	</SingleStepModal.Content>
</SingleStepModal.Root>
