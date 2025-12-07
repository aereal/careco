export interface TokenProvider {
  refresh(): Promise<void>;
  willExpire(): boolean;
  getToken(): string | null;
}
