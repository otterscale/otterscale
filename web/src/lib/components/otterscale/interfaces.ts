export interface Parameter {
    key: string;
    name: string;
    value?: string;
    values?: string[];
    description: string;
}

export interface Configuration {
    key: string;
    name: string;
    icon: string;
    steps: {
        description: string,
        parameters: Parameter[];
        advancedParameters?: Parameter[];
    }[];
    templates?: {
        name: string;
        parameters: {
            key: string;
            value: string;
        }[];
    }[];
}

export interface Instance {
    id: string;
    type: string;
    name: string;
    tags: string[]
    neighbors: string[];
    information: {},
    created: Date;
    updated: Date;
}