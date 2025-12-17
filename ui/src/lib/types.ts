export interface ProvidersResponse {
    Providers: string[];
    ProvidersMap: Record<string, string>;
}

export interface User {
    AvatarURL: string;
    Name: string;
    NickName: string;
    Email: string;
    Location: string;
    Description: string;
    UserID: string;
    Provider: string;
}
