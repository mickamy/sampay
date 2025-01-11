export type Either<L, R> = Left<L> | Right<R>;

export class Left<L> {
  constructor(public readonly value: L) {}

  isLeft(): this is Left<L> {
    return true;
  }

  isRight(): this is Right<never> {
    return false;
  }

  map<U>(_: (value: never) => U): Either<L, U> {
    return this;
  }

  mapLeft<U>(fn: (value: L) => U): Either<U, never> {
    return new Left(fn(this.value));
  }
}

export class Right<R> {
  constructor(public readonly value: R) {}

  isLeft(): this is Left<never> {
    return false;
  }

  isRight(): this is Right<R> {
    return true;
  }

  map<U>(fn: (value: R) => U): Either<never, U> {
    return new Right(fn(this.value));
  }

  mapLeft<U>(_: (value: never) => U): Either<U, R> {
    return this;
  }
}
