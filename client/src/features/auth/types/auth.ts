import { type User } from "../../../shared/types/user";

export interface AuthData {
    message: string;
    token: string;
    user: Partial<User>;
}

export interface AuthError {
    error: string;
    details: string;
}

export interface RegisterData {
    email: string;
    first_name: string;
    last_name: string;
    age: number;
    gender: "male" | "female" | "non_binary" | "prefer_not_to_say";
    password: string;
}


export interface LoginData {
    email: string;
    password: string;
}
