type TableRow = {
    id: number;
    name: string;
    email: string;
    role: 'Admin' | 'Manager' | 'User';
    status: 'active' | 'inactive';
    createdAt: string;
    isVerified: boolean;
    loginCount: number;
    rating: string;
    lastLoginDays: number;
};

export type { TableRow }
