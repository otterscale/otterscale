type AlertType = {
	level: string;
	title: string;
	message: string;
	action?: () => void;
};

export type { AlertType };
