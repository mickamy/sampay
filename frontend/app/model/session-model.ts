export interface Token {
  readonly value: string;
  readonly expiresAt: string;
}

export interface Tokens {
  readonly access: Token;
  readonly refresh: Token;
}

export interface Session {
  readonly tokens: Tokens;
}
