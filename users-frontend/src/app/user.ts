export type User = {
    user_id?: number;
    user_name: string; // 50 chars
    first_name: string;
    last_name: string;
    email: string;
    user_status: UserStatus;
    department?: string;
}

export enum UserStatus {
    Active = "A",
    Inactive = "I",
    Terminated = "T"
}

export function userStatusToString(user_status: UserStatus): string {
    if(user_status === UserStatus.Active) {
        return "Active"
    } 
    else if (user_status === UserStatus.Inactive) {
        return "Inactive"
    }
    else if (user_status === UserStatus.Terminated) {
        return "Terminated"
    }
    return ""
}