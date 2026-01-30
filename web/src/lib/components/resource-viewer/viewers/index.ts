import type { Component } from 'svelte';

type ViewerProps = { object: any; schema?: any };
type ViewerType = Component<ViewerProps>;

// Edit Form Types
type EditFormProps = { name: string; schema?: any; object?: any };
type EditFormType = Component<EditFormProps> | null;

// Delete Form Types
type DeleteFormProps = { name: string };
type DeleteFormType = Component<DeleteFormProps> | null;

// Viewer Components
// Delete Form Components
import DeleteWorkspaceDialog from '$lib/components/form/workspace/delete-workspace-dialog.svelte';
// Edit Form Components
import SheetEditWorkspace from '$lib/components/form/workspace/sheet-edit-workspace.svelte';

import Default from './default.svelte';
import Workspaces from './workspaces.svelte';

function getResourceViewer(resource: string): ViewerType {
	if (resource === 'workspaces') {
		return Workspaces as ViewerType;
	}
	return Default as ViewerType;
}

function getEditForm(resource: string): EditFormType {
	if (resource === 'workspaces') {
		return SheetEditWorkspace as unknown as EditFormType;
	}
	return null;
}

function getDeleteForm(resource: string): DeleteFormType {
	if (resource === 'workspaces') {
		return DeleteWorkspaceDialog as unknown as DeleteFormType;
	}
	return null;
}

export type { DeleteFormType, EditFormType, ViewerType };
export { getDeleteForm, getEditForm, getResourceViewer };
