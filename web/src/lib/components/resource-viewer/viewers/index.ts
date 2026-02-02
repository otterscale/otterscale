import type { Component } from 'svelte';

type ViewerProps = { object: any; schema?: any };
type ViewerType = Component<ViewerProps>;

// Edit Form Types
type EditFormProps = { name: string; schema?: any; object?: any; onsuccess?: () => void };
type EditFormType = Component<EditFormProps> | null;

// Delete Form Types
type DeleteFormProps = { name: string };
type DeleteFormType = Component<DeleteFormProps> | null;

// Viewer Components
// Delete Form Components
import DeleteDialog from '$lib/components/form/workspace/delete-dialog.svelte';
// Edit Form Components
import EditSheet from '$lib/components/form/workspace/edit-sheet.svelte';

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
		return EditSheet as unknown as EditFormType;
	}
	return null;
}

function getDeleteForm(resource: string): DeleteFormType {
	if (resource === 'workspaces') {
		return DeleteDialog as unknown as DeleteFormType;
	}
	return null;
}

export { getDeleteForm, getEditForm, getResourceViewer };
export type { DeleteFormType, EditFormType, ViewerType };
