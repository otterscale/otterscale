type Message = {
	senderId: string;
	message: string;
	sentAt: string;
};

type Choice = {
	index: number;
	finish_reason: string;
	text: string;
};

export type { Choice, Message };
