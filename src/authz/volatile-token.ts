export class VolatileToken {
  #expiresAt: Date | null = null;
  #token: string | null = null;

  update(token: string, expiresAt: Date): void {
    this.#expiresAt = expiresAt;
    this.#token = token;
  }

  getToken(date: Date): string | null {
    if (this.#token === null) {
      return null;
    }
    if (this.expiredOn(date)) {
      return null;
    }
    return this.#token;
  }

  expiredOn(date: Date): boolean {
    if (this.#expiresAt === null) {
      return true;
    }
    return this.#expiresAt > date;
  }
}
