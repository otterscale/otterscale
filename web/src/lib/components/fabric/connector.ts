export interface Connector {
    name: string;
    icon: string;
    parameters: {
        key: string;
        name: string;
        value: string;
        description: string;
    }[];
    extraParameters?: {
        key: string;
        name: string;
        value: string;
        description: string;
    }[];
    templates?: {
        name: string;
        parameters: {
            key: string;
            value: string;
        }[];
    }[];
}