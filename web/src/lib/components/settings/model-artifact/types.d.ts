type SortType = 'downloads' | 'likes';

type ModelTagCategory = 'region' | 'language' | 'license' | 'library' | 'pipeline_tag';

type ModelTag = {
	id: string;
	type: string;
	subType: string;
	label: string;
};

type HuggingFaceModel = {
	id: string;
	tags: string[];
	downloads: number;
	likes: number;
	createdAt: string;
};

export type { HuggingFaceModel, ModelTag, ModelTagCategory, SortType };
