import { type User } from "../../../shared/types/user";

export interface Connection {
    id: string;
    user_a_id: string;
    user_b_id: string;
    status: string;
    connected_at: string;
    user_a?: User;
    user_b?: User;
}

export interface Connections {
    connections: Connection[];
    count: number;
}

export interface ConnectionRequest {
    id: string;
    sender_id: string;
    receiver_id: string;
    status: string;
    message?: string | null;
    created_at: string;
    sender?: User;
    receiver?: User;
}

export interface ConnectionRequestResponse {
    message: string;
    request: ConnectionRequest;
}

export interface ConnectionNoContentResponse {
    message: string;
}

export interface ConnectionAcceptResponse {
    message: string;
    connection: Connection;
}

export interface ConnectionRequests {
    requests: ConnectionRequest[];
    count: number;
}

export interface ConnectionError {
    error: string;
    details: string;
}

export interface SendConnectionRequestBody {
  receiver_id: string;
  message?: string;
}


export interface SkipConnectionRequestBody {
  target_userId: string;
}