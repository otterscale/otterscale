interface UploadedFile {
	name: string;
	size: number;
	type: string;
	url: Promise<string>;
	uploadedAt: number;
	lastModifiedAt: number;
}

export type { UploadedFile };
